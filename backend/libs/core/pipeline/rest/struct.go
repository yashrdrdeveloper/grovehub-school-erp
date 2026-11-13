package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"net/url"

	validation "grovehub/libs/core/annotation/validation"
	coreerrors "grovehub/libs/core/errors"
	httprequest "grovehub/libs/core/http/request"
	corerequest "grovehub/libs/core/request"
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
func RunStruct[Request corerequest.Contract, Response any](writer http.ResponseWriter, request *http.Request, config StructConfig[Request, Response]) {
	if config.Service == nil {
		writeFailure(writer, request, config.Writer, errors.New("missing REST service callback"), config.Fallback, config.Fallback)
		return
	}
	var input Request
	sections, err := corerequest.Inspect(&input)
	if err != nil {
		writeFailure(writer, request, config.Writer, err, config.Fallback, config.Fallback)
		return
	}
	if sections.HasParams {
		if err = bindParamsInto(sections.Params, request); err != nil {
			writeFailure(writer, request, config.Writer, err, config.InvalidParams, config.Fallback)
			return
		}
	}
	if sections.HasQuery {
		if err = bindQueryInto(sections.Query, request, config.StrictQuery); err != nil {
			writeFailure(writer, request, config.Writer, err, config.InvalidQuery, config.Fallback)
			return
		}
	} else if config.StrictQuery && request.URL.RawQuery != "" {
		values, parseErr := url.ParseQuery(request.URL.RawQuery)
		if parseErr != nil {
			err = BindError{Source:"query",Field:"query",Reason:"malformed"}
		} else if len(values) > 0 {
			err = BindError{Source:"query",Field:"query",Reason:"not declared"}
		}
		if err != nil {
			writeFailure(writer, request, config.Writer, err, config.InvalidQuery, config.Fallback)
			return
		}
	}
	if sections.HasRequestBody {
		decode := config.Decode
		if decode == nil { decode = httprequest.DecodeJSON }
		if err = decode(request, sections.RequestBody.Addr().Interface()); err != nil {
			writeFailure(writer, request, config.Writer, &DecodeError{Err:err}, config.InvalidRequestBody, config.Fallback)
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
