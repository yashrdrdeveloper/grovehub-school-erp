package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	validation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	httprequest "grovehub/libs/core/http/request"
)

// StructConfig operates on one endpoint-specific request struct. Only the
// applicable exported fields Params, Query and RequestBody are declared.
type StructConfig[Request, Response any] struct {
	Rules validation.Rules[Request]
	Validate func(Request) error
	Service func(context.Context, Request) (Result[Response], error)
	Decode func(*http.Request, any) error
	StrictQuery bool
	InvalidParams coreerrors.Definition
	InvalidQuery coreerrors.Definition
	InvalidRequestBody coreerrors.Definition
	Fallback coreerrors.Definition
	Writer Writer[Response]
}

// RunStruct binds only the fields an endpoint declares, validates the complete
// request, invokes its service once, and writes exactly one HTTP response.
func RunStruct[Request, Response any](writer http.ResponseWriter, request *http.Request, config StructConfig[Request, Response]) {
	if config.Service == nil {
		writeFailure(writer, request, config.Writer, errors.New("missing REST service callback"), config.Fallback, config.Fallback)
		return
	}
	var input Request
	target := reflect.ValueOf(&input).Elem()
	if target.Kind() != reflect.Struct {
		writeFailure(writer, request, config.Writer, errors.New("request contract must be a struct"), config.Fallback, config.Fallback)
		return
	}
	for i := 0; i < target.NumField(); i++ {
		field := target.Type().Field(i)
		if !field.IsExported() || !target.Field(i).CanAddr() {
			writeFailure(writer, request, config.Writer, fmt.Errorf("invalid unexported request contract field"), config.Fallback, config.Fallback)
			return
		}
		part := target.Field(i)
		var err error
		var invalid coreerrors.Definition
		switch field.Name {
		case "Params":
			invalid = config.InvalidParams
			err = bindParamsInto(part, request)
		case "Query":
			invalid = config.InvalidQuery
			err = bindQueryInto(part, request, config.StrictQuery)
		case "RequestBody":
			invalid = config.InvalidRequestBody
			decode := config.Decode
			if decode == nil { decode = httprequest.DecodeJSON }
			err = decode(request, part.Addr().Interface())
			if err != nil { err = &DecodeError{Err: err} }
		default:
			err = fmt.Errorf("unsupported request contract field")
			invalid = config.Fallback
		}
		if err != nil {
			writeFailure(writer, request, config.Writer, err, invalid, config.Fallback)
			return
		}
	}
	normalized, err := validation.Process(input, config.Rules)
	if err != nil {
		invalid := config.InvalidRequestBody
		var fieldErr validation.Error
		if errors.As(err, &fieldErr) {
			switch {
			case strings.HasPrefix(fieldErr.Field, "Params."):
				invalid = config.InvalidParams
			case strings.HasPrefix(fieldErr.Field, "Query."):
				invalid = config.InvalidQuery
			}
		}
		writeFailure(writer, request, config.Writer, err, invalid, config.Fallback)
		return
	}
	if config.Validate != nil {
		if err = config.Validate(normalized); err != nil {
			writeFailure(writer, request, config.Writer, err, config.InvalidRequestBody, config.Fallback)
			return
		}
	}
	result, err := config.Service(request.Context(), normalized)
	if err != nil {
		writeFailure(writer, request, config.Writer, err, config.Fallback, config.Fallback)
		return
	}
	writeSuccess(writer, request, config.Writer, result, config.Fallback)
}
