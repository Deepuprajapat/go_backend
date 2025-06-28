package router

import (
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/handlers"
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
	propertyHandler := handlers.NewPropertyHandler(app)
	locationHandler := handlers.NewLocationHandler(app)

	// Public routes
	Router.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)

	// auth routes
	Router.Handle("/v1/api/auth/generate-token", imhttp.AppHandler(authHandler.GenerateToken)).Methods(http.MethodPost)
	Router.Handle("/v1/api/auth/refresh-token", imhttp.AppHandler(handlers.RefreshToken)).Methods(http.MethodPost)

	// project routes
	Router.Handle("/v1/api/projects/{project_id}", imhttp.AppHandler(projectHandler.GetProject)).Methods(http.MethodGet)
	Router.Handle("/v1/api/projects", imhttp.AppHandler(projectHandler.AddProject)).Methods(http.MethodPost)
	Router.Handle("/v1/api/projects/{project_id}", imhttp.AppHandler(projectHandler.UpdateProject)).Methods(http.MethodPatch)
	Router.Handle("/v1/api/projects/{project_id}", imhttp.AppHandler(projectHandler.DeleteProject)).Methods(http.MethodDelete)
	Router.Handle("/v1/api/projects", imhttp.AppHandler(projectHandler.ListProjects)).Methods(http.MethodGet)

	// property routes
	Router.Handle("/v1/api/projects/{project_id}/properties", imhttp.AppHandler(propertyHandler.GetPropertyOfProject)).Methods(http.MethodGet)
	Router.Handle("/v1/api/properties/{property_id}", imhttp.AppHandler(propertyHandler.GetProperty)).Methods(http.MethodGet)
	Router.Handle("/v1/api/properties/{property_id}", imhttp.AppHandler(propertyHandler.UpdateProperty)).Methods(http.MethodPatch)

	// location routes
	Router.Handle("/v1/api/locations", imhttp.AppHandler(locationHandler.ListLocations)).Methods(http.MethodGet)
}

/////   curl calls	/////
//// ----- Get Project -----
// curl -X GET http://localhost:9999/v1/api/projects/your-project-id

//// ----- Add Project -----
// curl -X POST http://localhost:9999/v1/api/projects \
// -H "Content-Type: application/json" \
// -d '{
//     "project_name": "Sample Project",
//     "project_url": "https://example.com/project",
//     "project_type": "Residential",
//     "locality": "Downtown",
//     "project_city": "Mumbai",
//     "developer_id": "dev123"
// }'

//// ----- Update Project -----
// curl -X PATCH http://localhost:9999/v1/api/projects/your-project-id \
// -H "Content-Type: application/json" \
// -d '{
//     "project_name": "Updated Project Name",
//     "description": "Updated project description",
//     "status": "ACTIVE",
//     "min_price": 5000000,
//     "max_price": 10000000,
//     "price_unit": "INR"
// }'
