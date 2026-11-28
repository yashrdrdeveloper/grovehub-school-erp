package platform

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	corehttp "grovehub/libs/core/http"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/services/workflow"
	sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/generated/registry"
	core "grovehub/libs/core/gateway"
)

const renameScope = "internal.school.registry.rename"
const renameIssuer = "internal-tool-school-rename"

func validSchoolName(s string) bool {
	return s == strings.TrimSpace(s) && utf8.ValidString(s) && utf8.RuneCountInString(s) >= 1 && utf8.RuneCountInString(s) <= 150 && !strings.ContainsFunc(s, unicode.IsControl)
}
func renameInput(call workflow.Invocation, school, path string) (sdk.SchoolRenameInput, bool) {
	var body sdk.SchoolRenameInput
	if call.SchoolID == nil || *call.SchoolID != school || call.Method != "PATCH" || call.Path != path || len(call.Headers) != 0 {
		return body, false
	}
	var raw map[string]json.RawMessage
	if json.Unmarshal(call.Input, &raw) != nil || len(raw) != 4 {
		return body, false
	}
	var id string
	if json.Unmarshal(raw["schoolId"], &id) != nil || id != school || json.Unmarshal(raw["expectedName"], &body.ExpectedName) != nil || json.Unmarshal(raw["name"], &body.Name) != nil || json.Unmarshal(raw["expectedVersion"], &body.ExpectedVersion) != nil {
		return body, false
	}
	return body, body.ExpectedVersion > 0 && validSchoolName(body.ExpectedName) && validSchoolName(body.Name) && body.ExpectedName != body.Name
}

func RenameBindings(c SchoolConfig) (map[workflow.SchoolBindingKey]workflow.SchoolBinding, error) {
	if !canonicalID(c.TemplateID) || !canonicalID(c.SchoolID) || !regexp.MustCompile("^[a-z][a-z0-9_-]{1,63}$").MatchString(c.EnvironmentCode) || c.KeyID == "" {
		return nil, errors.New("invalid rename binding")
	}
	operators := strings.Split(c.OperatorIDs, ",")
	if len(operators) > 100 {
		return nil, errors.New("too many rename operators")
	}
	for _, id := range operators {
		if !canonicalID(id) {
			return nil, errors.New("invalid rename operator")
		}
	}
	origin, err := url.Parse(c.Origin)
	if err != nil || origin.Hostname() == "" || origin.User != nil || origin.RawQuery != "" || origin.ForceQuery || origin.Fragment != "" || origin.RawPath != "" || (origin.Path != "" && origin.Path != "/") || (origin.Scheme != "https" && !(origin.Scheme == "http" && (origin.Hostname() == "127.0.0.1" || origin.Hostname() == "::1"))) {
		return nil, errors.New("invalid rename origin")
	}
	key, err := base64.RawURLEncoding.DecodeString(c.PrivateKey)
	if err != nil || len(key) != ed25519.PrivateKeySize {
		return nil, errors.New("invalid rename key")
	}
	path := "/school/restricted/registry/schools/" + c.SchoolID + "/name"
	validate := func(call workflow.Invocation) bool { _, ok := renameInput(call, c.SchoolID, path); return ok }
	base := http.DefaultTransport.(*http.Transport).Clone()
	base.Proxy = nil
	transport := &core.WorkloadTransport{Origin: origin, KeyID: c.KeyID, Issuer: renameIssuer, Subject: renameIssuer + ":" + c.SchoolID, Audience: "school", PrivateKey: ed25519.PrivateKey(key), Base: base, Policy: func(r *http.Request) (core.WorkloadClaims, error) {
		if r == nil || r.URL == nil || r.Method != "PATCH" || r.URL.Path != path || r.URL.RawPath != "" || r.URL.RawQuery != "" || r.URL.ForceQuery || len(r.Header.Values("Idempotency-Key")) != 1 || !canonicalID(r.Header.Get("Idempotency-Key")) || len(r.Header.Values("Cookie")) != 0 || len(r.Header.Values("X-GH-Identity-Proof")) != 0 {
			return core.WorkloadClaims{}, core.ErrWorkloadDenied
		}
		return core.WorkloadClaims{Scopes: []string{renameScope}}, nil
	}}
	client, err := core.NewWorkloadHTTPClient(transport)
	if err != nil {
		return nil, err
	}
	client.Transport = schoolResponseTransport{base: client.Transport}
	cfg := sdk.NewConfiguration()
	cfg.Servers = sdk.ServerConfigurations{{URL: strings.TrimRight(origin.String(), "/")}}
	cfg.HTTPClient = client
	api := sdk.NewAPIClient(cfg)
	adapter := func(ctx context.Context, call workflow.Invocation) (int, error) {
		body, ok := renameInput(call, c.SchoolID, path)
		if !ok || !canonicalID(call.ExecutionID) {
			return 0, core.ErrWorkloadDenied
		}
		if id := corehttp.ID(ctx); id != "" {
			body.CorrelationId = &id
		}
		resp, err := api.DefaultAPI.RenameSchool(ctx, c.SchoolID).IdempotencyKey(call.ExecutionID).SchoolRenameInput(body).Execute()
		if resp == nil {
			return 0, errors.New("rename unavailable")
		}
		if resp.Body != nil {
			defer resp.Body.Close()
		}
		if resp.StatusCode == 204 && err != nil {
			return 0, errors.New("invalid rename acknowledgment")
		}
		switch resp.StatusCode {
		case 204, 400, 403, 404, 409, 503:
			return resp.StatusCode, nil
		}
		return 0, errors.New("unexpected rename status")
	}
	return map[workflow.SchoolBindingKey]workflow.SchoolBinding{{EnvironmentCode: c.EnvironmentCode, TemplateID: c.TemplateID, SchoolID: c.SchoolID}: {Adapter: adapter, Validate: validate, Operators: operators, Mutation: true}}, nil
}
