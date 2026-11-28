package credential

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway/internal/src/contracts"

	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"
)

const maxCredentialLifetime = 2 * time.Minute

type signedClaims struct {
	Issuer string `json:"iss"`
	Subject string `json:"sub"`
	Audience string `json:"aud"`
	IssuedAt int64 `json:"iat"`
	ExpiresAt int64 `json:"exp"`
	Scopes []string `json:"scopes"`
	TenantID string `json:"tenant,omitempty"`
	ActorGID string `json:"actor,omitempty"`
	SessionID string `json:"session,omitempty"`
}

// Ed25519Credentials verifies a versioned internal credential, not a JWT.
// Keys and scope policies must be immutable while requests are being served.
type Ed25519Credentials struct {
	Keys map[string]contracts.TrustedWorkloadKey
}

func SignWorkload(keyID, issuer string, key ed25519.PrivateKey, claims contracts.WorkloadClaims) (string, error) {
	now := time.Now()
	if keyID == "" || len(keyID) > 128 || strings.ContainsAny(keyID, ". \t\r\n,") || issuer == "" || len(key) != ed25519.PrivateKeySize || claims.Subject == "" || claims.Audience == "" || !claims.ExpiresAt.After(now) || claims.ExpiresAt.After(now.Add(maxCredentialLifetime)) || (claims.TenantID == "") != (claims.ActorGID == "") {
		return "", contracts.ErrWorkloadDenied
	}
	payload, err := json.Marshal(signedClaims{Issuer: issuer, Subject: claims.Subject, Audience: claims.Audience, IssuedAt: now.Unix(), ExpiresAt: claims.ExpiresAt.Unix(), Scopes: claims.Scopes, TenantID: claims.TenantID, ActorGID: claims.ActorGID, SessionID: claims.SessionID})
	if err != nil {
		return "", err
	}
	body := "ghw1."+keyID+"."+base64.RawURLEncoding.EncodeToString(payload)
	token := body+"."+base64.RawURLEncoding.EncodeToString(ed25519.Sign(key, []byte(body)))
	if len(token) > 8192 {
		return "", contracts.ErrWorkloadDenied
	}
	return token, nil
}

func (v Ed25519Credentials) VerifyCredential(ctx context.Context, token string) (contracts.WorkloadClaims, error) {
	if ctx == nil || ctx.Err() != nil || len(token) > 8192 {
		return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
	}
	parts := strings.Split(token, ".")
	if len(parts) != 4 || parts[0] != "ghw1" {
		return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
	}
	key, ok := v.Keys[parts[1]]
	if !ok || len(key.PublicKey) != ed25519.PublicKeySize || key.Issuer == "" || key.Subject == "" {
		return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[3])
	if err != nil || !ed25519.Verify(key.PublicKey, []byte(strings.Join(parts[:3], ".")), signature) {
		return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[2])
	var claims signedClaims
	if err != nil || json.Unmarshal(payload, &claims) != nil {
		return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
	}
	now := time.Now().Unix()
	if claims.Issuer != key.Issuer || claims.Subject != key.Subject || claims.Audience == "" || claims.IssuedAt <= 0 || claims.IssuedAt > now || claims.ExpiresAt <= now || claims.ExpiresAt <= claims.IssuedAt || claims.ExpiresAt-claims.IssuedAt > int64(maxCredentialLifetime/time.Second) {
		return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
	}
	if (claims.TenantID == "") != (claims.ActorGID == "") || ((claims.ActorGID != "" || claims.SessionID != "") && !key.AllowDelegation) {
		return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
	}
	for _, scope := range claims.Scopes {
		allowed := false
		for _, trusted := range key.Scopes {
			if scope != "" && scope == trusted {
				allowed = true
				break
			}
		}
		if !allowed {
			return contracts.WorkloadClaims{}, contracts.ErrWorkloadDenied
		}
	}
	return contracts.WorkloadClaims{Subject: claims.Subject, Audience: claims.Audience, ExpiresAt: time.Unix(claims.ExpiresAt, 0), Scopes: claims.Scopes, TenantID: claims.TenantID, ActorGID: claims.ActorGID, SessionID: claims.SessionID}, nil
}
