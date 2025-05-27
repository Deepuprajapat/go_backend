package main

import (
	"log"
	"net/http"
	"os"

	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/database"
	"github.com/VI-IM/im_backend_go/internal/router"
)

func main() {
	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	client := database.NewClient(os.Getenv("DB_DSN"))
	defer client.Close()

	// Initialize router
	router.Init()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, router.Router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
