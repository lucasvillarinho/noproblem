package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
)

// ProblemDetails represents a problem detail as defined in RFC 9457
type ProblemDetails struct {
	Type     string         `json:"type,omitempty"`
	Title    string         `json:"title,omitempty"`
	Status   int            `json:"status,omitempty"`
	Detail   string         `json:"detail,omitempty"`
	Instance string         `json:"instance,omitempty"`
	Extra    map[string]any `json:"-"`
}

// MarshalJSON implements custom JSON marshaling to include extra fields
func (p ProblemDetails) MarshalJSON() ([]byte, error) {
	type Alias ProblemDetails
	aux := struct {
		Alias
	}{
		Alias: Alias(p),
	}

	data, err := json.Marshal(aux.Alias)
	if err != nil {
		return nil, err
	}

	if len(p.Extra) == 0 {
		return data, nil
	}

	var base map[string]interface{}
	if err := json.Unmarshal(data, &base); err != nil {
		return nil, err
	}

	maps.Copy(base, p.Extra)

	return json.Marshal(base)
}

// Option is a function that modifies a ProblemDetails
type Option func(*ProblemDetails)

// WithDetail sets the detail message
func WithDetail(detail string) Option {
	return func(p *ProblemDetails) {
		p.Detail = detail
	}
}

// WithInstance sets the instance URI
func WithInstance(instance string) Option {
	return func(p *ProblemDetails) {
		p.Instance = instance
	}
}

// WithExtra adds an extra field to the problem
func WithExtra(key string, value any) Option {
	return func(p *ProblemDetails) {
		if p.Extra == nil {
			p.Extra = make(map[string]any)
		}
		p.Extra[key] = value
	}
}

// NewProblem creates a new Problem with the given parameters and options
func NewProblem(problemType, title string, status int, opts ...Option) *ProblemDetails {
	p := &ProblemDetails{
		Type:   problemType,
		Title:  title,
		Status: status,
		Extra:  make(map[string]any),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// WriteHTTPResponse writes the problem as an HTTP response
func (p *ProblemDetails) WriteHTTPResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/problem+json")

	if p.Status > 0 {
		w.WriteHeader(p.Status)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(w)
	return encoder.Encode(p)
}

// Error implements the error interface
func (p ProblemDetails) Error() string {
	if p.Detail != "" {
		return fmt.Sprintf("%s: %s", p.Title, p.Detail)
	}
	return p.Title
}
