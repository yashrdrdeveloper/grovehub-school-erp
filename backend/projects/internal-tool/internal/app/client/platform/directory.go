package platform

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
	sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/registry"
	core "grovehub/libs/core/gateway"
)

type DirectoryConfig struct{ Origin, KeyID, PrivateKey, ReaderID string }

const directoryPath = "/school/restricted/registry/lookup"
const directoryScope = "internal.school.registry.read"
const directoryIssuer = "internal-tool-school-directory"

func NewDirectory(c DirectoryConfig) (workflow.Directory, error) {
	if c == (DirectoryConfig{}) {
		return nil, nil
	}
	origin, err := url.Parse(c.Origin)
	if err != nil || origin.Hostname() == "" || origin.User != nil || origin.RawQuery != "" || origin.ForceQuery || origin.Fragment != "" || origin.RawPath != "" || (origin.Path != "" && origin.Path != "/") || (origin.Scheme != "https" && !(origin.Scheme == "http" && (origin.Hostname() == "127.0.0.1" || origin.Hostname() == "::1"))) {
		return nil, errors.New("invalid school directory origin")
	}
	key, err := base64.RawURLEncoding.DecodeString(c.PrivateKey)
	if err != nil || len(key) != ed25519.PrivateKeySize || !canonicalID(c.ReaderID) || c.KeyID == "" {
		return nil, errors.New("invalid school directory credentials")
	}
	base := http.DefaultTransport.(*http.Transport).Clone()
	base.Proxy = nil
	transport := &core.WorkloadTransport{Origin: origin, KeyID: c.KeyID, Issuer: directoryIssuer, Subject: directoryIssuer + ":" + c.ReaderID, Audience: "school", PrivateKey: ed25519.PrivateKey(key), Base: base, Policy: func(r *http.Request) (core.WorkloadClaims, error) {
		if r == nil || r.URL == nil || r.Method != "POST" || r.URL.Path != directoryPath || r.URL.RawPath != "" || r.URL.RawQuery != "" || r.URL.ForceQuery || len(r.Header.Values("Cookie")) != 0 || len(r.Header.Values("X-GH-Identity-Proof")) != 0 {
			return core.WorkloadClaims{}, core.ErrWorkloadDenied
		}
		return core.WorkloadClaims{Scopes: []string{directoryScope}}, nil
	}}
	client, err := core.NewWorkloadHTTPClient(transport)
	if err != nil {
		return nil, err
	}
	client.Transport = directoryResponseTransport{base: client.Transport}
	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: strings.TrimRight(origin.String(), "/")}}
	cfg.HTTPClient = client
	api := sdk.NewAPIClient(cfg)
	return func(ctx context.Context, ids []string) (map[string]workflow.SchoolRecord, error) {
		result := map[string]workflow.SchoolRecord{}
		if len(ids) == 0 {
			return result, nil
		}
		if len(ids) > 100 {
			return nil, errors.New("too many school directory IDs")
		}
		requested := map[string]bool{}
		for _, id := range ids {
			if !canonicalID(id) || requested[id] {
				return nil, errors.New("invalid school directory IDs")
			}
			requested[id] = true
		}
		rows, resp, err := api.DefaultAPI.LookupSchools(ctx).SchoolLookupInput(sdk.SchoolLookupInput{SchoolIds: ids}).Execute()
		if resp != nil && resp.Body != nil {
			defer resp.Body.Close()
		}
		if err != nil || resp == nil || resp.StatusCode != 200 || rows == nil {
			return nil, errors.New("school directory unavailable")
		}
		for _, row := range rows.Records {
			if !requested[row.Id] || result[row.Id].ID != "" || row.Code == "" || row.Name == "" || row.Version < 1 || (row.Status != "ACTIVE" && row.Status != "DISABLED") {
				return nil, errors.New("invalid school directory response")
			}
			result[row.Id] = workflow.SchoolRecord{ID: row.Id, Code: row.Code, Name: row.Name, Status: row.Status, Version: row.Version}
		}
		return result, nil
	}, nil
}

type directoryResponseTransport struct{ base http.RoundTripper }

func (t directoryResponseTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := t.base.RoundTrip(r)
	if err != nil {
		return nil, err
	}
	if resp == nil || resp.Body == nil {
		return nil, errors.New("missing directory response")
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(io.LimitReader(resp.Body, 128*1024+1))
	if err != nil || len(raw) > 128*1024 {
		return nil, errors.New("oversized directory response")
	}
	resp.Body = io.NopCloser(strings.NewReader(string(raw)))
	return resp, nil
}
