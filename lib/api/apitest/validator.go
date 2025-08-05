package apitest

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

// ValidateResponseSchema validates rec against our internal OpenAPI schema.
// req is used to find the reference schema in the OpenAPI document.
func ValidateResponseSchema(rec *httptest.ResponseRecorder, req *http.Request, e *echo.Echo) error {
	route, pathParams, err := e.Binder.(*api.EchoBinder).SchemaRouter.FindRoute(req)
	if err != nil {
		return fmt.Errorf("find route: %w", err)
	}

	reqInput := &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
	}

	resInput := &openapi3filter.ResponseValidationInput{
		RequestValidationInput: reqInput,
		Status:                 rec.Code,
		Header:                 rec.Header(),
	}
	b := rec.Body.Bytes()
	defer func() { rec.Body = bytes.NewBuffer(b) }()
	resInput.SetBodyBytes(b)

	return openapi3filter.ValidateResponse(context.Background(), resInput)
}
