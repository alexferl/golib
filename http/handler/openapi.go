package handler

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
)

const (
	ApplicationJSON = "application/json"
	ApplicationXML  = "application/xml"
)

type OpenAPI struct {
	Config OpenAPIConfig
}

type OpenAPIConfig struct {
	ContentType  string
	ValidatorKey string
	// Set ExcludeRequestBody so ValidateRequest skips request body validation
	ExcludeRequestBody bool

	// Set ExcludeResponseBody so ValidateResponse skips response body validation
	ExcludeResponseBody bool

	// Set IncludeResponseStatus so ValidateResponse fails on response
	// status not defined in OpenAPI spec
	IncludeResponseStatus bool

	MultiError bool
}

var DefaultOpenAPIConfig = OpenAPIConfig{
	ContentType:           ApplicationJSON,
	ValidatorKey:          "validator",
	ExcludeRequestBody:    false,
	ExcludeResponseBody:   false,
	IncludeResponseStatus: true,
	MultiError:            true,
}

func NewOpenAPIHandler() *OpenAPI {
	c := DefaultOpenAPIConfig
	return NewOpenAPIWithConfig(c)
}

func NewOpenAPIWithConfig(config OpenAPIConfig) *OpenAPI {
	return &OpenAPI{Config: config}
}

func (h *OpenAPI) Validate(c echo.Context, code int, v any) error {
	input := c.Get(h.Config.ValidatorKey).(*openapi3filter.RequestValidationInput)
	responseValidationInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: input,
		Status:                 c.Response().Status,
		Header:                 c.Response().Header(),
		Options: &openapi3filter.Options{
			ExcludeRequestBody:    h.Config.ExcludeRequestBody,
			ExcludeResponseBody:   h.Config.ExcludeResponseBody,
			IncludeResponseStatus: h.Config.IncludeResponseStatus,
			MultiError:            h.Config.MultiError,
		},
	}

	var (
		b   []byte
		err error
	)
	if h.Config.ContentType == ApplicationJSON {
		b, err = json.Marshal(v)
	} else if h.Config.ContentType == ApplicationXML {
		b, err = xml.Marshal(v)
	} else {
		panic(fmt.Sprintf("error content-type %s not supported", h.Config.ContentType))
	}

	if err != nil {
		c.Logger().Errorf("error marshaling response: %v", err)
		return err
	}

	c.Response().Header().Add("Content-Type", h.Config.ContentType)

	responseValidationInput.SetBodyBytes(b)
	ctx := context.Background()
	err = openapi3filter.ValidateResponse(ctx, responseValidationInput)
	if err != nil {
		c.Logger().Debug("error validating response: %v", err)
		return err
	}

	return c.Blob(code, h.Config.ContentType, b)
}
