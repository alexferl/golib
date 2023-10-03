package client

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	netUrl "net/url"
	"time"
)

const (
	defaultTimeout             = 60 * time.Second
	defaultDialTimeout         = 5 * time.Second
	defaultTLSHandshakeTimeout = 5 * time.Second
)

type Client struct {
	opts        options
	overrideURL *string
}

// New creates a new client.
func New(opt ...Option) *Client {
	opts := defaultOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	c := &Client{
		opts: opts,
	}

	return c
}

// NewRequest creates a new HTTP request.
func (c *Client) NewRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	if c.opts.baseURL != nil {
		url = fmt.Sprintf("%s%s", *c.opts.baseURL, url)
	}

	if c.overrideURL != nil {
		url = *c.overrideURL
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	if c.opts.auth != nil {
		req.Header.Set("Authorization", *c.opts.auth)
	}

	if c.opts.contentType != nil {
		// Only need to set it for the following methods.
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
			req.Header.Set("Content-Type", *c.opts.contentType)
		}
	}

	if c.opts.headers != nil {
		for k, v := range *c.opts.headers {
			req.Header.Set(k, v)
		}
	}

	if len(c.opts.beforeRequest) > 0 {
		for _, f := range c.opts.beforeRequest {
			err = f(req)
			if err != nil {
				return nil, err
			}
		}
	}

	return req, nil
}

// Do sends an HTTP request and returns an HTTP response.
func (c *Client) Do(req *http.Request) (*Response, error) {
	r, err := c.opts.client.Do(req)
	if err != nil {
		return nil, err
	}

	resp := &Response{r}

	if len(c.opts.afterRequest) > 0 {
		for _, f := range c.opts.afterRequest {
			err = f(resp)
			if err != nil {
				return nil, err
			}
		}
	}

	return resp, nil
}

// Get sends a GET HTTP request and returns an HTTP response.
func (c *Client) Get(ctx context.Context, url string) (*Response, error) {
	r, err := c.NewRequest(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(r)
}

// Post sends a POST HTTP request and returns an HTTP response.
func (c *Client) Post(ctx context.Context, url string, body io.Reader) (*Response, error) {
	r, err := c.NewRequest(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	return c.Do(r)
}

// Put sends a PUT HTTP request and returns an HTTP response.
func (c *Client) Put(ctx context.Context, url string, body io.Reader) (*Response, error) {
	r, err := c.NewRequest(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	return c.Do(r)
}

// Patch sends a PATCH HTTP request and returns an HTTP response.
func (c *Client) Patch(ctx context.Context, url string, body io.Reader) (*Response, error) {
	r, err := c.NewRequest(ctx, http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}

	return c.Do(r)
}

// Delete sends a DELETE HTTP request and returns an HTTP response.
func (c *Client) Delete(ctx context.Context, url string) (*Response, error) {
	r, err := c.NewRequest(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	return c.Do(r)
}

// OverrideURL overrides the request's URL. Mostly for testing
// and using with httptest.Server.
func (c *Client) OverrideURL(url string) {
	c.overrideURL = &url
}

// DefaultHTTPClient returns a http.Client with sane defaults.
// It should be used as a starting point instead of an empty http.Client.
func DefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   defaultTimeout,
		Transport: DefaultHTTPTransport(),
	}
}

// DefaultHTTPTransport returns a http.Transport with sane defaults.
// It should be used as a starting point instead of an empty http.Transport.
func DefaultHTTPTransport() *http.Transport {
	return &http.Transport{
		DialContext: (&net.Dialer{
			Timeout: defaultDialTimeout,
		}).DialContext,
		TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
	}
}

// BuildURLWithQuery build a URL with query strings.
func BuildURLWithQuery(url string, qs map[string]string) (string, error) {
	s, err := netUrl.Parse(url)
	if err != nil {
		return "", err
	}

	q := netUrl.Values{}
	for k, v := range qs {
		q.Add(k, v)
	}

	s.RawQuery = q.Encode()

	return s.String(), nil
}
