package middleware

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/alexferl/golib/http/api/handler"
	"github.com/alexferl/golib/http/api/server"
)

//go:embed static/*
var staticFS embed.FS

func TestCache(t *testing.T) {
	s := server.New()
	s.GET("/static/*", handler.Static("/static/", staticFS, "static"), Cache("/static/", time.Minute*1))

	testCases := []struct {
		name string
		req  *http.Request
		resp *httptest.ResponseRecorder
	}{
		{"/foobar", httptest.NewRequest(http.MethodGet, "/static/foobar", nil), httptest.NewRecorder()},
		{"/foo/bar", httptest.NewRequest(http.MethodGet, "/static/foo/bar", nil), httptest.NewRecorder()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s.ServeHTTP(tc.resp, tc.req)
			assert.Equal(t, http.StatusOK, tc.resp.Code)
			assert.Contains(t, tc.resp.Header().Get("Cache-Control"), "public, max-age=60")
		})
	}
}

func TestCache_Sub(t *testing.T) {
	s := server.New()
	s.GET("/static/*", handler.Static("/static/", staticFS, "static"), Cache("/static/foo", time.Minute*1))

	testCases := []struct {
		name   string
		req    *http.Request
		resp   *httptest.ResponseRecorder
		header string
	}{
		{"/foobar", httptest.NewRequest(http.MethodGet, "/static/foobar", nil), httptest.NewRecorder(), ""},
		{"/foo/bar", httptest.NewRequest(http.MethodGet, "/static/foo/bar", nil), httptest.NewRecorder(), "public, max-age=60"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s.ServeHTTP(tc.resp, tc.req)
			assert.Equal(t, http.StatusOK, tc.resp.Code)
			assert.Contains(t, tc.resp.Header().Get("Cache-Control"), tc.header)
		})
	}
}
