package router

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

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

func corsMiddleware(_ *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		})
	}
}

// serveStaticFiles serves the React application static files from build directory
func serveStaticFiles(w http.ResponseWriter, r *http.Request) {
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

func serveReactApp(w http.ResponseWriter, r *http.Request) {
	// create a proxy for the frontend
	frontendProxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   "localhost:3000",
	})

	frontendProxy.ServeHTTP(w, r)
}

// Init initializes the router with all routes and middleware
func Init(app application.ApplicationInterface) {
	// Initialize handlers with controller
	handler := handlers.NewHandler(app)

	// Add middleware
	Router.Use(middleware.LoggingMiddleware)
	Router.Use(corsMiddleware(Router))

	// Public routes
	Router.HandleFunc("/health", handlers.HealthCheck).Methods(http.MethodGet)

	// auth routes
	Router.Handle("/v1/api/auth/generate-token", imhttp.AppHandler(handler.GenerateToken)).Methods(http.MethodPost)
	Router.Handle("/v1/api/auth/refresh-token", imhttp.AppHandler(handler.RefreshToken)).Methods(http.MethodPost)
	Router.Handle("/v1/api/auth/signup", imhttp.AppHandler(handler.Signup)).Methods(http.MethodPost)

	// project routes - specific routes must come before wildcard routes
	Router.Handle("/v1/api/projects/compare", imhttp.AppHandler(handler.CompareProjects)).Methods(http.MethodPost)
	Router.Handle("/v1/api/projects/names", middleware.Auth(imhttp.AppHandler(handler.GetProjectNames))).Methods(http.MethodGet)
	Router.Handle("/v1/api/projects/{project_id}", imhttp.AppHandler(handler.GetProject)).Methods(http.MethodGet)
	Router.Handle("/v1/api/projects", imhttp.AppHandler(handler.ListProjects)).Methods(http.MethodGet)

	Router.Handle("/v1/api/projects/compare", imhttp.AppHandler(handler.CompareProjects)).Methods(http.MethodPost)
	Router.Handle("/v1/api/projects/{slug}", imhttp.AppHandler(handler.GetProjectBySlug)).Methods(http.MethodGet)

	// project internal routes
	Router.Handle("/v1/api/internal/projects", imhttp.AppHandler(handler.AddProject)).Methods(http.MethodPost)                   // internal
	Router.Handle("/v1/api/internal/projects/{project_id}", imhttp.AppHandler(handler.UpdateProject)).Methods(http.MethodPatch)                        // internal
	Router.Handle("/v1/api/internal/projects/{project_id}", middleware.RequireDM(imhttp.AppHandler(handler.DeleteProject))).Methods(http.MethodDelete) // internal
	Router.Handle("/v1/api/internal/projects/filters", imhttp.AppHandler(handler.GetProjectFilters)).Methods(http.MethodGet)                           // internal

	// upload file routes
	Router.Handle("/v1/api/upload", imhttp.AppHandler(handler.UploadFile)).Methods(http.MethodPost)

	// property routes
	Router.Handle("/v1/api/projects/{project_id}/properties", imhttp.AppHandler(handler.GetPropertiesOfProject)).Methods(http.MethodGet)
	Router.Handle("/v1/api/properties/{property_id}", imhttp.AppHandler(handler.GetProperty)).Methods(http.MethodGet)
	Router.Handle("/v1/api/properties", imhttp.AppHandler(handler.ListProperties)).Methods(http.MethodGet)

	//internal routes
	// Protected property internal routes (require business_partner or superadmin role)
	Router.Handle("/v1/api/internal/properties", imhttp.AppHandler(handler.AddProperty)).Methods(http.MethodPost)
	Router.Handle("/v1/api/internal/properties/{property_id}", middleware.RequireBusinessPartner(imhttp.AppHandler(handler.UpdateProperty))).Methods(http.MethodPatch)
	Router.Handle("/v1/api/internal/properties/{property_id}", middleware.RequireBusinessPartner(imhttp.AppHandler(handler.DeleteProperty))).Methods(http.MethodDelete)

	// internal routes
	// Admin property route - business partners see only their properties, superadmins see all properties
	Router.Handle("/v1/api/internal/admin/dashboard/properties", middleware.RequireBusinessPartner(imhttp.AppHandler(handler.AdminListProperties))).Methods(http.MethodGet)

	// developer routes
	Router.Handle("/v1/api/developers", imhttp.AppHandler(handler.ListDevelopers)).Methods(http.MethodGet)
	Router.Handle("/v1/api/developers/{developer_id}", imhttp.AppHandler(handler.GetDeveloper)).Methods(http.MethodGet)
	Router.Handle("/v1/api/developers/{developer_id}", imhttp.AppHandler(handler.DeleteDeveloper)).Methods(http.MethodDelete)

	// location routes
	Router.Handle("/v1/api/locations", imhttp.AppHandler(handler.ListLocations)).Methods(http.MethodGet)
	Router.Handle("/v1/api/locations/{location_id}", imhttp.AppHandler(handler.GetLocation)).Methods(http.MethodGet)

	// location internal routes
	Router.Handle("/v1/api/internal/location", imhttp.AppHandler(handler.AddLocation)).Methods(http.MethodPost)
	Router.Handle("/v1/api/locations/{location_id}", imhttp.AppHandler(handler.DeleteLocation)).Methods(http.MethodDelete)

	// internal amenity routes
	Router.Handle("/v1/api/internal/amenities", imhttp.AppHandler(handler.GetAllCategoriesWithAmenities)).Methods(http.MethodGet)
	Router.Handle("/v1/api/internal/static-site-data", imhttp.AppHandler(handler.UpdateStaticSiteData)).Methods(http.MethodPatch)

	// blog routes
	Router.Handle("/v1/api/blogs", imhttp.AppHandler(handler.ListBlogs)).Methods(http.MethodGet)
	Router.Handle("/v1/api/blogs/{blog_id}", imhttp.AppHandler(handler.GetBlog)).Methods(http.MethodGet)

	// internal blog routes
	Router.Handle("/v1/api/internal/blogs", imhttp.AppHandler(handler.CreateBlog)).Methods(http.MethodPost)
	Router.Handle("/v1/api/internal/blogs/{blog_id}", imhttp.AppHandler(handler.DeleteBlog)).Methods(http.MethodDelete)
	Router.Handle("/v1/api/internal/blogs/{blog_id}", imhttp.AppHandler(handler.UpdateBlog)).Methods(http.MethodPatch)

	// URL availability checking route
	Router.Handle("/v1/api/internal/check-avialable-url", imhttp.AppHandler(handler.CheckURLExists)).Methods(http.MethodGet)

	// lead routes - public endpoints for lead creation and OTP operations
	Router.Handle("/v1/api/leads/send-otp", imhttp.AppHandler(handler.CreateLeadWithOTP)).Methods(http.MethodPost)
	Router.Handle("/v1/api/leads", imhttp.AppHandler(handler.CreateLead)).Methods(http.MethodPost)
	Router.Handle("/v1/api/leads/validate-otp", imhttp.AppHandler(handler.ValidateOTP)).Methods(http.MethodPatch)
	Router.Handle("/v1/api/leads/resend-otp", imhttp.AppHandler(handler.ResendOTP)).Methods(http.MethodPatch)

	// Protected lead routes - only dm role can access lead data
	Router.Handle("/v1/api/leads/get/by/{id}", middleware.RequireDM(imhttp.AppHandler(handler.GetLeadByID))).Methods(http.MethodGet)
	Router.Handle("/v1/api/leads", middleware.RequireDM(imhttp.AppHandler(handler.GetAllLeads))).Methods(http.MethodGet)

	//content routes
	Router.Handle("/v1/api/content/test/{url}", imhttp.AppHandler(handler.GetProjectSEOContent)).Methods(http.MethodGet)
	Router.Handle("/v1/api/content/text", imhttp.AppHandler(handler.GetPropertySEOContent)).Methods(http.MethodGet)
	Router.Handle("/v1/api/content/text/html", imhttp.AppHandler(handler.GetHTMLContent)).Methods(http.MethodGet)

	// generic search routes
	Router.Handle("/v1/api/s/{slug}", imhttp.AppHandler(handler.GetCustomSearchPage)).Methods(http.MethodGet)
	Router.Handle("/v1/api/links", imhttp.AppHandler(handler.GetLinks)).Methods(http.MethodGet)

	// internal route for generic search page
	Router.Handle("/v1/api/internal/custom-search-page", imhttp.AppHandler(handler.GetAllCustomSearchPages)).Methods(http.MethodGet)
	Router.Handle("/v1/api/internal/custom-search-page", imhttp.AppHandler(handler.AddCustomSearchPage)).Methods(http.MethodPost)
	Router.Handle("/v1/api/internal/custom-search-page/{id}", imhttp.AppHandler(handler.UpdateCustomSearchPage)).Methods(http.MethodPatch)
	Router.Handle("/v1/api/internal/custom-search-page/{id}", imhttp.AppHandler(handler.DeleteCustomSearchPage)).Methods(http.MethodDelete)

	// Catch-all route for React app - must be last to handle all non-API routes
	Router.PathPrefix("/").HandlerFunc(serveReactApp) // Proxy to local dev server
	//Router.PathPrefix("/").HandlerFunc(serveStaticFiles) // Serve static build files
}
