// Package router provides request handling capabilities.
package router

import (
	"net/http"

	"github.com/ambientkit/away"
	"github.com/ambientkit/away/router/paramconvert"
)

// Mux contains the router.
type Mux struct {
	router *away.Router

	// customServeHTTP is the serve function.
	customServeHTTP func(w http.ResponseWriter, r *http.Request, err error)
}

// New returns an instance of the router.
func New() *Mux {
	r := away.NewRouter()

	return &Mux{
		router: r,
	}
}

// SetServeHTTP sets the ServeHTTP function.
func (m *Mux) SetServeHTTP(csh func(w http.ResponseWriter, r *http.Request, err error)) {
	m.customServeHTTP = csh
}

// SetNotFound sets the NotFound function.
func (m *Mux) SetNotFound(notFound http.Handler) {
	m.router.NotFound = notFound
}

// Clear will remove a method and path from the router.
func (m *Mux) Clear(method string, path string) {
	m.router.Remove(method, paramconvert.BraceToColon(path))
}

// Count will return the number of routes from the router.
func (m *Mux) Count() int {
	return m.router.Count()
}

// ServeHTTP routes the incoming http.Request based on method and path
// extracting path parameters as it goes.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

// StatusError returns error with a status code.
func (m *Mux) StatusError(status int, err error) error {
	return StatusError{Code: status, Err: err}
}

// Error shows error page based on the status code.
func (m *Mux) Error(status int, w http.ResponseWriter, r *http.Request) {
	if m.customServeHTTP != nil {
		m.customServeHTTP(w, r, StatusError{Code: status, Err: nil})
		return
	}

	http.Error(w, http.StatusText(status), status)
}

// Param returns a URL parameter.
func (m *Mux) Param(r *http.Request, param string) string {
	return away.Param(r.Context(), param)
}

// Wrap a standard http handler so it can be used easily.
func (m *Mux) Wrap(handler http.HandlerFunc) func(w http.ResponseWriter, r *http.Request) (err error) {
	return func(w http.ResponseWriter, r *http.Request) (err error) {
		handler.ServeHTTP(w, r)
		return
	}
}

// Error represents a handler error. It provides methods for a HTTP status
// code and embeds the built-in error interface.
type Error interface {
	error
	Status() int
	Message() string
}

// StatusError represents an error with an associated HTTP status code.
type StatusError struct {
	Code     int
	Err      error
	Friendly string
}

// Error returns the error.
func (se StatusError) Error() string {
	if se.Err != nil {
		return se.Err.Error()
	}

	return ""
}

// Status returns a HTTP status code.
func (se StatusError) Status() int {
	return se.Code
}

// Message returns a optional user friendly error message.
func (se StatusError) Message() string {
	return se.Friendly
}
