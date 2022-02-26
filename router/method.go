package router

import (
	"net/http"

	"github.com/ambientkit/away/router/ambhandler"
	"github.com/ambientkit/away/router/paramconvert"
)

func (m *Mux) handle(method string, path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.router.Handle(method, paramconvert.BraceToColon(path), ambhandler.Handler{
		HandlerFunc:     fn,
		CustomServeHTTP: m.customServeHTTP,
	})
}

// Delete registers a pattern with the router.
func (m *Mux) Delete(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(http.MethodDelete, path, fn)
}

// Get registers a pattern with the router.
func (m *Mux) Get(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(http.MethodGet, path, fn)
}

// Handle registers a method and pattern with the router.
func (m *Mux) Handle(method string, path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(method, path, fn)
}

// Head registers a pattern with the router.
func (m *Mux) Head(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(http.MethodHead, path, fn)
}

// Options registers a pattern with the router.
func (m *Mux) Options(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(http.MethodOptions, path, fn)
}

// Patch registers a pattern with the router.
func (m *Mux) Patch(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(http.MethodPatch, path, fn)
}

// Post registers a pattern with the router.
func (m *Mux) Post(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(http.MethodPost, path, fn)
}

// Put registers a pattern with the router.
func (m *Mux) Put(path string, fn func(http.ResponseWriter, *http.Request) (int, error)) {
	m.handle(http.MethodPut, path, fn)
}
