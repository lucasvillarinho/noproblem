package main

import (
	"encoding/json"
	"net/http/httptest"
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

func TestProblemError(t *testing.T) {
	// Test with detail
	problem1 := NewProblem("https://example.com/problems/test", "Test Problem", 400,
		WithDetail("This is a test problem"))

	expected1 := "Test Problem: This is a test problem"
	if problem1.Error() != expected1 {
		t.Errorf("Expected Error() to return '%s', got '%s'", expected1, problem1.Error())
	}

	// Test without detail
	problem2 := NewProblem("https://example.com/problems/test", "Test Problem", 400)

	expected2 := "Test Problem"
	if problem2.Error() != expected2 {
		t.Errorf("Expected Error() to return '%s', got '%s'", expected2, problem2.Error())
	}
}

func TestProblemHTTPResponse(t *testing.T) {
	problem := NewProblem("https://example.com/problems/test", "Test Problem", 400,
		WithDetail("This is a test problem"))

	recorder := httptest.NewRecorder()
	err := problem.WriteHTTPResponse(recorder)

	if err != nil {
		t.Fatalf("WriteHTTPResponse failed: %v", err)
	}

	// Check status code
	if recorder.Code != 400 {
		t.Errorf("Expected status code 400, got %d", recorder.Code)
	}

	// Check content type
	contentType := recorder.Header().Get("Content-Type")
	if contentType != "application/problem+json" {
		t.Errorf("Expected Content-Type 'application/problem+json', got '%s'", contentType)
	}

	// Check body
	var response map[string]interface{}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["type"] != "https://example.com/problems/test" {
		t.Errorf("Expected type in response")
	}
}
