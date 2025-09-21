# `noproblem`

A Go implementation of **RFC 7807/9457 - Problem Details for HTTP APIs**.

## RFC Background

This library implements the Problem Details standard based on:

- **[RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807)** (2016) - Original Problem Details specification
- **[RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html)** (2023) - Updated specification that obsoletes RFC 7807

## Features

- ✅ Full RFC 9457 compliance
- ✅ Go idiomatic Option pattern
- ✅ JSON serialization with `application/problem+json` media type
- ✅ Extensible with custom fields
- ✅ Zero external dependencies

## Installation

```bash
go get github.com/lucasvillarinho/noproblem
```

## Quick Start

```go
// Basic problem creation using Option pattern
problem := NewProblem("https://example.com/problems/validation", "Validation Error", 400,
    WithDetail("Email is required"),
    WithInstance("/api/users"),
    WithExtra("field", "email"))

// Write as HTTP response
problem.WriteHTTPResponse(w)
```

## Usage Examples

```go
// Multiple options
problem := NewProblem("https://example.com/problems/validation", "Validation Error", 400,
    WithDetail("Multiple validation errors"),
    WithInstance("/api/users/123"),
    WithExtra("email", "Email is required"),
    WithExtra("age", "Age must be greater than 0"))

// Minimal problem
simple := NewProblem("https://example.com/problems/not-found", "Not Found", 404)

// With custom fields
custom := NewProblem("https://example.com/problems/custom", "Custom Problem", 422,
    WithDetail("Custom error occurred"),
    WithExtra("error_code", "CUSTOM_001"),
    WithExtra("retry_after", 60))
```

## Testing

```bash
go test -v
```

## RFC Compliance

**RFC 9457 features implemented:**

- `application/problem+json` media type
- All standard members: `type`, `title`, `status`, `detail`, `instance`
- Extension with additional members via `Extra` map
- Proper JSON serialization with custom marshaling
- HTTP response integration

**Improvements over RFC 7807:**

- Enhanced standardization
- Better interoperability guidance
- Maintains full backward compatibility

## License

MIT License
