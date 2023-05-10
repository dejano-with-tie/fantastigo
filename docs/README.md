[//]: # (@formatter:off)
# Content

1. [Project Structure](#project-structure)
2. [Http](#http)

## Project Structure

Project layout is inspired and follows go [community recommendations](https://github.com/golang-standards/project-layout).

### [`/cmd`](../cmd)

Main applications for this project. It consists of two executable applications, hence two subdirectories can be found
within cmd dir. `/cmd/ector` and `/cmd/fleet`.

### [`/config`](../config)

Configuration files.

### [`/docker`](../docker)

TODO

### [`/docs`](../docs)

Contains documentation.

### [`/internal`](../internal)

Implementation code, not meant to be used by external apps.

### [`/pkg`](../pkg)

Code that can be used by external applications. Be careful what you place here. E.g. http client libs encourage so
called modular-monolith anti-pattern in microservice architecture.

### [`/scripts`](../scripts)

Scripts to perform various build, install, analysis, etc operations.

## Http

#### Echo Framework

#### Documentation, OpenApi & Swagger

#### Api-first Approach

#### Validation

#### Error handling

Go's error handling doesn't relay on traditional try/catch mechanism. Instead, errors are returned as normal values,
i.e. `error` type.
On the other hand, there is a `panic` built-in function which stops the program, and can be recovered from with
the `recover` built-in function. Panic should be used sparingly. E.g. when creation of necessary dependencies fail and
the program is unable to start.

> Echo framework provides [recoverer](https://echo.labstack.com/middleware/recover/) middleware which prevents http
> server from shutting down if `panic` occurs within request.

Go's errors can be formed in a couple of ways, can be expanded with additional context and _wrapped_ when
needed. [Good blog post](https://go.dev/blog/go1.13-errors) on the subject.

```go
fmt.Errorf("error reading config file <path=%s>", path) // allows formatting
errors.New("should respond with status=500, code=unknown") // only plain text
```

Function `fmt.Errorf` is especially interesting because it can be used to add additional context to the existing error,
or _wrap_ it.

```go
var ErrAccessDenied = errors.New("Access denied") // sentinel error

// somewhere else, expand the sentinel err
return fmt.Errorf("User with <id=%s> error: %v", id, ErrAccessDenied) // add additional context with '%v'
return fmt.Errorf("User with <id=%s> error: %w", id, ErrAccessDenied) // _wrap_ the error with '%w'

```

We can also create custom error types, enabling us to perform checks and react. Custom error types are defined by
implementing the `Error` interface:

```go
type NotFound struct { Id string }
func (n *NotFound) Error() string {
    return fmt.Sprintf("Entity not found <id=%s>", n.Id)
}
```

---
Echo's [error handling](https://echo.labstack.com/guide/error-handling/) is centralized and boils down to
a [single function](https://github.com/labstack/echo/blob/v4.10.2/echo.go#L418) which is essentially responsible for
writing the error response.

Echo's error response model has only a single property `message`, which unfortunately doesn't necessarily provide enough
context. This project implements
a [custom error](../internal/common/server/errhandler.go) handler with enriched response model.

```go
type (
    HttpErrResponse struct {
        Status      int                       `json:"-"`
        Code        string                    `json:"code"`
        Validations []ValidationFieldResponse `json:"validations,omitempty"`
        Message     string                    `json:"message,omitempty"`
        Trace       string                    `json:"trace,omitempty"`
    }
    ValidationFieldResponse struct {
        Property      string `json:"property"`
        Key           string `json:"error"`
        RejectedValue string `json:"rejectedValue"`
        Message       string `json:"message,omitempty"`
    }
)
```

The base idea is to have a transport-agnostic [custom Error type (AppErr)](../internal/common/apperr/error.go) which
explains the
error reason with the machine-readable error `code` property.
Attributes specific to the transport (http) protocol are later on calculated based on `code` value in the err handler
implementation, and all together is written back to the http client.

```json
{
  "code": "bad-request",
  "validations": [
    {
      "property": "Capacity",
      "error": "required",
      "rejectedValue": "0",
      "message": "Property validation for 'Capacity' failed validation tag 'required' tag"
    }
  ],
  "message": "Request validation failed"
}
```

Implementation of a handler boils down to the error type checking and mapping the error to the `HttpErrResponse`.

```go
// default
status := getHttpStatus(apperr.ErrCodeUnknown)
code := apperr.ErrCodeUnknown
message := err.Error()
var validations []ValidationFieldResponse

var appErr *apperr.AppErr
var echoHttpErr *echo.HTTPError
var validationErrs validator.ValidationErrors
if errors.As(err, &validationErrs) {
//...
} else if errors.As(err, &appErr) {
//...
} else if errors.As(err, &echoHttpErr) {
//...
}

```

#### Versioning

Api versioning boils down to the general practice, it is agnostic of the technology used. Therefore, there are a couple
of options out there, most common ones are Api versioning via Http _headers_ or via _url_. There are articles on the web
going in dept on the subject. Personally, I prefer _url_ approach for simplicity, and _it just works_.