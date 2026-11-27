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
	auth "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/services/identity/client/authentication"
)

const maxResponseBytes = 1 << 20

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
	bounded := *client
	bounded.Transport = boundedTransport{base: client.Transport}
	cfg.HTTPClient = &bounded
	return &Client{api: sdk.NewAPIClient(cfg)}
}

type boundedTransport struct { base http.RoundTripper }

func (t boundedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	base := t.base
	if base == nil { base = http.DefaultTransport }
	resp, err := base.RoundTrip(req)
	if err != nil { return nil, err }
	if resp == nil || resp.Body == nil { return nil, fmt.Errorf("empty audit response") }
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 600 || (resp.StatusCode >= 300 && resp.StatusCode < 400) {
		return nil, fmt.Errorf("unexpected audit response")
	}
	raw, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseBytes+1))
	if err != nil || len(raw) > maxResponseBytes || !json.Valid(raw) {
		return nil, fmt.Errorf("invalid audit response")
	}
	resp.Body = io.NopCloser(bytes.NewReader(raw))
	return resp, nil
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
	ctx = auth.WithIdentityProof(ctx, proof)
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
