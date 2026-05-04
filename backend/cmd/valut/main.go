package main

import (
	"backend/internal/config"
	"backend/internal/db"
	"backend/internal/handlers"
	"backend/internal/midleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	database, err := db.Init(cfg)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	defer database.Close()

	handlers.SetDB(database)

	r := mux.NewRouter()

	r.HandleFunc("/signup", handlers.SignUpHandler).Methods("POST")
	r.HandleFunc("/login", handlers.LogInHandler).Methods("POST")

	protected := r.PathPrefix("/vault").Subrouter()
	protected.Use(midleware.JWTMidleware)
	protected.HandleFunc("", handlers.CreateEntryHandler).Methods("POST")
	protected.HandleFunc("/{sign}/{date}", handlers.GetEntryHandler).Methods("GET")

	port := cfg.ServerPort
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
