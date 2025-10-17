package api

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func init() {
	openapi3.SchemaErrorDetailsDisabled = true
	openapi3.DefineStringFormatValidator("uuid", uuidFormatValidator{})
	for _, a := range []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",   // .docx
		"application/vnd.openxmlformats-officedocument.presentationml.presentation", // .pptx
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",         // .xlsx
		"application/msword",            // .doc
		"application/vnd.ms-excel",      // .xls
		"application/vnd.ms-powerpoint", // .ppt
		"application/pdf",               // .pdf
		"image/jpeg",                    // .jpg .jpeg
		"image/png",                     // .png
	} {
		openapi3filter.RegisterBodyDecoder(a, openapi3filter.FileBodyDecoder)
	}
}

// EchoBinder implements echo.Binder.
type EchoBinder struct {
	SchemaRouter  routers.Router
	defaultBinder echo.Binder
}

func newEchoBinder(openapiBytes []byte) (*EchoBinder, error) {
	doc, err := openapi3.NewLoader().LoadFromData(openapiBytes)
	if err != nil {
		return nil, fmt.Errorf("load openapi blob: %w", err)
	}

	// Skip scheme & host validation.
	doc.Servers = nil

	s, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, fmt.Errorf("new schema router: %w", err)
	}

	return &EchoBinder{SchemaRouter: s, defaultBinder: new(echo.DefaultBinder)}, nil
}

func (b *EchoBinder) Bind(i any, c echo.Context) error {
	if err := validateRequest(c.Request(), b.SchemaRouter); err != nil {
		return err
	}

	if err := b.defaultBinder.Bind(i, c); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}
	return nil
}

func validateRequest(req *http.Request, router routers.Router) error {
	route, pathParams, err := router.FindRoute(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "openapi route tidak ditemukan")
	}

	if err := openapi3filter.ValidateRequest(req.Context(), &openapi3filter.RequestValidationInput{
		Request:    req,
		PathParams: pathParams,
		Route:      route,
		Options: &openapi3filter.Options{
			MultiError: true,
		},
	}); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, requestMultiErrorMessage(err.(openapi3.MultiError)))
	}
	return nil
}

func requestMultiErrorMessage(merr openapi3.MultiError) string {
	msgs := []string{}
	for _, err := range merr {
		rerr := err.(*openapi3filter.RequestError)
		prefix := ""
		if rerr.Parameter != nil {
			prefix = fmt.Sprintf(`parameter "%s" `, rerr.Parameter.Name)
		}

		if err = rerr.Err; err == nil {
			if strings.HasPrefix(rerr.Reason, "header Content-Type has unexpected value") {
				msgs = append(msgs, "header Content-Type harus dalam format yang sesuai")
			} else {
				msgs = append(msgs, rerr.Reason)
			}
		} else if errs, ok := err.(openapi3.MultiError); ok {
			msgs2 := []string{}
			for _, err2 := range errs {
				msgs2 = append(msgs2, schemaErrorMessage(err2.(*openapi3.SchemaError)))
			}
			msgs = append(msgs, prefix+strings.Join(msgs2, " | "))
		} else if serr, ok := err.(*openapi3.SchemaError); ok {
			msgs = append(msgs, prefix+schemaErrorMessage(serr))
		} else if perr, ok := err.(*openapi3filter.ParseError); ok {
			if prefix == "" {
				prefix = "request body "
				if len(perr.Path()) > 0 {
					prefix = fmt.Sprintf(`parameter "%v" `, perr.Path()[0])
				}
			}
			msgs = append(msgs, prefix+"harus dalam format yang sesuai")
		} else if herr, ok := err.(*echo.HTTPError); ok && herr.Code == 413 {
			msgs = append(msgs, "request body terlalu besar")
		} else if err.Error() == "value is required but missing" {
			if prefix == "" {
				prefix = "request body "
			}
			msgs = append(msgs, prefix+"harus diisi")
		} else if err.Error() == "empty value is not allowed" {
			msgs = append(msgs, prefix+"tidak boleh kosong")
		} else {
			msgs = append(msgs, err.Error())
		}
	}
	return strings.Join(msgs, " | ")
}

var quotedRegexp = regexp.MustCompile(`"(.+?)"`)

func schemaErrorMessage(err *openapi3.SchemaError) string {
	var msg string

	switch err.SchemaField {
	case "format":
		msg = "harus dalam format " + err.Schema.Format
	case "enum":
		vals := make([]string, len(err.Schema.Enum))
		for i, v := range err.Schema.Enum {
			vals[i] = fmt.Sprintf("%#v", v)
		}
		msg = "harus salah satu dari " + strings.Join(vals, ", ")
	case "minimum":
		msg = "harus tidak kurang dari " + strconv.FormatFloat(*err.Schema.Min, 'f', -1, 64)
	case "maximum":
		msg = "harus tidak lebih dari " + strconv.FormatFloat(*err.Schema.Max, 'f', -1, 64)
	case "minLength":
		msg = fmt.Sprintf("harus %d karakter atau lebih", err.Schema.MinLength)
	case "maxLength":
		msg = fmt.Sprintf("harus %d karakter atau kurang", *err.Schema.MaxLength)
	case "minItems":
		msg = fmt.Sprintf("harus %d item atau lebih", err.Schema.MinItems)
	case "maxItems":
		msg = fmt.Sprintf("harus %d item atau kurang", *err.Schema.MaxItems)
	case "uniqueItems":
		msg = "item tidak boleh duplikat"
	case "required":
		msg = "harus diisi"
	case "properties":
		if strings.HasSuffix(err.Reason, "is unsupported") {
			field := append(err.JSONPointer(), quotedRegexp.FindStringSubmatch(err.Reason)[1])
			return fmt.Sprintf(`parameter "%s" tidak didukung`, strings.Join(field, "."))
		}
		return err.Error()
	case "minProperties":
		msg = fmt.Sprintf("harus %d property atau lebih", err.Schema.MinProps)
	case "maxProperties":
		msg = fmt.Sprintf("harus %d property atau kurang", *err.Schema.MaxProps)
	case "type":
		msg = "harus dalam tipe " + strings.Join(*err.Schema.Type, ", ")
	case "nullable":
		msg = "tidak boleh null"
	case "allOf", "anyOf", "oneOf":
		msg = "parameter tidak valid"
	default:
		return err.Error()
	}

	if field := err.JSONPointer(); len(field) > 0 {
		msg = fmt.Sprintf(`parameter "%s" %s`, strings.Join(field, "."), msg)
	} else if err.SchemaField == "minProperties" || err.SchemaField == "maxProperties" {
		msg = "request body " + msg
	}
	return msg
}

type uuidFormatValidator struct{}

func (uuidFormatValidator) Validate(value string) error {
	if uuid.Validate(value) != nil {
		return errors.New("invalid UUID")
	}
	return nil
}
