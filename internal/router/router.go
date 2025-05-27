package router

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/handlers"
	"github.com/VI-IM/im_backend_go/internal/middleware"
	"github.com/gorilla/mux"
)

var (
	// Router is the shared router instance
	Router = mux.NewRouter()
)

// Init initializes the router with all routes and middleware
func Init() {
	// Apply middleware
	Router.Use(middleware.Logging)
	Router.Use(middleware.Recover)

	// Public routes
	Router.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)
	Router.HandleFunc("/auth/token", handlers.GenerateToken).Methods(http.MethodPost)

	// Protected routes
	admin := Router.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.Auth)
	// admin.HandleFunc("/users", handlers.CreateUser).Methods(http.MethodPost)
	// admin.HandleFunc("/users", handlers.ListUsers).Methods(http.MethodGet)
	// admin.HandleFunc("/users/{id}", handlers.GetUser).Methods(http.MethodGet)
	// admin.HandleFunc("/users/{id}", handlers.UpdateUser).Methods(http.MethodPut)
	// admin.HandleFunc("/users/{id}", handlers.DeleteUser).Methods(http.MethodDelete)
}
