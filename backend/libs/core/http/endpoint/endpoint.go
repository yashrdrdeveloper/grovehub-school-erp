// Package endpoint dispatches typed in-process client calls through the
// owning service's registered HTTP controller routes. It never opens a socket.
package endpoint

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	coreerrors "grovehub/libs/core/errors"
	headercontext "grovehub/libs/core/headercontext"
	httpheaders "grovehub/libs/core/http/headers"
)

const maxResponseBytes int64 = 8 << 20

var placeholder = regexp.MustCompile(`\{([A-Za-z][A-Za-z0-9_]*)\}`)

// Call executes one HTTP controller handler in process, using the same
// routing, decoding, normalization, validation and response envelope as
// external HTTP calls. pathDTO supplies path:"name" tagged fields.
//
// The handler must be a registered, trusted service router. This is not
// a network transport and does not automatically run outer gateway middleware.
func Call[T any](ctx context.Context, handler http.Handler, method, route string, pathDTO any, body any) (T, int, error) {
	var zero T
	if ctx == nil { return zero, 0, errors.New("HTTP endpoint call requires context") }
	if err := ctx.Err(); err != nil { return zero, 0, err }
	if handler == nil { return zero, 0, errors.New("HTTP endpoint call requires a handler") }
	path, err := expandPath(route, pathDTO)
	if err != nil { return zero, 0, err }

	var payload io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil { return zero, 0, fmt.Errorf("encode endpoint body: %w", err) }
		payload = bytes.NewReader(encoded)
	}
	request, err := http.NewRequestWithContext(ctx, method, "http://grovehub.internal"+path, payload)
	if err != nil { return zero, 0, fmt.Errorf("construct endpoint request: %w", err) }
	if body != nil { request.Header.Set("Content-Type", "application/json") }

	// Validated caller context is preserved, including tracing, transactions
	// and trusted metadata. Headers are forwarded for future HTTP middleware.
	values := headercontext.FromContext(ctx)
	if values.IdempotencyKey != "" { request.Header.Set(httpheaders.IdempotencyKey, values.IdempotencyKey) }
	if values.CorrelationID != "" { request.Header.Set(httpheaders.CorrelationID, values.CorrelationID) }
	if values.RequestID != "" { request.Header.Set(httpheaders.RequestID, values.RequestID) }

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	status := response.StatusCode
	if status == http.StatusNoContent { return zero, status, nil }
	if status < 200 || status >= 600 {
		return zero, status, fmt.Errorf("invalid endpoint response status %d", status)
	}

	data, err := io.ReadAll(io.LimitReader(response.Body, maxResponseBytes+1))
	if err != nil { return zero, status, fmt.Errorf("read endpoint response: %w", err) }
	if int64(len(data)) > maxResponseBytes {
		return zero, status, errors.New("endpoint response exceeds maximum size")
	}

	var envelope struct {
		Success       bool            `json:"success"`
		Data          json.RawMessage `json:"data"`
		ErrorCode     coreerrors.Code `json:"errorCode"`
		ErrorMessage  string          `json:"errorMessage"`
		HTTPErrorCode int             `json:"httpErrorCode"`
		ErrorMetaData json.RawMessage `json:"errorMetaData"`
	}
	if err := json.Unmarshal(data, &envelope); err != nil {
		return zero, status, fmt.Errorf("invalid endpoint response envelope: %w", err)
	}
	if status >= http.StatusBadRequest {
		if envelope.Success || envelope.ErrorCode == "" || envelope.HTTPErrorCode != status {
			return zero, status, errors.New("invalid endpoint error envelope")
		}
		definition := coreerrors.NewDefinition(envelope.ErrorCode, envelope.ErrorMessage, status)
		if len(envelope.ErrorMetaData) == 0 || bytes.Equal(bytes.TrimSpace(envelope.ErrorMetaData), []byte("null")) {
			return zero, status, coreerrors.NewError(definition)
		}
		var metadata any
		if err := json.Unmarshal(envelope.ErrorMetaData, &metadata); err != nil {
			return zero, status, fmt.Errorf("invalid endpoint error metadata: %w", err)
		}
		return zero, status, coreerrors.NewErrorWithMetadata(definition, metadata)
	}
	if !envelope.Success {
		return zero, status, errors.New("unexpected endpoint success envelope")
	}
	if len(envelope.Data) == 0 || bytes.Equal(bytes.TrimSpace(envelope.Data), []byte("null")) {
		return zero, status, nil
	}
	var result T
	if err := json.Unmarshal(envelope.Data, &result); err != nil {
		return zero, status, fmt.Errorf("decode endpoint result: %w", err)
	}
	return result, status, nil
}

// expandPath resolves route placeholders exclusively from tagged DTO fields.
// No value may change the route's structure or be injected into a query string.
func expandPath(route string, pathDTO any) (string, error) {
	if !strings.HasPrefix(route, "/") || strings.ContainsAny(route, "?#") {
		return "", errors.New("endpoint route must be a path without query or fragment")
	}
	placeholders := placeholder.FindAllStringSubmatch(route, -1)
	if len(placeholders) == 0 {
		if strings.ContainsAny(route, "{}") { return "", errors.New("invalid endpoint route placeholder") }
		return route, nil
	}
	value := reflect.ValueOf(pathDTO)
	if !value.IsValid() { return "", errors.New("endpoint path DTO is missing") }
	for value.Kind() == reflect.Pointer {
		if value.IsNil() { return "", errors.New("endpoint path DTO is nil") }
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return "", errors.New("endpoint path DTO must be a struct")
	}
	fields := make(map[string]string)
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := strings.Split(field.Tag.Get("path"), ",")[0]
		if tag == "" || tag == "-" { continue }
		if !field.IsExported() { return "", fmt.Errorf("endpoint path field %s must be exported", field.Name) }
		if _, exists := fields[tag]; exists { return "", fmt.Errorf("duplicate endpoint path field %s", tag) }
		raw := value.Field(i)
		if raw.Kind() == reflect.Pointer {
			if raw.IsNil() { return "", fmt.Errorf("endpoint path field %s is nil", tag) }
			raw = raw.Elem()
		}
		var rendered string
		switch raw.Kind() {
		case reflect.String:
			rendered = raw.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			rendered = strconv.FormatInt(raw.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			rendered = strconv.FormatUint(raw.Uint(), 10)
		default:
			return "", fmt.Errorf("unsupported endpoint path field type for %s", tag)
		}
		if rendered == "" || rendered == "." || rendered == ".." {
			return "", fmt.Errorf("endpoint path field %s is empty or unsafe", tag)
		}
		fields[tag] = url.PathEscape(rendered)
	}
	result := placeholder.ReplaceAllStringFunc(route, func(match string) string {
		return fields[match[1:len(match)-1]]
	})
	if strings.ContainsAny(result, "{}") || strings.Contains(result, "//") {
		return "", errors.New("unresolved or unsafe endpoint path")
	}
	if result == route { return "", errors.New("unresolved endpoint path placeholders") }
	for _, match := range placeholders {
		if _, exists := fields[match[1]]; !exists {
			return "", fmt.Errorf("endpoint path field %s is missing", match[1])
		}
	}
	return result, nil
}
