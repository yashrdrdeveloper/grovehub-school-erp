package authentication

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestGeneratedProtectedAndGrantOperations(t *testing.T) {
	const sessionID = "e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
	const proof = "route-selected-test-proof"
	const status = `{"status":"OK"}`
	const grant = `{"verified":true,"purpose":"PASSWORD_RESET","expiresIn":600}`
	type invoke func(*Client, context.Context) (Result, error)
	cases := []struct {
		name, method, path, body, data string
		call invoke
	}{
		{"refresh", "POST", "/auth/refresh", "", `{"sessionId":"`+sessionID+`","expiresIn":900,"refreshed":true}`, func(c *Client, ctx context.Context)(Result,error){return c.Refresh(ctx,proof)}},
		{"recovery validation", "GET", "/auth/recovery/validate", "", grant, func(c *Client, ctx context.Context)(Result,error){return c.RecoveryValidate(ctx,proof)}},
		{"invitation validation", "GET", "/auth/invitations/validate", "", `{"verified":true,"purpose":"ACCOUNT_SETUP","expiresIn":600}`, func(c *Client, ctx context.Context)(Result,error){return c.InvitationValidate(ctx,proof)}},
		{"recovery exchange", "POST", "/auth/recovery/exchange", "", grant, func(c *Client, ctx context.Context)(Result,error){return c.RecoveryExchange(ctx,proof)}},
		{"invitation exchange", "POST", "/auth/invitations/exchange", "", `{"verified":true,"purpose":"ACCOUNT_SETUP","expiresIn":600}`, func(c *Client, ctx context.Context)(Result,error){return c.InvitationExchange(ctx,proof)}},
		{"reset status", "GET", "/auth/reset-password/status", "", `{"verified":true,"purpose":"RESET_GRANT","expiresIn":600}`, func(c *Client, ctx context.Context)(Result,error){return c.ResetStatus(ctx,proof)}},
		{"reset password", "POST", "/auth/reset-password", `{"newPassword":"test-password"}`, status, func(c *Client, ctx context.Context)(Result,error){return c.ResetPassword(ctx,ResetPassword{NewPassword:"test-password"},proof)}},
		{"setup status", "GET", "/accounts/setup/status", "", `{"verified":true,"expiresIn":600}`, func(c *Client, ctx context.Context)(Result,error){return c.SetupStatus(ctx,proof)}},
		{"setup", "POST", "/accounts/setup", `{"password":"test-password"}`, `{"accountId":"`+sessionID+`","personGid":"`+sessionID+`","status":"ACTIVE","setupCompleted":true}`, func(c *Client, ctx context.Context)(Result,error){return c.Setup(ctx,Setup{Password:"test-password"},proof)}},
		{"logout", "POST", "/auth/logout", "", status, func(c *Client, ctx context.Context)(Result,error){return c.Logout(ctx,proof)}},
		{"logout all", "POST", "/auth/logout-all", "", `{"status":"ALL_SESSIONS_REVOKED","revokedCount":1}`, func(c *Client, ctx context.Context)(Result,error){return c.LogoutAll(ctx,proof)}},
		{"change password", "POST", "/auth/change-password", `{"newPassword":"test-password"}`, status, func(c *Client, ctx context.Context)(Result,error){return c.ChangePassword(ctx,ChangePassword{NewPassword:"test-password"},proof)}},
		{"reauthenticate", "POST", "/auth/reauthenticate", `{"password":"test-password"}`, status, func(c *Client, ctx context.Context)(Result,error){return c.Reauthenticate(ctx,Reauthenticate{Password:"test-password"},proof)}},
		{"sessions", "GET", "/auth/sessions", "", `{"items":[]}`, func(c *Client, ctx context.Context)(Result,error){return c.Sessions(ctx,proof)}},
		{"revoke", "DELETE", "/auth/sessions/"+sessionID, "", `{"sessionId":"`+sessionID+`","status":"REVOKED"}`, func(c *Client, ctx context.Context)(Result,error){return c.RevokeSession(ctx,sessionID,proof)}},
		{"me", "GET", "/auth/me", "", `{"accountId":"`+sessionID+`","personGid":"`+sessionID+`","username":null,"primaryRoleId":null,"secondaryRoleIds":[],"email":"test@example.com","status":"ACTIVE","version":1,"createdAt":"2026-01-01T00:00:00Z","updatedAt":"2026-01-01T00:00:00Z"}`, func(c *Client, ctx context.Context)(Result,error){return c.Me(ctx,proof)}},
	}
	for _, tc := range cases {
		t.Run(tc.name,func(t *testing.T){
			response := `{"success":true,"data":`+tc.data+`}`
			client := newTestClient(t,func(w http.ResponseWriter,r *http.Request){
				if r.Method!=tc.method || r.URL.Path!="/identity/restricted"+tc.path { t.Errorf("unexpected operation %s %s",r.Method,r.URL.Path) }
				if r.Header.Get("X-GH-Identity-Proof")!=proof || r.Header.Get("Cookie")!="" { t.Error("incorrect credential forwarding") }
				body,err:=io.ReadAll(r.Body)
				if err!=nil || strings.TrimSpace(string(body))!=tc.body { t.Error("incorrect generated request body") }
				w.Header().Set("Content-Type","application/json")
				_,_=io.WriteString(w,response)
			})
			result,err:=tc.call(client,context.Background())
			if err!=nil || result.Status!=200 || string(result.Body)!=response { t.Fatalf("unexpected response: %v",err) }
		})
	}
}

func TestGeneratedRevokeRejectsInvalidSessionID(t *testing.T) {
	client:=newTestClient(t,func(w http.ResponseWriter,r *http.Request){t.Error("invalid session ID reached Platform")})
	if _,err:=client.RevokeSession(context.Background(),"../invalid","proof");err==nil {t.Fatal("expected invalid session error")}
}
