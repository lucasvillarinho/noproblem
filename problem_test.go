package noproblem

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestProblemBasicFields(t *testing.T) {
	problem := NewProblem("https://example.com/problems/test", "Test Problem", 400)

	if problem.Type != "https://example.com/problems/test" {
		t.Errorf("Expected Type to be 'https://example.com/problems/test', got '%s'", problem.Type)
	}

	if problem.Title != "Test Problem" {
		t.Errorf("Expected Title to be 'Test Problem', got '%s'", problem.Title)
	}

	if problem.Status != 400 {
		t.Errorf("Expected Status to be 400, got %d", problem.Status)
	}
}

func TestProblemOptions(t *testing.T) {
	problem := NewProblem("https://example.com/problems/test", "Test Problem", 400,
		WithDetail("This is a test problem"),
		WithInstance("/test/123"),
		WithExtra("custom_field", "custom_value"))

	if problem.Detail != "This is a test problem" {
		t.Errorf("Expected Detail to be 'This is a test problem', got '%s'", problem.Detail)
	}

	if problem.Instance != "/test/123" {
		t.Errorf("Expected Instance to be '/test/123', got '%s'", problem.Instance)
	}

	if problem.Extra["custom_field"] != "custom_value" {
		t.Errorf("Expected custom_field to be 'custom_value', got '%v'", problem.Extra["custom_field"])
	}
}

func TestProblemJSONSerialization(t *testing.T) {
	problem := NewProblem("https://example.com/problems/test", "Test Problem", 400,
		WithDetail("This is a test problem"),
		WithInstance("/test/123"),
		WithExtra("error_code", "TEST_001"),
		WithExtra("retry_after", 60))

	data, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("Failed to marshal problem: %v", err)
	}

	var unmarshaled map[string]interface{}
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal problem: %v", err)
	}

	// Check standard fields
	if unmarshaled["type"] != "https://example.com/problems/test" {
		t.Errorf("Expected type field in JSON")
	}

	if unmarshaled["title"] != "Test Problem" {
		t.Errorf("Expected title field in JSON")
	}

	if unmarshaled["status"] != float64(400) {
		t.Errorf("Expected status field in JSON")
	}

	// Check extra fields
	if unmarshaled["error_code"] != "TEST_001" {
		t.Errorf("Expected error_code in JSON")
	}

	if unmarshaled["retry_after"] != float64(60) {
		t.Errorf("Expected retry_after in JSON")
	}
}

func TestProblemJSONFieldOrder(t *testing.T) {
	validationErrors := []map[string]string{
		{"field": "application_name", "detail": "Application name is required"},
		{"field": "domain", "detail": "Domain must be a valid URL"},
		{"field": "callback_url", "detail": "Callback is required and must be a valid URL"},
		{"field": "session_model.session_duration", "detail": "Session duration must be at least 3600 seconds (1 hour)"},
	}

	problem := NewProblem(
		"https://developer.mozilla.org/pt-BR/docs/Web/HTTP/Reference/Status/400",
		"Your request is not valid.",
		http.StatusUnprocessableEntity,
		WithInstance("/api/v1/settings"),
		WithExtra("errors", validationErrors),
	)

	jsonBytes, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("Failed to marshal problem: %v", err)
	}

	jsonString := string(jsonBytes)

	expectedOrder := []string{
		`"type":`,
		`"title":`,
		`"status":`,
		`"instance":`,
		`"errors":`,
	}

	lastIndex := -1
	for _, field := range expectedOrder {
		index := strings.Index(jsonString, field)
		if index == -1 {
			t.Errorf("Field %s not found in JSON", field)
			continue
		}

		if index <= lastIndex {
			t.Errorf("Field %s is not in expected order (index %d, previous %d)", field, index, lastIndex)
		}

		lastIndex = index
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Errorf("Generated JSON is not valid: %v", err)
	}

	if result["type"] != "https://developer.mozilla.org/pt-BR/docs/Web/HTTP/Reference/Status/400" {
		t.Errorf("Unexpected type value")
	}

	if result["title"] != "Your request is not valid." {
		t.Errorf("Unexpected title value")
	}

	if result["status"] != float64(422) {
		t.Errorf("Unexpected status value")
	}

	if result["instance"] != "/api/v1/settings" {
		t.Errorf("Unexpected instance value")
	}
}
