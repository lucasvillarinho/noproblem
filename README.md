# `🚩 noproblem`

A Go implementation of **RFC 7807/9457 - Problem Details for HTTP APIs**.

<a href="https://github.com/lucasvillarinho/noproblem/blob/master/LICENSE"><img src="https://img.shields.io/github/license/lucasvillarinho/noproblem" alt="License"></a>
<a href="https://github.com/lucasvillarinho/noproblem"><img src="https://img.shields.io/github/go-mod/go-version/lucasvillarinho/noproblem" alt="GitHub go.mod Go version"></a>
<a href="https://goreportcard.com/report/github.com/lucasvillarinho/noproblem"><img src="https://goreportcard.com/badge/github.com/lucasvillarinho/noproblem" alt="Go Report Card"></a>

## RFC Background

This library implements the Problem Details standard based on:

- **[RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807)** (2016) - Original Problem Details specification
- **[RFC 9457](https://www.rfc-editor.org/rfc/rfc9457.html)** (2023) - Updated specification that obsoletes RFC 7807

## Features

- Full RFC 9457 compliance
- Go idiomatic Option pattern
- JSON serialization with `application/problem+json` media type
- Extensible with custom fields
- Zero external dependencies

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

<details>
<summary><b>Echo</b></summary>

```go
package main

import (
    "net/http"

    "github.com/labstack/echo/v4"
    np "github.com/lucasvillarinho/noproblem"
)

func main() {
    e := echo.New()

    e.POST("/users", func(c echo.Context) error {
        problem := np.NewProblem(
            "https://example.com/problems/validation",
            "Validation Error",
            400,
            np.WithDetail("Email is required"),
            np.WithInstance(c.Request().URL.Path),
            np.WithExtra("field", "email"),
        )

        c.Response().Header().Set("Content-Type", np.ContentTypeProblemJSON)
        c.Response().WriteHeader(problem.Status)
        return c.JSON(problem.Status, problem)
    })

    e.Logger.Fatal(e.Start(":8080"))
}
```

</details>

<details>
<summary><b>Gin</b></summary>

```go
package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
    np "github.com/lucasvillarinho/noproblem"
)

func main() {
    r := gin.Default()

    r.POST("/users", func(c *gin.Context) {
        problem := np.NewProblem(
            "https://example.com/problems/validation",
            "Validation Error",
            400,
            np.WithDetail("Email is required"),
            np.WithInstance(c.Request.URL.Path),
            np.WithExtra("field", "email"),
        )

        c.Header("Content-Type", np.ContentTypeProblemJSON)
        c.JSON(problem.Status, problem)
    })

    r.Run(":8080")
}
```

</details>

<details>
<summary><b>Fiber</b></summary>

```go
package main

import (
    "github.com/gofiber/fiber/v2"
    np "github.com/lucasvillarinho/noproblem"
)

func main() {
    app := fiber.New()

    app.Post("/users", func(c *fiber.Ctx) error {
        problem := np.NewProblem(
            "https://example.com/problems/validation",
            "Validation Error",
            400,
            np.WithDetail("Email is required"),
            np.WithInstance(c.Path()),
            np.WithExtra("field", "email"),
        )

        c.Set("Content-Type", np.ContentTypeProblemJSON)
        return c.Status(problem.Status).JSON(problem)
    })

    app.Listen(":8080")
}
```

</details>

<details>
<summary><b>Standard HTTP</b></summary>

```go
package main

import (
    "log"
    "net/http"

    np "github.com/lucasvillarinho/noproblem"
)

func main() {
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            problem := np.NewProblem(
                "https://example.com/problems/method-not-allowed",
                "Method Not Allowed",
                405,
                np.WithDetail("Only POST method is allowed"),
                np.WithInstance(r.URL.Path),
            )

            w.Header().Set("Content-Type", np.ContentTypeProblemJSON)
            w.WriteHeader(problem.Status)
            problem.WriteHTTPResponse(w)
            return
        }

        problem := np.NewProblem(
            "https://example.com/problems/validation",
            "Validation Error",
            400,
            np.WithDetail("Email is required"),
            np.WithInstance(r.URL.Path),
            np.WithExtra("field", "email"),
        )

        w.Header().Set("Content-Type", np.ContentTypeProblemJSON)
        w.WriteHeader(problem.Status)
        problem.WriteHTTPResponse(w)
    })

    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

</details>

## Testing

```bash
go test -v
```

## RFC Compliance

**RFC 9457/7807 features implemented:**

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
