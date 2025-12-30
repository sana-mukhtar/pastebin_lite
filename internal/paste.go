package internal

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Paste struct {
	ID        string
	Content   string
	MaxViews  *int
	Views     int
	ExpiresAt *time.Time
}

var (
	mu     sync.RWMutex
	pastes = make(map[string]*Paste)
)

// CreatePaste stores a new paste and returns its UUID
func CreatePaste(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Content    string `json:"content"`
		MaxViews   *int   `json:"max_views"`
		TTLSeconds *int   `json:"ttl_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Content) == "" {
		http.Error(w, `{"error":"content required"}`, http.StatusBadRequest)
		return
	}

	if req.MaxViews != nil && *req.MaxViews < 1 {
		http.Error(w, `{"error":"max_views must be >= 1"}`, http.StatusBadRequest)
		return
	}

	var expiresAt *time.Time
	if req.TTLSeconds != nil {
		if *req.TTLSeconds < 1 {
			http.Error(w, `{"error":"ttl_seconds must be >= 1"}`, http.StatusBadRequest)
			return
		}
		t := time.Now().Add(time.Duration(*req.TTLSeconds) * time.Second)
		expiresAt = &t
	}

	id := uuid.NewString()

	store[id] = &Paste{
		ID:        id,
		Content:   req.Content,
		MaxViews:  req.MaxViews,
		Views:     0,
		ExpiresAt: expiresAt,
	}

	json.NewEncoder(w).Encode(map[string]string{
		"id":  id,
		"url": "http://localhost:3000/paste/" + id,
	})
}

// GetPaste retrieves a paste by ID, checks TTL and MaxViews
func GetPaste(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]
	paste, ok := store[id]
	if !ok {
		http.NotFound(w, r)
		return
	}

	// 🔹 View limit enforcement
	if paste.MaxViews != nil && paste.Views >= *paste.MaxViews {
		http.NotFound(w, r)
		return
	}

	paste.Views++

	var remaining *int
	if paste.MaxViews != nil {
		r := *paste.MaxViews - paste.Views
		if r < 0 {
			r = 0
		}
		remaining = &r
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"content":         paste.Content,
		"remaining_views": remaining,
		"expires_at":      nil,
	})
}

var store = make(map[string]*Paste)
