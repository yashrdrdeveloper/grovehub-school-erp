// Package activity adapts the generated Audit SDK to the School ingress.
package activity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/audit"
	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/platformhttp"
)

const maxResponseBytes = platformhttp.MaxResponseBytes

var ErrInvalidRequest = errors.New("invalid audit search request")

type Client struct {
	api *sdk.APIClient
}

type Result struct {
	Status int
	Body []byte
}

// New reuses the configured signed workload client. The wrapper enforces the
// response limit before the generated SDK reads or decodes any response.
func New(origin *url.URL, client *http.Client) *Client {
	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: strings.TrimRight(origin.String(), "/")}}
	cfg.HTTPClient = platformhttp.BoundedClient(client)
	return &Client{api: sdk.NewAPIClient(cfg)}
}

// Search validates against consumer-local generated types, delegates all
// business authorization to Platform, and preserves the original JSON bytes
// rather than round-tripping audit payloads through interface{} / float64.
func (c *Client) Search(ctx context.Context, input json.RawMessage, proof string) (Result, error) {
	var request sdk.SearchRequest
	if len(input) == 0 || !strings.HasPrefix(strings.TrimSpace(string(input)), "{") {
		return Result{}, ErrInvalidRequest
	}
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		return Result{}, ErrInvalidRequest
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return Result{}, ErrInvalidRequest
	}
	ctx = platformhttp.WithIdentityProof(ctx, proof)
	result, resp, err := c.api.DefaultAPI.SearchAuditActivities(ctx).SearchRequest(request).Execute()
	if resp == nil { return Result{}, fmt.Errorf("audit unavailable") }
	defer resp.Body.Close()
	if resp.StatusCode < 300 && (err != nil || result == nil || !result.Success) {
		return Result{}, fmt.Errorf("invalid audit success response")
	}
	raw, readErr := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if readErr != nil || len(raw) > maxResponseBytes || !json.Valid(raw) {
		return Result{}, fmt.Errorf("invalid audit response")
	}
	return Result{Status: resp.StatusCode, Body: raw}, nil
}
