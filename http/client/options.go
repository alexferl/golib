package client

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
)

type options struct {
	client        *http.Client
	accept        *string
	contentType   *string
	auth          *string
	baseURL       *string
	headers       *map[string]string
	beforeRequest []func(req *http.Request, v ...any) error
	afterRequest  []func(resp *Response, v ...any) error
	fmtError      bool
	logRequest    bool
	logResponse   bool
}

var (
	accept      = "application/json, */*;q=0.5"
	contentType = "application/json"
	logger      = zerolog.Logger{}
)

var ErrNotFound = errors.New("not found")

var defaultOptions = options{
	client:      DefaultHTTPClient(),
	accept:      &accept,
	contentType: &contentType,
	logRequest:  false,
	logResponse: false,
}

type Option interface {
	apply(*options)
}

// funcOption wraps a function that modifies options into an
// implementation of the Option interface.
type funcOption struct {
	f func(*options)
}

func (fco *funcOption) apply(do *options) {
	fco.f(do)
}

func newFuncOption(f func(*options)) *funcOption {
	return &funcOption{
		f: f,
	}
}

// WithHTTPClient sets the http.Client.
func WithHTTPClient(client *http.Client) Option {
	return newFuncOption(func(o *options) {
		o.client = client
	})
}

// WithBaseURL sets a base URL that will be prepended to every
// subsequent URL when making requests.
func WithBaseURL(s string) Option {
	return newFuncOption(func(o *options) {
		o.baseURL = &s
	})
}

// WithAuth sets the Authorization header that will be sent with all requests.
func WithAuth(s string) Option {
	return newFuncOption(func(o *options) {
		o.auth = &s
	})
}

// WithAuthBearer sets the Authorization header with the Bearer prefix
// that will be sent with all requests.
func WithAuthBearer(s string) Option {
	return newFuncOption(func(o *options) {
		f := fmt.Sprintf("Bearer %s", s)
		o.auth = &f
	})
}

// WithAccept sets the Accept header that will be sent with
// all request.
func WithAccept(s string) Option {
	return newFuncOption(func(o *options) {
		o.accept = &s
	})
}

// WithContentType sets the ContentType header that will be sent with
// POST, PUT and PATCH requests.
func WithContentType(s string) Option {
	return newFuncOption(func(o *options) {
		o.contentType = &s
	})
}

// WithHeaders sets headers that will be sent with all requests.
func WithHeaders(m map[string]string) Option {
	return newFuncOption(func(o *options) {
		o.headers = &m
	})
}

// WithBeforeRequest sets a function that will run before every request
// and has access to the current request to inspect and/or modify it.
// Can be called multiple times to chain functions.
func WithBeforeRequest(f func(req *http.Request, v ...any) error) Option {
	return newFuncOption(func(o *options) {
		o.beforeRequest = append(o.beforeRequest, f)
	})
}

// WithAfterRequest sets a function that will run after every request
// and has access to the current response to inspect and/or modify it.
// Can be called multiple times to chain functions.
func WithAfterRequest(f func(resp *Response, v ...any) error) Option {
	return newFuncOption(func(o *options) {
		o.afterRequest = append(o.afterRequest, f)
	})
}

// WithFormatErrors will format the >= 400 responses and return them in the following format:
// 'http.client: <status_code> response for <url>: <body>'.
// 404 responses will return ErrNotFound instead as sometimes a 404 is expected.
func WithFormatErrors(b bool) Option {
	return newFuncOption(func(o *options) {
		if b && !o.fmtError {
			o.fmtError = true
			// prepend so it runs first
			var fn []func(resp *Response, v ...any) error
			fn = append(fn, fmtError)
			o.afterRequest = append(fn, o.afterRequest...)
		}
	})
}

// WithDebugRequest will log the request to the supplied zerolog.Logger.
func WithDebugRequest(l zerolog.Logger) Option {
	return newFuncOption(func(o *options) {
		if !o.logRequest {
			logger = l
			o.logRequest = true
			// prepend so it runs first
			var fn []func(req *http.Request, v ...any) error
			fn = append(fn, logRequest)
			o.beforeRequest = append(fn, o.beforeRequest...)
		}
	})
}

// WithDebugResponse will log the response to the supplied zerolog.Logger.
func WithDebugResponse(l zerolog.Logger) Option {
	return newFuncOption(func(o *options) {
		if !o.logResponse {
			logger = l
			o.logResponse = true
			// prepend so it runs first
			var fn []func(resp *Response, v ...any) error
			fn = append(fn, logResponse)
			o.afterRequest = append(fn, o.afterRequest...)
		}
	})
}

func fmtError(resp *Response, _ ...any) error {
	if resp.StatusCode >= http.StatusBadRequest {
		// some requests expect a 404
		if resp.StatusCode == http.StatusNotFound {
			return ErrNotFound
		}

		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		return fmt.Errorf("http.client: %d response for %s: %s", resp.StatusCode, resp.Request.URL, respBody)
	}

	return nil
}

func logRequest(req *http.Request, _ ...any) error {
	var body []byte
	if req.Body != nil {
		b, err := io.ReadAll(req.Body)
		if err != nil {
			return err
		}

		body = b
		req.Body = io.NopCloser(bytes.NewBuffer(b))
	}

	logger.Debug().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Str("body", string(body)).
		Msg("http.client request")

	return nil
}

func logResponse(resp *Response, _ ...any) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(body))

	logger.Debug().
		Int("status_code", resp.StatusCode).
		Str("method", resp.Request.Method).
		Str("url", resp.Request.URL.String()).
		Str("body", string(body)).
		Msg("http.client response")

	return nil
}
