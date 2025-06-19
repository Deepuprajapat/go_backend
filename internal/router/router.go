package router

import (
	"net/http"
	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/handlers"
	"github.com/VI-IM/im_backend_go/internal/middleware"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

var (
	// Router is the shared router instance
	Router = mux.NewRouter()
)

// Init initializes the router with all routes and middleware
func Init(app application.ApplicationInterface) {
	// Initialize handlers with controller
	authHandler := handlers.NewAuthHandler(app)
	projectHandler := handlers.NewProjectHandler(app)

	// Apply middleware
	Router.Use(middleware.Cors)
	Router.Use(middleware.Logging)
	Router.Use(middleware.Recover)

	// Public routes
	Router.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)

	// auth routes
	Router.Handle("/auth/generate-token", imhttp.AppHandler(authHandler.GenerateToken)).Methods(http.MethodPost)
	Router.Handle("/auth/refresh-token", imhttp.AppHandler(handlers.RefreshToken)).Methods(http.MethodPost)

	// project routes
	Router.Handle("/projects/{project_id}", imhttp.AppHandler(projectHandler.GetProject)).Methods(http.MethodGet)
	Router.Handle("/projects", middleware.Auth(imhttp.AppHandler(projectHandler.AddProject))).Methods(http.MethodPost)
}

