package internal

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Paste struct {
	Content   string
	ExpiresAt *time.Time
	MaxViews  *int
	Views     int
}

var (
	mu     sync.RWMutex
	pastes = make(map[string]*Paste)
)

// CreatePaste stores a new paste and returns its UUID
func CreatePaste(content string, ttlSeconds *int, maxViews *int) string {
	mu.Lock()
	defer mu.Unlock()

	id := uuid.New().String() // generate UUID

	var expiresAt *time.Time
	if ttlSeconds != nil {
		t := time.Now().Add(time.Duration(*ttlSeconds) * time.Second)
		expiresAt = &t
	}

	pastes[id] = &Paste{
		Content:   content,
		ExpiresAt: expiresAt,
		MaxViews:  maxViews,
		Views:     0,
	}
	return id
}

// GetPaste retrieves a paste by ID, checks TTL and MaxViews
func GetPaste(id string) (*Paste, error) {
	mu.Lock()
	defer mu.Unlock()

	p, ok := pastes[id]
	if !ok {
		return nil, errors.New("paste not found")
	}

	// Check TTL
	if p.ExpiresAt != nil && time.Now().After(*p.ExpiresAt) {
		delete(pastes, id)
		return nil, errors.New("paste expired")
	}

	// Check max views
	if p.MaxViews != nil && p.Views >= *p.MaxViews {
		delete(pastes, id)
		return nil, errors.New("paste view limit reached")
	}

	// Count the view
	p.Views++

	return p, nil
}
