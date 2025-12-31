package internal

import (
	"encoding/json"
	"net/http"
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
		"url": "http://localhost:3000/paste/" + id,
	})
}

func GetPaste(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	mu.Lock()
	paste, ok := pastes[id]
	if !ok {
		mu.Unlock()
		http.NotFound(w, r)
		return
	}

	// TTL check
	if paste.ExpiresAt != nil && time.Now().After(*paste.ExpiresAt) {
		delete(pastes, id)
		mu.Unlock()
		http.NotFound(w, r)
		return
	}

	// Max views check
	if paste.MaxViews > 0 && paste.Views >= paste.MaxViews {
		mu.Unlock()
		http.NotFound(w, r)
		return
	}

	paste.Views++

	var remaining *int
	if paste.MaxViews > 0 {
		r := paste.MaxViews - paste.Views
		if r < 0 {
			r = 0
		}
		remaining = &r
	}

	mu.Unlock()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"content":         paste.Content,
		"remaining_views": remaining,
		"expires_at":      paste.ExpiresAt,
	})
}
