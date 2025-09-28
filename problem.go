package noproblem

import (
	"encoding/json"
	"fmt"
)

const ContentTypeProblemJSON = "application/problem+json"

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
// with guaranteed field order according to RFC 9457
func (p ProblemDetails) MarshalJSON() ([]byte, error) {
	base := []byte(`{`)

	fields := []struct {
		Key   string
		Value any
	}{
		{"type", p.Type},
		{"title", p.Title},
		{"status", p.Status},
		{"detail", p.Detail},
		{"instance", p.Instance},
	}

	first := true
	for _, f := range fields {
		if f.Value == "" || f.Value == 0 {
			continue
		}

		val, err := json.Marshal(f.Value)
		if err != nil {
			return nil, err
		}

		if !first {
			base = append(base, ',')
		}
		base = append(base, fmt.Sprintf(`"%s":%s`, f.Key, val)...)
		first = false
	}

	for k, v := range p.Extra {
		val, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}

		if !first {
			base = append(base, ',')
		}
		base = append(base, fmt.Sprintf(`"%s":%s`, k, val)...)
		first = false
	}

	base = append(base, '}')
	return base, nil
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
