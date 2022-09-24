package main

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/hexops/autogold"

	"github.com/larschri/pestotrap/dirindex"
	"github.com/larschri/searchpage"
)

//go:embed testdata/jsons
var testdata embed.FS

var handler *searchpage.Handler

func TestMain(m *testing.M) {
	jsons, err := fs.Sub(testdata, "testdata/jsons")
	if err != nil {
		panic(err)
	}

	idx, err := dirindex.NewMemOnly()
	if err != nil {
		panic(err)
	}

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	if err := dirindex.Start(jsons, idx, ticker.C); err != nil {
		panic(err)
	}

	handler = searchpage.New(&config, idx)
	os.Exit(m.Run())
}

func withHxRequestHeader(r *http.Request) *http.Request {
	r.Header["Hx-Request"] = []string{"true"}
	return r
}

func TestRequests(t *testing.T) {
	tests := map[string]*http.Request{
		"root":    httptest.NewRequest(http.MethodGet, "/", nil),
		"search":  httptest.NewRequest(http.MethodGet, "/q?search=anotherfield1", nil),
		"search2": withHxRequestHeader(httptest.NewRequest(http.MethodGet, "/q?search=anotherfield1", nil)),
		"invalid": withHxRequestHeader(httptest.NewRequest(http.MethodGet, "/q?search=xx:%20", nil)),
	}

	for name, req := range tests {
		t.Run(name, func(t *testing.T) {

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			bs, err := io.ReadAll(w.Result().Body)
			if err != nil {
				t.Fatal(err)
			}

			autogold.Equal(t, map[string]interface{}{
				"body":    string(bs),
				"status":  w.Code,
				"headers": w.Header(),
			})
		})
	}
}
