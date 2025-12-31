package main

import (
	"log"
	"net/http"
	"pastebin/internal"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Connect to Postgres
	internal.InitDB()
	internal.AutoMigrate()

	r := mux.NewRouter()
	r.Use(corsMiddleware)

	// Routes
	r.HandleFunc("/api/healthz", healthHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/pastes", internal.CreatePasteHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/pastes/{id}", internal.GetPaste).Methods("GET", "OPTIONS")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

// Simple CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
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
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok": true}`))
}
