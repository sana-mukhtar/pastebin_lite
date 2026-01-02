package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Paste struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	MaxViews  int        `json:"max_views,omitempty"`
	TTL       int        `json:"ttl_seconds,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Views     int        `json:"views"`
}

// Validate input
func validatePasteInput(content string, maxViews, ttl int) (string, bool) {
	if content == "" {
		return "content cannot be empty", false
	}
	if maxViews < 0 {
		return "max_views cannot be negative", false
	}
	if ttl < 0 {
		return "ttl_seconds cannot be negative", false
	}
	return "", true
}

// Deterministic time for testing
func now(r *http.Request) time.Time {
	if os.Getenv("TEST_MODE") == "1" {
		if v := r.Header.Get("x-test-now-ms"); v != "" {
			if ms, err := strconv.ParseInt(v, 10, 64); err == nil {
				return time.UnixMilli(ms)
			}
		}
	}
	return time.Now()
}

// CreatePasteHandler stores a paste in Postgres
func CreatePasteHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Content  string `json:"content"`
		MaxViews int    `json:"max_views"`
		TTL      int    `json:"ttl_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if msg, ok := validatePasteInput(req.Content, req.MaxViews, req.TTL); !ok {
		http.Error(w, `{"error":"`+msg+`"}`, http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	nowTime := time.Now()
	var expiresAt *time.Time

	if req.TTL > 0 {
		exp := nowTime.Add(time.Duration(req.TTL) * time.Second)
		expiresAt = &exp
	}

	_, err := DB.Exec(
		`INSERT INTO pastes (id, content, max_views, ttl_seconds, created_at, expires_at, views)
		 VALUES ($1,$2,$3,$4,$5,$6,$7)`,
		id, req.Content, req.MaxViews, req.TTL, nowTime, expiresAt, 0,
	)
	if err != nil {
		log.Println("DB insert error:", err)
		http.Error(w, `{"error":"failed to save paste"}`, http.StatusInternalServerError)
		return
	}

	frontendBase := os.Getenv("FRONTEND_URL")
	if frontendBase == "" {
		frontendBase = "http://localhost:3000"
	}
	// Ensure no trailing slash
	frontendBase = strings.TrimRight(frontendBase, "/")

	json.NewEncoder(w).Encode(map[string]string{
		"id":  id,
		"url": frontendBase + "/p/" + id
	})
}

// GetPaste retrieves a paste from Postgres
func GetPaste(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	var paste Paste
	err := DB.QueryRow(`SELECT id, content, max_views, ttl_seconds, created_at, expires_at, views FROM pastes WHERE id=$1`, id).
		Scan(&paste.ID, &paste.Content, &paste.MaxViews, &paste.TTL, &paste.CreatedAt, &paste.ExpiresAt, &paste.Views)
	if err == sql.ErrNoRows {
		jsonError(w)
		return
	} else if err != nil {
		http.Error(w, `{"error":"DB error"}`, http.StatusInternalServerError)
		return
	}

	// TTL check
	if paste.ExpiresAt != nil && now(r).After(*paste.ExpiresAt) {
		_, _ = DB.Exec("DELETE FROM pastes WHERE id=$1", id)
		jsonError(w)
		return
	}

	// Max views check
	if paste.MaxViews > 0 && paste.Views >= paste.MaxViews {
		jsonError(w)
		return
	}

	// Increment view
	paste.Views++
	_, _ = DB.Exec("UPDATE pastes SET views=$1 WHERE id=$2", paste.Views, id)

	var remaining *int
	if paste.MaxViews > 0 {
		r := paste.MaxViews - paste.Views
		if r < 0 {
			r = 0
		}
		remaining = &r
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"content":         paste.Content,
		"remaining_views": remaining,
		"expires_at":      paste.ExpiresAt,
	})
}

// Send 404 JSON
func jsonError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "paste unavailable",
	})
}
