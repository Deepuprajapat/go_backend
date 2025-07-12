package router

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/handlers"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/gorilla/mux"
)

var (
	// Router is the shared router instance
	Router = mux.NewRouter()
)

// serveReactApp serves the React application static files
func serveReactApp(w http.ResponseWriter, r *http.Request) {
	// Path to React build directory
	buildDir := "./build"

	// Check if file exists in build directory
	filePath := filepath.Join(buildDir, r.URL.Path)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// If file doesn't exist, serve index.html for client-side routing
		filePath = filepath.Join(buildDir, "index.html")
	}

	http.ServeFile(w, r, filePath)
}

// Init initializes the router with all routes and middleware
func Init(app application.ApplicationInterface) {
	// Initialize handlers with controller
	handler := handlers.NewHandler(app)

	

	// Public routes
	Router.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)

	// auth routes
	Router.Handle("/v1/api/auth/generate-token", imhttp.AppHandler(handler.GenerateToken)).Methods(http.MethodPost)
	Router.Handle("/v1/api/auth/refresh-token", imhttp.AppHandler(handler.RefreshToken)).Methods(http.MethodPost)
	Router.Handle("/v1/api/auth/signup", imhttp.AppHandler(handler.Signup)).Methods(http.MethodPost)

	// project routes
	Router.Handle("/v1/api/projects/{project_id}", imhttp.AppHandler(handler.GetProject)).Methods(http.MethodGet)
	Router.Handle("/v1/api/projects", imhttp.AppHandler(handler.AddProject)).Methods(http.MethodPost)
	Router.Handle("/v1/api/projects/{project_id}", imhttp.AppHandler(handler.UpdateProject)).Methods(http.MethodPatch)
	Router.Handle("/v1/api/projects/{project_id}", imhttp.AppHandler(handler.DeleteProject)).Methods(http.MethodDelete)
	Router.Handle("/v1/api/projects", imhttp.AppHandler(handler.ListProjects)).Methods(http.MethodGet)
	Router.Handle("/v1/api/projects/compare", imhttp.AppHandler(handler.CompareProjects)).Methods(http.MethodPost)

	// upload file routes
	Router.Handle("/v1/api/upload", imhttp.AppHandler(handler.UploadFile)).Methods(http.MethodPost)

	// property routes
	Router.Handle("/v1/api/projects/{project_id}/properties", imhttp.AppHandler(handler.GetPropertiesOfProject)).Methods(http.MethodGet)
	Router.Handle("/v1/api/properties/{property_id}", imhttp.AppHandler(handler.GetProperty)).Methods(http.MethodGet)
	Router.Handle("/v1/api/properties/{property_id}", imhttp.AppHandler(handler.UpdateProperty)).Methods(http.MethodPatch)
	Router.Handle("/v1/api/properties", imhttp.AppHandler(handler.AddProperty)).Methods(http.MethodPost)
	Router.Handle("/v1/api/properties", imhttp.AppHandler(handler.ListProperties)).Methods(http.MethodGet)
	Router.Handle("/v1/api/properties/{property_id}", imhttp.AppHandler(handler.DeleteProperty)).Methods(http.MethodDelete)

	// developer routes
	Router.Handle("/v1/api/developers", imhttp.AppHandler(handler.ListDevelopers)).Methods(http.MethodGet)
	Router.Handle("/v1/api/developers/{developer_id}", imhttp.AppHandler(handler.GetDeveloper)).Methods(http.MethodGet)
	Router.Handle("/v1/api/developers/{developer_id}", imhttp.AppHandler(handler.DeleteDeveloper)).Methods(http.MethodDelete)

	// location routes
	Router.Handle("/v1/api/locations", imhttp.AppHandler(handler.ListLocations)).Methods(http.MethodGet)
	Router.Handle("/v1/api/locations/{location_id}", imhttp.AppHandler(handler.GetLocation)).Methods(http.MethodGet)
	Router.Handle("/v1/api/locations/{location_id}", imhttp.AppHandler(handler.DeleteLocation)).Methods(http.MethodDelete)

	// amenity routes
	// get GetAllCategoriesWithAmenities
	// add CategoryWithAmenities -- done
	// patch static site data

	Router.Handle("/v1/api/amenities", imhttp.AppHandler(handler.GetAllCategoriesWithAmenities)).Methods(http.MethodGet)
	Router.Handle("/v1/api/static-site-data", imhttp.AppHandler(handler.UpdateStaticSiteData)).Methods(http.MethodPatch)

	// blog routes
	Router.Handle("/v1/api/blogs", imhttp.AppHandler(handler.ListBlogs)).Methods(http.MethodGet)
	Router.Handle("/v1/api/blogs/{blog_id}", imhttp.AppHandler(handler.GetBlog)).Methods(http.MethodGet)
	Router.Handle("/v1/api/blogs", imhttp.AppHandler(handler.CreateBlog)).Methods(http.MethodPost)
	Router.Handle("/v1/api/blogs/{blog_id}", imhttp.AppHandler(handler.DeleteBlog)).Methods(http.MethodDelete)
	Router.Handle("/v1/api/blogs/{blog_id}", imhttp.AppHandler(handler.UpdateBlog)).Methods(http.MethodPatch)

	//lead routes

	//content routes
	Router.Handle("/v1/api/content/test/{url}", imhttp.AppHandler(handler.GetProjectSEOContent)).Methods(http.MethodGet)
	Router.Handle("/v1/api/content/text", imhttp.AppHandler(handler.GetPropertySEOContent)).Methods(http.MethodGet)
	Router.Handle("/v1/api/content/text/html", imhttp.AppHandler(handler.GetHTMLContent)).Methods(http.MethodGet)

	
	Router.PathPrefix("/").HandlerFunc(serveReactApp)

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
