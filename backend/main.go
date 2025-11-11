package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"event-registration-backend/config"
	"event-registration-backend/firestore"
	"event-registration-backend/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize Firestore
	ctx := context.Background()
	if err := firestore.InitializeFirestore(ctx, cfg.ServiceAccountPath, cfg.GCPProjectID, cfg.ClientID); err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	log.Printf("Firestore initialized with client_id: %s", firestore.ClientID)

	// Set admin password in handlers
	handlers.SetAdminPassword(cfg.AdminPassword)

	// Setup router
	r := mux.NewRouter()

	// Public API routes
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/sessions", handlers.GetSessions).Methods("GET")
	api.HandleFunc("/speakers", handlers.GetSpeakers).Methods("GET")
	api.HandleFunc("/attendees/count", handlers.GetAttendeeCount).Methods("GET")
	api.HandleFunc("/attendees/register", handlers.RegisterAttendee).Methods("POST")

	// Admin API routes
	admin := api.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/login", handlers.AdminLogin).Methods("POST")
	admin.HandleFunc("/attendees", handlers.GetAttendees).Methods("GET")
	admin.HandleFunc("/stats", handlers.GetAttendeeStats).Methods("GET")
	admin.HandleFunc("/speakers", handlers.CreateOrUpdateSpeaker).Methods("POST")
	admin.HandleFunc("/sessions", handlers.CreateOrUpdateSession).Methods("POST")

	// Serve static files from frontend/dist
	frontendDir := cfg.FrontendDir
	if frontendDir == "" {
		frontendDir = "./frontend/dist"
	}

	// Check if frontend directory exists
	if _, err := os.Stat(frontendDir); err == nil {
		// Create a file server for static assets
		// Serve all static files (assets, images, etc.) directly
		fileServer := http.FileServer(http.Dir(frontendDir))
		
		// Serve static files - these routes must be registered before the catch-all
		// Serve assets directory - strip "/" prefix so fileServer looks in frontendDir/assets/
		r.PathPrefix("/assets/").Handler(http.StripPrefix("/", fileServer))
		// Serve vite.svg and other root-level static files
		r.PathPrefix("/vite.svg").Handler(fileServer)
		
		// Serve index.html for all non-API and non-asset routes (SPA routing)
		// Use MatcherFunc to exclude /api and /assets paths
		r.PathPrefix("/").MatcherFunc(func(r *http.Request, rm *mux.RouteMatch) bool {
			// Don't match API routes
			if filepath.HasPrefix(r.URL.Path, "/api") {
				return false
			}
			// Don't match asset routes (already handled above)
			if filepath.HasPrefix(r.URL.Path, "/assets/") {
				return false
			}
			// Don't match vite.svg (already handled above)
			if r.URL.Path == "/vite.svg" {
				return false
			}
			return true
		}).HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			// Check if it's a file request (has extension) - try to serve the actual file
			ext := filepath.Ext(req.URL.Path)
			if ext != "" && ext != "/" {
				filePath := filepath.Join(frontendDir, req.URL.Path)
				if _, err := os.Stat(filePath); err == nil {
					http.ServeFile(w, req, filePath)
					return
				}
			}
			// Otherwise serve index.html for SPA routing
			indexPath := filepath.Join(frontendDir, "index.html")
			http.ServeFile(w, req, indexPath)
		})
		log.Printf("Serving static files from: %s", frontendDir)
	} else {
		log.Printf("Frontend directory not found at %s, skipping static file serving", frontendDir)
	}

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

