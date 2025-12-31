package main

import (
	"log"
	"net/http"
	"os"
	"pastebin/internal"

	"github.com/gorilla/mux"
)

func main() {
	internal.InitDB()

	r := mux.NewRouter()

	// CORS middleware
	r.Use(corsMiddleware)

	// Routes
	r.HandleFunc("/api/healthz", healthHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/pastes", internal.CreatePasteHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/pastes/{id}", internal.GetPaste).Methods("GET", "OPTIONS")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, x-test-now-ms")
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}
