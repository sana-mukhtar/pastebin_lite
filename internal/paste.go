package internal

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync"
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

var (
	mu     sync.RWMutex
	pastes = make(map[string]*Paste)
)

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
	now := time.Now()

	paste := &Paste{
		ID:        id,
		Content:   req.Content,
		MaxViews:  req.MaxViews,
		TTL:       req.TTL,
		CreatedAt: now,
		Views:     0,
	}

	if req.TTL > 0 {
		exp := now.Add(time.Duration(req.TTL) * time.Second)
		paste.ExpiresAt = &exp
	}

	mu.Lock()
	pastes[id] = paste
	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]string{
		"id":  id,
		"url": "http://localhost:3000/p/" + id,
	})
}

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

func GetPaste(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	mu.Lock()
	paste, ok := pastes[id]
	if !ok {
		mu.Unlock()
		jsonError(w)
		return
	}

	// TTL check
	if paste.ExpiresAt != nil && now(r).After(*paste.ExpiresAt) {
		delete(pastes, id)
		mu.Unlock()
		jsonError(w)
		return
	}

	// Max views check
	if paste.MaxViews > 0 && paste.Views >= paste.MaxViews {
		mu.Unlock()
		jsonError(w)
		return
	}

	// Count view
	paste.Views++

	var remaining *int
	if paste.MaxViews > 0 {
		r := paste.MaxViews - paste.Views
		if r < 0 {
			r = 0
		}
		remaining = &r
	}

	var expiresAt *time.Time
	if paste.TTL > 0 {
		expiresAt = paste.ExpiresAt
	}

	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"content":         paste.Content,
		"remaining_views": remaining, // null if unlimited
		"expires_at":      expiresAt, // null if no TTL
	})
}

func jsonError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "paste unavailable",
	})
}
