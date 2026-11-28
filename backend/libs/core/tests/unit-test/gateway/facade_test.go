package gateway_test

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"testing"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

func TestGatewayFacade(t *testing.T) {
	public, _, err := ed25519.GenerateKey(nil)
	if err != nil { t.Fatal(err) }
	_, err = gateway.NewEd25519Credentials([]gateway.WorkloadKeyConfig{{KeyID:"key",PublicKey:base64.RawURLEncoding.EncodeToString(public),Issuer:"issuer",Subject:"subject",Scopes:[]string{"read"}}})
	if err != nil { t.Fatal(err) }
	if _, ok := gateway.WorkloadFrom(context.Background()); ok { t.Fatal("unexpected claims") }
	if gateway.Public != "public" || gateway.Protected != "protected" || gateway.Restricted != "restricted" { t.Fatal("route tier changed") }
}
