// Package workload constructs cryptographically verified requests for tests.
package workload

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

func Key(id, issuer, subject string, scopes []string) (core.WorkloadKeyConfig, ed25519.PrivateKey, error) {
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return core.WorkloadKeyConfig{}, nil, err
	}

	return core.WorkloadKeyConfig{
		KeyID: id, Issuer: issuer, Subject: subject,
		PublicKey: base64.RawURLEncoding.EncodeToString(public), Scopes: scopes,
	}, private, nil
}

func Request(pattern, path string, key core.WorkloadKeyConfig, private ed25519.PrivateKey, claims core.WorkloadClaims) (*http.Request, error) {
	method, _, _ := strings.Cut(pattern, " ")
	claims.Subject = key.Subject
	claims.Audience = "platform"
	claims.ExpiresAt = time.Now().Add(time.Minute)
	token, err := core.SignWorkload(key.KeyID, key.Issuer, private, claims)
	if err != nil {
		return nil, err
	}

	request := httptest.NewRequest(method, path, nil)
	request.Pattern = pattern
	request.Header.Set("Authorization", "Bearer "+token)
	return request, nil
}

func VerifiedContext(subject string) (context.Context, error) {
	const pattern = "GET /fixture"
	const scope = "fixture.read"
	key, private, err := Key("fixture", subject, subject, []string{scope})
	if err != nil {
		return nil, err
	}

	credentials, err := core.NewEd25519Credentials([]core.WorkloadKeyConfig{key})
	if err != nil {
		return nil, err
	}

	request, err := Request(pattern, "/fixture", key, private, core.WorkloadClaims{Scopes: []string{scope}})
	if err != nil {
		return nil, err
	}

	verifier := core.ScopedWorkloadVerifier{
		Credentials: credentials, Audience: "platform",
		Routes: map[string]core.RoutePolicy{pattern: {Scope: scope}},
	}
	return verifier.VerifyWorkload(request.Context(), request)
}
