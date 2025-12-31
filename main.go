package main

import (
	"log"
	"net/http"
	"pastebin/internal"

	"github.com/gorilla/mux"
)

func main() {
	internal.InitDB()

	r := mux.NewRouter()
	r.Use(cors)

	r.HandleFunc("/api/healthz", health).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/pastes", internal.CreatePasteHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/pastes/{id}", internal.GetPaste).Methods("GET", "OPTIONS")

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"ok":true}`))
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
