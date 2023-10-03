# http.client
A module to help make HTTP requests to REST APIs.

## example
```go
package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexferl/golib/http/client"
)

func before(req *http.Request) error {
	fmt.Printf("sending request to: %s\n", req.URL.String())
	return nil
}

func main() {
	hc := client.DefaultHTTPClient()
	hc.Timeout = 42 * time.Second // set default timeout for every requests
	contentType := "application/json"
	auth := "token"
	url := "https://example.com"
	headers := map[string]string{
		"X-Custom-Header1": "1",
		"X-Custom-Header2": "2",
	}

	c := client.New(
		client.WithHTTPClient(hc),
		client.WithBaseURL(url),
		client.WithAuthBearer(auth),
		client.WithContentType(contentType),
		client.WithHeaders(headers),
		client.WithBeforeRequest(before),
	)

	// Creating and sending requests manually
	req, err := c.NewRequest(context.Background(), http.MethodGet, "/users/1", nil)
	if err != nil {
		panic(err)
	}

	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("received body: %s\n", b)

	// Using method helpers to create and send requests automatically
	// Override the timeout on the http.Client with a context, this is
	// the preferred method over changing the http.Client timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	resp, err = c.Post(ctx, "/users", bytes.NewBuffer([]byte("payload")))
	if err != nil {
		panic(err)
	}

	fmt.Printf("response code: %d\n", resp.StatusCode)
}
```
