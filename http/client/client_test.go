package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	netUrl "net/url"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rzajac/zltest"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	hc := &http.Client{Timeout: 42 * time.Second}
	acc := "application/xml"
	typ := "application/xml"
	auth := "token"
	url := "https://example.com"
	headers := map[string]string{
		"X-Test1": "1",
		"X-Test2": "2",
	}

	c := New(
		WithHTTPClient(hc),
		WithBaseURL(url),
		WithAuthBearer(auth),
		WithAccept(acc),
		WithContentType(typ),
		WithHeaders(headers),
	)

	assert.Equal(t, hc, c.opts.client)
	assert.Equal(t, acc, *c.opts.accept)
	assert.Equal(t, typ, *c.opts.contentType)
	assert.Equal(t, fmt.Sprintf("Bearer %s", auth), *c.opts.auth)
	assert.Equal(t, url, *c.opts.baseURL)
	assert.Equal(t, headers, *c.opts.headers)
}

func TestClient_NewRequest(t *testing.T) {
	url := "https://example.com"
	payload := []byte("body")

	c := New()
	r, err := c.NewRequest(context.Background(), http.MethodPost, url, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	body, err := io.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodPost, r.Method)
	assert.Equal(t, url, r.URL.String())
	assert.Equal(t, payload, body)
}

func before(req *http.Request) error {
	url, _ := netUrl.Parse("https://example.com/1")
	req.URL = url

	return nil
}

func TestClient_NewRequest_BeforeRequest(t *testing.T) {
	url := "https://example.com"

	c := New(WithBeforeRequest(before))
	r, err := c.NewRequest(context.Background(), http.MethodGet, url, nil)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, r.Method)
	assert.Equal(t, "https://example.com/1", r.URL.String())
}

func beforeChain(req *http.Request) error {
	url, _ := netUrl.Parse("https://example.com/2")
	req.URL = url

	return nil
}

func TestClient_NewRequest_BeforeRequest_Chain(t *testing.T) {
	url := "https://example.com"

	c := New(WithBeforeRequest(before), WithBeforeRequest(beforeChain))
	r, err := c.NewRequest(context.Background(), http.MethodGet, url, nil)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, r.Method)
	assert.Equal(t, "https://example.com/2", r.URL.String())
}

func beforeErr(_ *http.Request) error {
	return fmt.Errorf("some error")
}

func TestClient_NewRequest_BeforeRequest_Error(t *testing.T) {
	url := "https://example.com"

	c := New(WithBeforeRequest(beforeErr))
	_, err := c.NewRequest(context.Background(), http.MethodGet, url, nil)

	assert.Error(t, err)
}

func TestClient_NewRequest_NilContext(t *testing.T) {
	url := "https://example.com"
	c := New()
	_, err := c.NewRequest(nil, http.MethodGet, url, nil)

	assert.Error(t, err)
}

func TestClient_NewRequest_Error(t *testing.T) {
	c := New()
	_, err := c.NewRequest(context.Background(), ";", "", nil)

	assert.Error(t, err)
}

func TestClient_NewRequestWithOptions(t *testing.T) {
	hc := &http.Client{Timeout: 42 * time.Second}
	typ := "application/json"
	auth := "token"
	url := "https://example.com"
	headers := map[string]string{
		"X-Test1": "1",
		"X-Test2": "2",
	}
	endpoint := "/endpoint"
	payload := []byte("body")

	c := New(
		WithHTTPClient(hc),
		WithBaseURL(url),
		WithAuthBearer(auth),
		WithAuth(auth), // should overwrite the previous call
		WithContentType(typ),
		WithHeaders(headers),
	)

	r, err := c.NewRequest(context.Background(), http.MethodPost, endpoint, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	body, err := io.ReadAll(r.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodPost, r.Method)
	assert.Equal(t, url+endpoint, r.URL.String())
	assert.Equal(t, payload, body)
	assert.Equal(t, typ, r.Header.Get("Content-Type"))
	assert.Equal(t, auth, r.Header.Get("Authorization"))
	assert.Equal(t, headers["X-Test1"], r.Header.Get("X-Test1"))
	assert.Equal(t, headers["X-Test2"], r.Header.Get("X-Test2"))
}

func Server(statusCode int, body []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write(body)
	}))
}

func TestClient_Get(t *testing.T) {
	server := Server(http.StatusOK, nil)
	defer server.Close()

	c := New()
	resp, err := c.Get(context.Background(), server.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
}

type Hello struct {
	Hello string `json:"hello"`
}

func TestClient_Get_JSON(t *testing.T) {
	body := &Hello{Hello: "world"}
	b, err := json.Marshal(body)
	assert.NoError(t, err)
	server := Server(http.StatusOK, b)
	defer server.Close()

	c := New()
	resp, err := c.Get(context.Background(), server.URL)
	assert.NoError(t, err)

	hello := &Hello{}
	err = resp.UnmarshalJSON(hello)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	assert.Equal(t, body, hello)
}

func TestClient_Get_Error(t *testing.T) {
	c := New()
	_, err := c.Get(context.Background(), "\n")

	assert.Error(t, err)
}

func after(resp *Response) error {
	resp.StatusCode = 201
	return nil
}

func TestClient_Get_AfterRequest(t *testing.T) {
	server := Server(http.StatusOK, nil)
	defer server.Close()

	c := New(WithAfterRequest(after))
	resp, err := c.Get(context.Background(), server.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func afterChain(resp *Response) error {
	resp.StatusCode = http.StatusTeapot
	return nil
}

func TestClient_Get_AfterRequest_Chain(t *testing.T) {
	server := Server(http.StatusOK, nil)
	defer server.Close()

	c := New(WithAfterRequest(after), WithAfterRequest(afterChain))
	resp, err := c.Get(context.Background(), server.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	assert.Equal(t, http.StatusTeapot, resp.StatusCode)
}

func afterErr(_ *Response) error {
	return fmt.Errorf("some error")
}

func TestClient_Get_AfterRequest_Error(t *testing.T) {
	server := Server(http.StatusOK, nil)
	defer server.Close()

	c := New(WithAfterRequest(afterErr))
	_, err := c.Get(context.Background(), server.URL)
	assert.Error(t, err)
}

func TestClient_Post(t *testing.T) {
	payload := []byte("body")
	server := Server(http.StatusOK, payload)
	defer server.Close()

	c := New()
	resp, err := c.Post(context.Background(), server.URL, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodPost, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	assert.Equal(t, payload, body)
}

func TestClient_Post_Error(t *testing.T) {
	c := New()
	_, err := c.Post(context.Background(), "\n", nil)

	assert.Error(t, err)
}

func TestClient_Put(t *testing.T) {
	payload := []byte("body")
	server := Server(http.StatusOK, payload)
	defer server.Close()

	c := New()
	resp, err := c.Put(context.Background(), server.URL, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodPut, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	assert.Equal(t, payload, body)
}

func TestClient_Put_Error(t *testing.T) {
	c := New()
	_, err := c.Put(context.Background(), "\n", nil)

	assert.Error(t, err)
}

func TestClient_Patch(t *testing.T) {
	payload := []byte("body")
	server := Server(http.StatusOK, payload)
	defer server.Close()

	c := New()
	resp, err := c.Patch(context.Background(), server.URL, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodPatch, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	assert.Equal(t, payload, body)
}

func TestClient_Patch_Error(t *testing.T) {
	c := New()
	_, err := c.Patch(context.Background(), "\n", nil)

	assert.Error(t, err)
}

func TestClient_Delete(t *testing.T) {
	server := Server(http.StatusOK, nil)
	defer server.Close()

	c := New()
	resp, err := c.Delete(context.Background(), server.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodDelete, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
}

func TestClient_Delete_Error(t *testing.T) {
	c := New()
	_, err := c.Delete(context.Background(), "\n")

	assert.Error(t, err)
}

func TestClient_OverrideURL(t *testing.T) {
	override := "https://example.com/override"
	c := New()
	c.OverrideURL(override)
	req, err := c.NewRequest(context.Background(), http.MethodGet, "https://example.com", nil)
	assert.NoError(t, err)

	assert.Equal(t, override, req.URL.String())
}

func TestDefaultClient(t *testing.T) {
	c := DefaultHTTPClient()

	assert.Equal(t, c.Timeout, defaultTimeout)
}

func TestDefaultClientTransport(t *testing.T) {
	c := DefaultHTTPTransport()

	assert.Equal(t, c.TLSHandshakeTimeout, defaultTLSHandshakeTimeout)
}

func TestBuildURLWithQuery(t *testing.T) {
	url := "https://example.com"
	qs := map[string]string{
		"k":      "v",
		"answer": "42",
	}

	s, err := BuildURLWithQuery(url, qs)
	if err != nil {
		assert.NoError(t, err)
	}

	assert.Equal(t, "https://example.com?answer=42&k=v", s)
}

func TestBuildURLWithQuery_Error(t *testing.T) {
	_, err := BuildURLWithQuery("\n", nil)

	assert.Error(t, err)
}

func TestClient_Do_Error(t *testing.T) {
	c := New()
	_, err := c.Do(&http.Request{})

	assert.Error(t, err)
}

func TestClient_WithFormatErrors(t *testing.T) {
	response := []byte("bad request")
	server := Server(http.StatusBadRequest, response)
	defer server.Close()

	c := New(WithFormatErrors(true))
	_, err := c.Get(context.Background(), server.URL)
	assert.Error(t, err)

	m := fmt.Sprintf("http.client: %d response for %s: %s", http.StatusBadRequest, server.URL, response)
	assert.Equal(t, m, err.Error())
}

func TestClient_WithFormatErrors_404(t *testing.T) {
	server := Server(http.StatusNotFound, nil)
	defer server.Close()

	c := New(WithFormatErrors(true))
	_, err := c.Get(context.Background(), server.URL)
	assert.Error(t, err)

	assert.Equal(t, ErrNotFound, err)
}

func TestClient_Get_WithDebugRequest(t *testing.T) {
	zl := zltest.New(t)
	log := zerolog.New(zl).With().Timestamp().Logger()

	server := Server(http.StatusOK, nil)
	defer server.Close()

	c := New(WithDebugRequest(log))
	resp, err := c.Get(context.Background(), server.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	le := zl.LastEntry()
	le.ExpStr("method", http.MethodGet)
	le.ExpStr("url", server.URL)
	le.ExpStr("body", "")
}

func TestClient_Post_WithDebugRequest(t *testing.T) {
	zl := zltest.New(t)
	log := zerolog.New(zl).With().Timestamp().Logger()

	payload := []byte("body")
	server := Server(http.StatusOK, payload)
	defer server.Close()

	c := New(WithDebugRequest(log))
	resp, err := c.Post(context.Background(), server.URL, bytes.NewBuffer(payload))
	assert.NoError(t, err)

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodPost, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	assert.Equal(t, payload, body)
	le := zl.LastEntry()
	le.ExpStr("method", http.MethodPost)
	le.ExpStr("url", server.URL)
	le.ExpStr("body", string(payload))
}

func TestClient_WithDebugResponse(t *testing.T) {
	zl := zltest.New(t)
	log := zerolog.New(zl).With().Timestamp().Logger()

	response := []byte("response")
	server := Server(http.StatusOK, response)
	defer server.Close()

	c := New(WithDebugResponse(log))
	resp, err := c.Get(context.Background(), server.URL)
	assert.NoError(t, err)

	assert.Equal(t, http.MethodGet, resp.Request.Method)
	assert.Equal(t, server.URL, resp.Request.URL.String())
	le := zl.LastEntry()
	le.ExpNum("status_code", http.StatusOK)
	le.ExpStr("method", http.MethodGet)
	le.ExpStr("url", server.URL)
	le.ExpStr("body", string(response))
}
