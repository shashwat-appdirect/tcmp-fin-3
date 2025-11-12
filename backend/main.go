package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"github.com/tcmp-fin-3/backend/firestore"
	"github.com/tcmp-fin-3/backend/handlers"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize Firestore
	ctx := context.Background()
	if err := firestore.InitFirestore(ctx); err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	log.Println("Firestore initialized successfully")

	// Setup router
	r := mux.NewRouter()

	// API routes (must be registered before static file serving)
	api := r.PathPrefix("/api").Subrouter()

	// Public routes
	api.HandleFunc("/attendees", handlers.GetAttendees).Methods("GET")
	api.HandleFunc("/attendees", handlers.RegisterAttendee).Methods("POST")
	api.HandleFunc("/sessions", handlers.GetSessions).Methods("GET")
	api.HandleFunc("/speakers", handlers.GetSpeakers).Methods("GET")

	// Admin routes
	api.HandleFunc("/admin/login", handlers.AdminLogin).Methods("POST")
	api.HandleFunc("/admin/attendees", handlers.AdminMiddleware(handlers.GetAdminAttendees)).Methods("GET")
	api.HandleFunc("/admin/attendees/{id}", handlers.AdminMiddleware(handlers.DeleteAttendee)).Methods("DELETE")
	api.HandleFunc("/admin/sessions", handlers.AdminMiddleware(handlers.AddOrUpdateSession)).Methods("POST")
	api.HandleFunc("/admin/sessions/{id}", handlers.AdminMiddleware(handlers.DeleteSession)).Methods("DELETE")
	api.HandleFunc("/admin/speakers", handlers.AdminMiddleware(handlers.AddOrUpdateSpeaker)).Methods("POST")
	api.HandleFunc("/admin/speakers/{id}", handlers.AdminMiddleware(handlers.DeleteSpeaker)).Methods("DELETE")
	api.HandleFunc("/admin/stats", handlers.AdminMiddleware(handlers.GetStats)).Methods("GET")

	// Serve static files (frontend build) - must be after API routes
	staticDir := "./static"
	if _, err := os.Stat(staticDir); err == nil {
		// Serve static files for all non-API routes
		r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(staticDir))))
		log.Println("Serving static files from ./static")
	} else {
		log.Println("Static directory not found, API-only mode")
	}

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

