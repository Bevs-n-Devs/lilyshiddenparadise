package testutil

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/Bevs-n-Devs/lilyshiddenparadise/handlers"
)

// TestEnv holds the test environment
type TestEnv struct {
	DB        *MockDB
	Templates *template.Template
	Server    *httptest.Server
}

// Global test environment that can be used across tests
var TestEnvironment TestEnv

// SetupTestMain initializes the test environment
func SetupTestMain() {
	// Create a mock database
	mockDB := NewMockDB()

	// Set up test landlord
	mockDB.CreateNewLandlord("test@example.com", "password123")

	// Set up test tenant
	mockDB.CreateNewTenant("tenant@example.com", "password123", "Single Room", "2025-01-01", "1", "1000", "USD")

	// Set up test tenant application
	mockDB.SaveTenantApplicationForm(
		"John Doe",
		"1990-01-01",
		"AB123456",
		"1234567890",
		"applicant@example.com",
		"Software Engineer",
		"Tech Corp",
		"0987654321",
		"Jane Doe",
		"1234567890",
		"123 Emergency St",
		"No",
		"",
		"No",
		"",
		"No",
		"No",
		"Yes",
		"ABC123",
		"No",
		"",
		"No",
		"",
		"No",
		"",
	)

	// Load templates
	templatesDir := filepath.Join("..", "templates")
	templates, err := template.ParseGlob(filepath.Join(templatesDir, "*.html"))
	if err != nil {
		// If templates can't be loaded, create a minimal template for testing
		templates = template.New("test")
		template.Must(templates.New("home.html").Parse("Home Page"))
		template.Must(templates.New("landlord_dashboard.html").Parse("Landlord Dashboard"))
		template.Must(templates.New("tenant_dashboard.html").Parse("Tenant Dashboard"))
		template.Must(templates.New("login_landlord.html").Parse("Landlord Login"))
		template.Must(templates.New("login_tenant.html").Parse("Tenant Login"))
	}

	// Set up the global handlers.Templates variable
	handlers.Templates = templates

	// Create the test environment
	TestEnvironment = TestEnv{
		DB:        mockDB,
		Templates: templates,
	}
}

// TeardownTestMain cleans up the test environment
func TeardownTestMain() {
	if TestEnvironment.Server != nil {
		TestEnvironment.Server.Close()
	}
}

// NewTestRequest creates a new HTTP request for testing
func NewTestRequest(method, path string, body string) *http.Request {
	req := httptest.NewRequest(method, path, nil)
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Body = http.NoBody // Replace with actual body if needed
	}
	return req
}

// ExecuteTestRequest executes a test request and returns the response
func ExecuteTestRequest(t *testing.T, handler http.HandlerFunc, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	return rr
}

// InitTestEnv initializes the test environment for a single test
func InitTestEnv() {
	// Set up test environment if it hasn't been initialized
	if TestEnvironment.DB == nil {
		SetupTestMain()
	}
}

// TestMainWrapper is the entry point for all tests in the package that imports this
func TestMainWrapper(m *testing.M) {
	if m == nil {
		// This is being called from a test, not as the actual TestMain
		InitTestEnv()
		return
	}
	
	// Set up test environment
	SetupTestMain()

	// Run tests
	code := m.Run()

	// Clean up
	TeardownTestMain()

	// Exit with the test status code
	os.Exit(code)
}
