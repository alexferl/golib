package client

import (
	"encoding/json"
	"io"
	"net/http"
)

type Response struct {
	*http.Response
}

func (r *Response) UnmarshalJSON(v any) error {
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	return nil
}
