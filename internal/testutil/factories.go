package testutil

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/google/uuid"
)

// LocationFactory creates test locations
type LocationFactory struct {
	client *ent.Client
}

// NewLocationFactory creates a new location factory
func NewLocationFactory(client *ent.Client) *LocationFactory {
	return &LocationFactory{client: client}
}

// Create creates a test location
func (f *LocationFactory) Create(t *testing.T) *ent.Location {
	ctx := context.Background()
	
	location, err := f.client.Location.Create().
		SetID(uuid.New().String()).
		SetLocalityName("Test Locality").
		SetCity("Test City").
		SetState("Test State").
		SetCountry("India").
		SetPincode("123456").
		SetPhoneNumber("9876543210").
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test location: %v", err)
	}
	
	return location
}

// CreateWithCustomData creates a location with custom data
func (f *LocationFactory) CreateWithCustomData(t *testing.T, data map[string]interface{}) *ent.Location {
	ctx := context.Background()
	
	builder := f.client.Location.Create().
		SetID(uuid.New().String()).
		SetLocalityName(getStringValue(data, "locality_name", "Test Locality")).
		SetCity(getStringValue(data, "city", "Test City")).
		SetState(getStringValue(data, "state", "Test State")).
		SetCountry(getStringValue(data, "country", "India")).
		SetPincode(getStringValue(data, "pincode", "123456")).
		SetPhoneNumber(getStringValue(data, "phone_number", "9876543210"))
	
	location, err := builder.Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test location: %v", err)
	}
	
	return location
}

// DeveloperFactory creates test developers
type DeveloperFactory struct {
	client *ent.Client
}

// NewDeveloperFactory creates a new developer factory
func NewDeveloperFactory(client *ent.Client) *DeveloperFactory {
	return &DeveloperFactory{client: client}
}

// Create creates a test developer
func (f *DeveloperFactory) Create(t *testing.T) *ent.Developer {
	ctx := context.Background()
	
	developer, err := f.client.Developer.Create().
		SetID(uuid.New().String()).
		SetName("Test Developer").
		SetLegalName("Test Developer Pvt Ltd").
		SetIdentifier("DEV001").
		SetEstablishedYear(2020).
		SetIsVerified(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test developer: %v", err)
	}
	
	return developer
}

// CreateWithCustomData creates a developer with custom data
func (f *DeveloperFactory) CreateWithCustomData(t *testing.T, data map[string]interface{}) *ent.Developer {
	ctx := context.Background()
	
	builder := f.client.Developer.Create().
		SetID(uuid.New().String()).
		SetName(getStringValue(data, "name", "Test Developer")).
		SetLegalName(getStringValue(data, "legal_name", "Test Developer Pvt Ltd")).
		SetIdentifier(getStringValue(data, "identifier", "DEV001")).
		SetEstablishedYear(getIntValue(data, "established_year", 2020)).
		SetIsVerified(getBoolValue(data, "is_verified", true))
	
	developer, err := builder.Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test developer: %v", err)
	}
	
	return developer
}

// ProjectFactory creates test projects
type ProjectFactory struct {
	client *ent.Client
}

// NewProjectFactory creates a new project factory
func NewProjectFactory(client *ent.Client) *ProjectFactory {
	return &ProjectFactory{client: client}
}

// Create creates a test project with default dependencies
func (f *ProjectFactory) Create(t *testing.T) *ent.Project {
	// Create required dependencies
	locationFactory := NewLocationFactory(f.client)
	location := locationFactory.Create(t)
	
	developerFactory := NewDeveloperFactory(f.client)
	developer := developerFactory.Create(t)
	
	return f.CreateWithDependencies(t, location, developer)
}

// CreateWithDependencies creates a project with provided dependencies
func (f *ProjectFactory) CreateWithDependencies(t *testing.T, location *ent.Location, developer *ent.Developer) *ent.Project {
	ctx := context.Background()
	
	project, err := f.client.Project.Create().
		SetID(uuid.New().String()).
		SetName("Test Project").
		SetDescription("Test project description").
		SetStatus("ongoing").
		SetProjectType("RESIDENTIAL").
		SetLocation(location).
		SetDeveloper(developer).
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test project: %v", err)
	}
	
	return project
}

// PropertyFactory creates test properties
type PropertyFactory struct {
	client *ent.Client
}

// NewPropertyFactory creates a new property factory
func NewPropertyFactory(client *ent.Client) *PropertyFactory {
	return &PropertyFactory{client: client}
}

// Create creates a test property with default dependencies
func (f *PropertyFactory) Create(t *testing.T) *ent.Property {
	// Create required dependencies
	projectFactory := NewProjectFactory(f.client)
	project := projectFactory.Create(t)
	
	return f.CreateWithProject(t, project)
}

// CreateWithProject creates a property with provided project
func (f *PropertyFactory) CreateWithProject(t *testing.T, project *ent.Project) *ent.Property {
	ctx := context.Background()
	
	property, err := f.client.Property.Create().
		SetID(uuid.New().String()).
		SetName("Test Property").
		SetPropertyType("apartment").
		SetPrice(1000000).
		SetProject(project).
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test property: %v", err)
	}
	
	return property
}

// LeadsFactory creates test leads
type LeadsFactory struct {
	client *ent.Client
}

// NewLeadsFactory creates a new leads factory
func NewLeadsFactory(client *ent.Client) *LeadsFactory {
	return &LeadsFactory{client: client}
}

// Create creates a test lead
func (f *LeadsFactory) Create(t *testing.T) *ent.Leads {
	ctx := context.Background()
	
	lead, err := f.client.Leads.Create().
		SetName("Test Lead").
		SetEmail("testlead@example.com").
		SetPhoneNumber("9876543210").
		SetMessage("Test lead message").
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test lead: %v", err)
	}
	
	return lead
}

// BlogsFactory creates test blogs
type BlogsFactory struct {
	client *ent.Client
}

// NewBlogsFactory creates a new blogs factory
func NewBlogsFactory(client *ent.Client) *BlogsFactory {
	return &BlogsFactory{client: client}
}

// Create creates a test blog
func (f *BlogsFactory) Create(t *testing.T) *ent.Blogs {
	ctx := context.Background()
	
	blog, err := f.client.Blogs.Create().
		SetID(uuid.New().String()).
		SetBlogURL("test-blog-post").
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test blog: %v", err)
	}
	
	return blog
}

// Helper functions for extracting values from data map
func getStringValue(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultValue
}

func getIntValue(data map[string]interface{}, key string, defaultValue int) int {
	if val, ok := data[key]; ok {
		if i, ok := val.(int); ok {
			return i
		}
	}
	return defaultValue
}

func getBoolValue(data map[string]interface{}, key string, defaultValue bool) bool {
	if val, ok := data[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return defaultValue
}

// RandomString generates a random string of given length
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomEmail generates a random email address
func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomString(8))
}

// RandomPhoneNumber generates a random Indian phone number
func RandomPhoneNumber() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("9%09d", seededRand.Intn(1000000000))
}