package client

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
)

func TestWithHTTPClient(t *testing.T) {
	hc := &http.Client{Timeout: 42 * time.Second}
	c := New(WithHTTPClient(hc))

	assert.Equal(t, hc.Timeout, c.opts.client.Timeout)
}

func TestAccept(t *testing.T) {
	typ := accept
	c := New(WithAccept(typ))

	assert.Equal(t, typ, *c.opts.accept)
}

func TestWithContentType(t *testing.T) {
	typ := contentType
	c := New(WithContentType(typ))

	assert.Equal(t, typ, *c.opts.contentType)
}

func TestWithAuth(t *testing.T) {
	auth := "value"
	c := New(WithAuth(auth))

	assert.Equal(t, auth, *c.opts.auth)
}

func TestWithAuthBearer(t *testing.T) {
	auth := "value"
	c := New(WithAuthBearer(auth))

	assert.Equal(t, fmt.Sprintf("Bearer %s", auth), *c.opts.auth)
}

func TestWithBaseURL(t *testing.T) {
	url := "https://example.com"
	c := New(WithBaseURL(url))

	assert.Equal(t, url, *c.opts.baseURL)
}

func TestWithHeaders(t *testing.T) {
	m := map[string]string{
		"X-Test1": "1",
		"X-Test2": "2",
	}
	c := New(WithHeaders(m))

	assert.Equal(t, m, *c.opts.headers)
}

func TestWithBeforeRequest(t *testing.T) {
	f := func(req *http.Request, v ...any) error { return nil }
	var fns []func(req *http.Request, v ...any) error
	fns = append(fns, f)
	c := New(WithBeforeRequest(f))

	assert.ObjectsAreEqual(f, c.opts.beforeRequest)
}

func TestWithAfterRequest(t *testing.T) {
	f := func(resp *Response, v ...any) error { return nil }
	var fns []func(resp *Response, v ...any) error
	fns = append(fns, f)

	c := New(WithAfterRequest(f))

	assert.ObjectsAreEqual(f, c.opts.afterRequest)
}

func TestWithDebugRequest(t *testing.T) {
	c := New(WithDebugRequest(zerolog.Logger{}))

	assert.ObjectsAreEqual(logRequest, c.opts.beforeRequest[0])
}

func TestWithDebugResponse(t *testing.T) {
	c := New(WithDebugResponse(zerolog.Logger{}))

	assert.ObjectsAreEqual(logResponse, c.opts.afterRequest[0])
}
