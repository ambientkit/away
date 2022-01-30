package away_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/ambientkit/away"
	"github.com/stretchr/testify/assert"
)

var tests = []struct {
	RouteMethod  string
	RoutePattern string

	Method string
	Path   string
	Match  bool
	Params map[string]string
}{
	// simple path matching
	{
		"GET", "/one",
		"GET", "/one", true, nil,
	},
	{
		"GET", "/two",
		"GET", "/two", true, nil,
	},
	{
		"GET", "/three",
		"GET", "/three", true, nil,
	},
	// methods
	{
		"get", "/methodcase",
		"GET", "/methodcase", true, nil,
	},
	{
		"Get", "/methodcase",
		"get", "/methodcase", true, nil,
	},
	{
		"GET", "/methodcase",
		"get", "/methodcase", true, nil,
	},
	{
		"GET", "/method1",
		"POST", "/method1", false, nil,
	},
	{
		"DELETE", "/method2",
		"GET", "/method2", false, nil,
	},
	{
		"GET", "/method3",
		"PUT", "/method3", false, nil,
	},
	// all methods
	{
		"*", "/all-methods",
		"GET", "/all-methods", true, nil,
	},
	{
		"*", "/all-methods",
		"POST", "/all-methods", true, nil,
	},
	{
		"*", "/all-methods",
		"PUT", "/all-methods", true, nil,
	},
	// nested
	{
		"GET", "/parent/child/one",
		"GET", "/parent/child/one", true, nil,
	},
	{
		"GET", "/parent/child/two",
		"GET", "/parent/child/two", true, nil,
	},
	{
		"GET", "/parent/child/three",
		"GET", "/parent/child/three", true, nil,
	},
	// slashes
	{
		"GET", "slashes/one",
		"GET", "/slashes/one", true, nil,
	},
	{
		"GET", "/slashes/two",
		"GET", "slashes/two", true, nil,
	},
	{
		"GET", "slashes/three/",
		"GET", "/slashes/three", true, nil,
	},
	{
		"GET", "/slashes/four",
		"GET", "slashes/four/", true, nil,
	},
	// prefix
	{
		"GET", "/prefix/",
		"GET", "/prefix/anything/else", true, nil,
	},
	{
		"GET", "/not-prefix",
		"GET", "/not-prefix/anything/else", false, nil,
	},
	{
		"GET", "/prefixdots...",
		"GET", "/prefixdots/anything/else", true, nil,
	},
	{
		"GET", "/prefixdots...",
		"GET", "/prefixdots", true, nil,
	},
	// path params
	{
		"GET", "/path-param/:id",
		"GET", "/path-param/123", true, map[string]string{"id": "123"},
	},
	{
		"GET", "/path-params/:era/:group/:member",
		"GET", "/path-params/60s/beatles/lennon", true, map[string]string{
			"era":    "60s",
			"group":  "beatles",
			"member": "lennon",
		},
	},
	{
		"GET", "/path-params-prefix/:era/:group/:member/",
		"GET", "/path-params-prefix/60s/beatles/lennon/yoko", true, map[string]string{
			"era":    "60s",
			"group":  "beatles",
			"member": "lennon",
		},
	},
	// misc no matches
	{
		"GET", "/not/enough",
		"GET", "/not/enough/items", false, nil,
	},
	{
		"GET", "/not/enough/items",
		"GET", "/not/enough", false, nil,
	},
}

func TestWay(t *testing.T) {
	for _, test := range tests {
		r := away.NewRouter()
		match := false
		var ctx context.Context
		r.Handle(test.RouteMethod, test.RoutePattern, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			match = true
			ctx = r.Context()
		}))
		req, err := http.NewRequest(test.Method, test.Path, nil)
		if err != nil {
			t.Errorf("NewRequest: %s", err)
		}
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if match != test.Match {
			t.Errorf("expected match %v but was %v: %s %s", test.Match, match, test.Method, test.Path)
		}
		if len(test.Params) > 0 {
			for expK, expV := range test.Params {
				// check using helper
				actualValStr := away.Param(ctx, expK)
				if actualValStr != expV {
					t.Errorf("Param: context value %s expected \"%s\" but was \"%s\"", expK, expV, actualValStr)
				}
			}
		}
	}
}

func TestMultipleRoutesDifferentMethods(t *testing.T) {
	r := away.NewRouter()
	var match string
	r.Handle(http.MethodGet, "/route", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		match = "GET /route"
	}))
	r.Handle(http.MethodDelete, "/route", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		match = "DELETE /route"
	}))
	r.Handle(http.MethodPost, "/route", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		match = "POST /route"
	}))

	req, err := http.NewRequest(http.MethodGet, "/route", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	r.ServeHTTP(httptest.NewRecorder(), req)
	if match != "GET /route" {
		t.Errorf("unexpected: %s", match)
	}

	req, err = http.NewRequest(http.MethodDelete, "/route", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	r.ServeHTTP(httptest.NewRecorder(), req)
	if match != "DELETE /route" {
		t.Errorf("unexpected: %s", match)
	}

	req, err = http.NewRequest(http.MethodPost, "/route", nil)
	if err != nil {
		t.Errorf("NewRequest: %s", err)
	}
	r.ServeHTTP(httptest.NewRecorder(), req)
	if match != "POST /route" {
		t.Errorf("unexpected: %s", match)
	}

}

type route struct {
	pattern string
	method  string
	segs    []string
	handler http.Handler
	prefix  bool
}

type routeList []route

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

func TestAway(t *testing.T) {
	var arr routeList = []route{
		{
			method:  "GET",
			pattern: "/:slug",
		},
		{
			method:  "GET",
			pattern: "/",
		},
		{
			method:  "GET",
			pattern: "/cool/balloon",
		},
		{
			method:  "GET",
			pattern: "/cool/:slug",
		},
		{
			method:  "GET",
			pattern: "/cool/another",
		},
		{
			method:  "GET",
			pattern: "/rss/:ok",
		},
		{
			method:  "GET",
			pattern: "/cool",
		},
		{
			method:  "GET",
			pattern: "/rss.xml",
		},
	}

	assert.Equal(t, arr[0].pattern, "/:slug")
	sort.Sort(arr)
	assert.Equal(t, arr[0].pattern, "/")
	assert.Equal(t, arr[len(arr)-1].pattern, "/:slug")
}
