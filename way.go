// Package away is a fork of way: https://github.com/matryer/way
package away

import (
	"context"
	"net/http"
	"sort"
	"strings"
)

// wayContextKey is the context key type for storing
// parameters in context.Context.
type wayContextKey string

// Router routes HTTP requests.
type Router struct {
	routes routeList
	// NotFound is the http.Handler to call when no routes
	// match. By default uses http.NotFoundHandler().
	NotFound http.Handler
}

// NewRouter makes a new Router.
func NewRouter() *Router {
	return &Router{
		NotFound: http.NotFoundHandler(),
	}
}

func (r *Router) pathSegments(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

// Remove an entry from the router.
func (r *Router) Remove(method string, p string) {
	for index, v := range r.routes {
		if v.pattern == p && strings.EqualFold(v.method, method) {
			r.routes = removeIndex(r.routes, index)
		}
	}
}

// Count returns the number of routes.
func (r *Router) Count() int {
	return len(r.routes)
}

func removeIndex(s []*route, index int) []*route {
	return append(s[:index], s[index+1:]...)
}

// Handle adds a handler with the specified method and pattern.
// Method can be any HTTP method string or "*" to match all methods.
// Pattern can contain path segments such as: /item/:id which is
// accessible via the Param function.
// If pattern ends with trailing /, it acts as a prefix.
func (r *Router) Handle(method, pattern string, handler http.Handler) {
	route := &route{
		pattern: pattern,
		method:  strings.ToLower(method),
		segs:    r.pathSegments(pattern),
		handler: handler,
		prefix:  strings.HasSuffix(pattern, "/") || strings.HasSuffix(pattern, "..."),
	}
	r.routes = append(r.routes, route)

	// Sort so the routes are in the proper order.
	sort.Sort(r.routes)
}

// HandleFunc is the http.HandlerFunc alternative to http.Handle.
func (r *Router) HandleFunc(method, pattern string, fn http.HandlerFunc) {
	r.Handle(method, pattern, fn)
}

// ServeHTTP routes the incoming http.Request based on method and path
// extracting path parameters as it goes.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := strings.ToLower(req.Method)
	segs := r.pathSegments(req.URL.Path)
	for _, route := range r.routes {
		if route.method != method && route.method != "*" {
			continue
		}
		if ctx, ok := route.match(req.Context(), r, segs); ok {
			route.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	r.NotFound.ServeHTTP(w, req)
}

// Param gets the path parameter from the specified Context.
// Returns an empty string if the parameter was not found.
func Param(ctx context.Context, param string) string {
	vStr, ok := ctx.Value(wayContextKey(param)).(string)
	if !ok {
		return ""
	}
	return vStr
}

type route struct {
	pattern string
	method  string
	segs    []string
	handler http.Handler
	prefix  bool
}

type routeList []*route

func (s routeList) Len() int {
	return len(s)
}

func (s routeList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s routeList) Less(i, j int) bool {
	var si string = s[i].pattern
	var sj string = s[j].pattern
	var siLower = strings.ToLower(si)
	var sjLower = strings.ToLower(sj)
	if strings.HasPrefix(sjLower, "/:") {
		return true
	} else if strings.HasPrefix(siLower, "/:") {
		return false
	} else if strings.Contains(sjLower, ":") && !strings.Contains(siLower, ":") {
		return true
	} else if !strings.Contains(sjLower, ":") && strings.Contains(siLower, ":") {
		return false
	}

	if siLower == sjLower {
		return si < sj
	}
	return siLower < sjLower
}

func (r *route) match(ctx context.Context, router *Router, segs []string) (context.Context, bool) {
	if len(segs) > len(r.segs) && !r.prefix {
		return nil, false
	}
	for i, seg := range r.segs {
		if i > len(segs)-1 {
			return nil, false
		}
		isParam := false
		if strings.HasPrefix(seg, ":") {
			isParam = true
			seg = strings.TrimPrefix(seg, ":")
		}
		if !isParam { // verbatim check
			if strings.HasSuffix(seg, "...") {
				if strings.HasPrefix(segs[i], seg[:len(seg)-3]) {
					return ctx, true
				}
			}
			if seg != segs[i] {
				return nil, false
			}
		}
		if isParam {
			ctx = context.WithValue(ctx, wayContextKey(seg), segs[i])
		}
	}
	return ctx, true
}
