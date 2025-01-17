package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// defaultServeHTTP is the default ServeHTTP function that receives the status and error from
// the function call.
var defaultServeHTTP = func(w http.ResponseWriter, r *http.Request, err error) {
	if err != nil {
		switch e := err.(type) {
		case Error:
			http.Error(w, e.Error(), e.Status())
		default:
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
		}
	}
}

func TestParams(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	outParam := ""
	mux.Get("/user/{name}", func(w http.ResponseWriter, r *http.Request) (err error) {
		outParam = mux.Param(r, "name")
		return nil
	})

	r := httptest.NewRequest("GET", "/user/john", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "john", outParam)
}

func TestInstance(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	outParam := ""
	mux.Get("/user/{name}", func(w http.ResponseWriter, r *http.Request) (err error) {
		outParam = mux.Param(r, "name")
		return nil
	})

	r := httptest.NewRequest("GET", "/user/john", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "john", outParam)
}

func TestPostForm(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	form := url.Values{}
	form.Add("username", "jsmith")

	outParam := ""
	mux.Post("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		r.ParseForm()
		outParam = r.FormValue("username")
		return nil
	})

	r := httptest.NewRequest("POST", "/user", strings.NewReader(form.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "jsmith", outParam)
}

func TestPostJSON(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	j, err := json.Marshal(map[string]interface{}{
		"username": "jsmith",
	})
	assert.Nil(t, err)

	outParam := ""
	mux.Post("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		b, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		r.Body.Close()
		outParam = string(b)
		assert.Equal(t, `{"username":"jsmith"}`, string(b))
		return nil
	})

	r := httptest.NewRequest("POST", "/user", bytes.NewBuffer(j))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, `{"username":"jsmith"}`, outParam)
}

func TestGet(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Get("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestDelete(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Delete("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("DELETE", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestHead(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Head("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("HEAD", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestOptions(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Options("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("OPTIONS", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestPatch(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Patch("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("PATCH", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func TestPut(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Put("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("PUT", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
}

func Test404(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Get("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("GET", "/badroute", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, false, called)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func Test500NoError(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := true

	mux.Get("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return StatusError{Code: http.StatusInternalServerError, Err: nil}
	})

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func Test500WithError(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := true
	specificError := errors.New("specific error")

	mux.Get("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return StatusError{Code: http.StatusInternalServerError, Err: specificError}
	})

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, true, called)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, w.Body.String(), specificError.Error()+"\n")
}

func Test400(t *testing.T) {
	notFound := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	},
	)

	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)
	mux.SetNotFound(notFound)

	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestNotFound(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	mux.Error(http.StatusNotFound, w, r)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestBadRequest(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	r := httptest.NewRequest("GET", "/unknown", nil)
	w := httptest.NewRecorder()
	mux.Error(http.StatusBadRequest, w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestClear(t *testing.T) {
	mux := New()
	mux.SetServeHTTP(defaultServeHTTP)

	called := false

	mux.Get("/user", func(w http.ResponseWriter, r *http.Request) (err error) {
		called = true
		return nil
	})

	r := httptest.NewRequest("GET", "/user", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.True(t, called)

	called = false
	mux.Clear("GET", "/user")

	r = httptest.NewRequest("GET", "/user", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.False(t, called)
}
