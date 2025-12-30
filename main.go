package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"pastebin/internal"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/api/healthz", healthHandler).Methods("GET")
	r.HandleFunc("/api/pastes", createPasteHandler).Methods("POST")
	r.HandleFunc("/api/pastes/{id}", getPasteHandler).Methods("GET")
	r.HandleFunc("/p/{id}", viewPasteHandler).Methods("GET")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}

func createPasteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		Content    string `json:"content"`
		TTLSeconds *int   `json:"ttl_seconds"`
		MaxViews   *int   `json:"max_views"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	if body.Content == "" {
		http.Error(w, `{"error":"content required"}`, http.StatusBadRequest)
		return
	}

	id := internal.CreatePaste(body.Content, body.TTLSeconds, body.MaxViews)
	url := fmt.Sprintf("http://localhost:8080/p/%s", id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":  id,
		"url": url,
	})
}

func getPasteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/pastes/"):]
	p, err := internal.GetPaste(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	var remainingViews *int
	if p.MaxViews != nil {
		remaining := *p.MaxViews - p.Views
		remainingViews = &remaining
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"content":         p.Content,
		"remaining_views": remainingViews,
		"expires_at":      p.ExpiresAt,
	})
}

func viewPasteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/p/"):]
	p, err := internal.GetPaste(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	tmpl := `<html><body><pre>{{.}}</pre></body></html>`
	t := template.Must(template.New("paste").Parse(tmpl))
	t.Execute(w, template.HTMLEscapeString(p.Content))
}
