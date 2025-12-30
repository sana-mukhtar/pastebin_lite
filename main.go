package main

import (
	"log"
	"net/http"
	"pastebin/internal"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Routes (OPTIONS explicitly allowed)
	r.HandleFunc("/api/healthz", healthHandler).
		Methods("GET", "OPTIONS")

	r.HandleFunc("/api/pastes", internal.CreatePasteHandler).
		Methods("POST", "OPTIONS")

	r.HandleFunc("/api/pastes/{id}", internal.GetPaste).
		Methods("GET", "OPTIONS")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}
