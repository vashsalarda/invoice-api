package server

import (
	"io"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHandler(t *testing.T) {
	// Create a Fiber app for testing
	app := fiber.New()
	// Inject the Fiber app into the server
	s := &FiberServer{App: app}
	// Define a route in the Fiber app
	app.Get("/", s.InvoiceAPIHandler)
	// Create a test HTTP request
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}
	// Perform the request
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}
	// Your test assertions...
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}
	expected := "{\"message\":\"Invoice API\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

// mockDB implements database.Service for testing HealthHandler
type mockDB struct{}

func (m *mockDB) Health() map[string]string {
	return map[string]string{"message": "It's healthy"}
}

func TestHealthHandler(t *testing.T) {
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDB{}}

	app.Get("/health", s.HealthHandler)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	expected := "{\"message\":\"It's healthy\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}

// mockDBDown simulates database being down by returning a different message
type mockDBDown struct{}

func (m *mockDBDown) Health() map[string]string {
	return map[string]string{"message": "db down"}
}

func TestHealthHandler_DBDown(t *testing.T) {
	app := fiber.New()
	s := &FiberServer{App: app, db: &mockDBDown{}}

	app.Get("/health", s.HealthHandler)

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("error creating request. Err: %v", err)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("error making request to server. Err: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", resp.Status)
	}

	expected := "{\"message\":\"db down\"}"
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("error reading response body. Err: %v", err)
	}
	if expected != string(body) {
		t.Errorf("expected response body to be %v; got %v", expected, string(body))
	}
}
