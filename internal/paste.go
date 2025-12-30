package internal

import (
	"encoding/json"
	"fmt"
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

	// maxViews is optional; if provided, must be >= 1
	if maxViews != 0 && maxViews < 1 {
		return "max_views must be >= 1", false
	}

	// ttl is optional; if provided, must be >= 1
	if ttl != 0 && ttl < 1 {
		return "ttl_seconds must be >= 1", false
	}

	return "", true
}

// CreatePaste stores a new paste and returns its UUID
func CreatePasteHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Content  string `json:"content"`
		MaxViews int    `json:"max_views"`
		TTL      int    `json:"ttl_seconds"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	var msg string
	var ok bool
	// Validate using separate function
	if msg, ok = validatePasteInput(req.Content, req.MaxViews, req.TTL); !ok {
		http.Error(w, `{"error":"`+msg+`"}`, http.StatusBadRequest)
		return
	}

	fmt.Println("Creating paste with content length:", msg, ok)

	id := uuid.New().String()
	paste := &Paste{
		ID:        id,
		Content:   req.Content,
		MaxViews:  req.MaxViews,
		TTL:       req.TTL,
		CreatedAt: time.Now(),
		Views:     0,
	}

	if req.TTL > 0 {
		exp := paste.CreatedAt.Add(time.Duration(req.TTL) * time.Second)
		paste.ExpiresAt = &exp
	}

	pastes[id] = paste

	w.Header().Set("Content-Type", "application/json")
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

	if paste.MaxViews != 0 && paste.Views >= paste.MaxViews {
		http.NotFound(w, r)
		return
	}

	paste.Views++

	var remaining *int
	if paste.MaxViews != 0 {
		r := paste.MaxViews - paste.Views
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
