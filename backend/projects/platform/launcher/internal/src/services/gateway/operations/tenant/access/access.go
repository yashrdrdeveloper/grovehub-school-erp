package access

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
	schoolgateway "grovehub/internal/session"
	"grovehub/launcher/internal/src/services/gateway/operations/workload/helpers"
)

type TenantAccessConfig struct{ Origin, KeyID, PrivateKey string }

const TenantAccessSubject = "platform-tenant-access"

const TenantAccessScope = "platform.school.access-check"

func NewTenantAccess(c TenantAccessConfig) (schoolgateway.TenantCheck, error) {
	lookup, err := NewTenantState(c)
	if err != nil {
		return nil, err
	}

	return schoolgateway.ActiveTenantCheck(lookup), nil
}

func NewTenantState(c TenantAccessConfig) (schoolgateway.TenantStateLookup, error) {
	return NewTenantStateWithDependencies(c, core.NewWorkloadHTTPClient, http.NewRequestWithContext)
}

// NewTenantStateWithDependencies composes the signed tenant lookup with
// explicit HTTP client and request construction adapters.
func NewTenantStateWithDependencies(c TenantAccessConfig, clientFactory func(*core.WorkloadTransport) (*http.Client, error), requestFactory func(context.Context, string, string, io.Reader) (*http.Request, error)) (schoolgateway.TenantStateLookup, error) {
	if c == (TenantAccessConfig{}) {
		return nil, nil
	}

	origin, err := url.Parse(c.Origin)
	if err != nil ||
		origin.Hostname() == "" ||
		origin.User != nil ||
		origin.RawPath != "" ||
		origin.RawQuery != "" ||
		origin.ForceQuery ||
		origin.Fragment != "" ||
		(origin.Path != "" && origin.Path != "/") ||
		(origin.Scheme != "https" && !(origin.Scheme == "http" && (origin.Hostname() == "127.0.0.1" || origin.Hostname() == "::1"))) {
		return nil, errors.New("invalid tenant access origin")
	}

	private, err := base64.RawURLEncoding.DecodeString(c.PrivateKey)
	if err != nil || len(private) != ed25519.PrivateKeySize || c.KeyID == "" {
		return nil, errors.New("invalid tenant access credential")
	}

	base := http.DefaultTransport.(*http.Transport).Clone()
	base.Proxy = nil
	client, err := clientFactory(&core.WorkloadTransport{Origin: origin, KeyID: c.KeyID, Issuer: TenantAccessSubject, Subject: TenantAccessSubject, Audience: "school", PrivateKey: ed25519.PrivateKey(private), Base: base, Policy: TenantAccessPolicy})
	if err != nil {
		return nil, err
	}

	return func(ctx context.Context, id string) (schoolgateway.TenantState, error) {
		if !helpers.ValidSchoolID(id) {
			return schoolgateway.TenantState{}, schoolgateway.ErrTenantDenied
		}

		ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		req, err := requestFactory(ctx, "GET", strings.TrimRight(origin.String(), "/")+"/school/restricted/registry/schools/"+id+"/access-check", nil)
		if err != nil {
			return schoolgateway.TenantState{}, schoolgateway.ErrTenantDenied
		}

		resp, err := client.Do(req)
		if err != nil {
			return schoolgateway.TenantState{}, schoolgateway.ErrTenantDenied
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(io.LimitReader(resp.Body, 1))
		if err != nil || (resp.StatusCode != 204 && resp.StatusCode != 403) || len(body) != 0 {
			return schoolgateway.TenantState{}, schoolgateway.ErrTenantDenied
		}

		values := resp.Header.Values("X-GH-Session-Generation")
		if len(values) != 1 {
			return schoolgateway.TenantState{}, schoolgateway.ErrTenantDenied
		}

		generation, err := strconv.ParseInt(values[0], 10, 64)
		if err != nil || generation < 1 || strconv.FormatInt(generation, 10) != values[0] {
			return schoolgateway.TenantState{}, schoolgateway.ErrTenantDenied
		}

		return schoolgateway.TenantState{Generation: generation, Active: resp.StatusCode == 204}, nil
	}, nil
}

// TenantAccessPolicy authorizes only canonical, body-free registry reads.
func TenantAccessPolicy(r *http.Request) (core.WorkloadClaims, error) {
	if r == nil || r.URL == nil {
		return core.WorkloadClaims{}, core.ErrWorkloadDenied
	}

	id := strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/school/restricted/registry/schools/"), "/access-check")
	if !helpers.ValidSchoolID(id) ||
		r.URL.Path != "/school/restricted/registry/schools/"+id+"/access-check" ||
		r.Method != "GET" ||
		r.URL.RawQuery != "" ||
		r.URL.RawPath != "" ||
		r.URL.ForceQuery ||
		r.ContentLength != 0 {
		return core.WorkloadClaims{}, core.ErrWorkloadDenied
	}

	return core.WorkloadClaims{Scopes: []string{TenantAccessScope}}, nil
}
