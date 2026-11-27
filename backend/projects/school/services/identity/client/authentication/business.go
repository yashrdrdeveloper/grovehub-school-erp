package authentication

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	sdk "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/school/internal/generated/identity"
)

var ErrInvalidBusinessRequest = errors.New("invalid identity business request")

// Business routes are selected by School's static scope allowlist, never by a
// caller-supplied upstream URL. Platform owns validation and authorization.
func (c *Client) Business(ctx context.Context, scope, id string, input json.RawMessage, proof, idempotencyKey string) (Result, error) {
	ctx = WithIdentityProof(ctx, proof)
	switch scope {
	case "identity.accounts.create":
		body, err := decodeBusinessBody[sdk.CreateAccountBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		if !validIdempotencyKey(idempotencyKey) { return Result{}, ErrInvalidBusinessRequest }
		request := c.api.DefaultAPI.CreateAccount(ctx).CreateAccountBody(body)
		if idempotencyKey != "" { request = request.IdempotencyKey(idempotencyKey) }
		_, resp, err := request.Execute()
		return sdkResult(resp, err)
	case "identity.accounts.search":
		body, err := decodeBusinessBody[sdk.SearchBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.SearchAccounts(ctx).SearchBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.accounts.roles":
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.ListAssignableRoles(ctx).Execute()
		return sdkResult(resp, err)
	case "identity.accounts.get":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.GetAccount(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.accounts.update":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.UpdateAccountBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.UpdateAccount(ctx, id).UpdateAccountBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.accounts.suspend":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.SuspendAccount(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.accounts.activate":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.ActivateAccount(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.accounts.deactivate":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.DeactivateAccount(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.accounts.resend":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.ResendInvitation(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.accountroles.get":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.GetAccountRoles(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.accountroles.preview":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.ReplaceAccountRolesBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.PreviewAccountRoles(ctx, id).ReplaceAccountRolesBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.history":
		body, err := decodeBusinessBody[sdk.SearchBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.SearchApprovalHistory(ctx).SearchBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.inbox":
		body, err := decodeBusinessBody[sdk.SearchBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.SearchApprovalInbox(ctx).SearchBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.count":
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.CountPendingApprovals(ctx).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.get":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.GetApproval(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.lifecycle":
		body, err := decodeBusinessBody[sdk.LifecycleApprovalBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		if !validIdempotencyKey(idempotencyKey) { return Result{}, ErrInvalidBusinessRequest }
		request := c.api.DefaultAPI.SubmitLifecycleApproval(ctx).LifecycleApprovalBody(body)
		if idempotencyKey != "" { request = request.IdempotencyKey(idempotencyKey) }
		_, resp, err := request.Execute()
		return sdkResult(resp, err)
	case "identity.approvals.permissions":
		body, err := decodeBusinessBody[sdk.PermissionApprovalBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		if !validIdempotencyKey(idempotencyKey) { return Result{}, ErrInvalidBusinessRequest }
		request := c.api.DefaultAPI.SubmitPermissionApproval(ctx).PermissionApprovalBody(body)
		if idempotencyKey != "" { request = request.IdempotencyKey(idempotencyKey) }
		_, resp, err := request.Execute()
		return sdkResult(resp, err)
	case "identity.approvals.assignments":
		body, err := decodeBusinessBody[sdk.AssignmentApprovalBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		if !validIdempotencyKey(idempotencyKey) { return Result{}, ErrInvalidBusinessRequest }
		request := c.api.DefaultAPI.SubmitAssignmentApproval(ctx).AssignmentApprovalBody(body)
		if idempotencyKey != "" { request = request.IdempotencyKey(idempotencyKey) }
		_, resp, err := request.Execute()
		return sdkResult(resp, err)
	case "identity.approvals.approve":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.ApprovalDecisionBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.ApproveRequest(ctx, id).ApprovalDecisionBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.reject":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.ApprovalDecisionBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.RejectRequest(ctx, id).ApprovalDecisionBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.cancel":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.ApprovalVersionBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.CancelRequest(ctx, id).ApprovalVersionBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.approvals.execute":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.ApprovalVersionBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.ExecuteRequest(ctx, id).ApprovalVersionBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.capabilities.me":
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.DiscoverCapabilities(ctx).Execute()
		return sdkResult(resp, err)
	case "identity.capabilities.permissions":
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.GetPermissionCatalog(ctx).Execute()
		return sdkResult(resp, err)
	case "identity.roles.create":
		body, err := decodeBusinessBody[sdk.CreateRoleBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.CreateRoleDraft(ctx).CreateRoleBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.roles.update":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.UpdateRoleBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.UpdateRoleMetadata(ctx, id).UpdateRoleBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.roles.search":
		body, err := decodeBusinessBody[sdk.SearchBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.SearchRoles(ctx).SearchBody(body).Execute()
		return sdkResult(resp, err)
	case "identity.roles.get":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.GetRole(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.roles.permissions":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		if len(bytes.TrimSpace(input)) != 0 { return Result{}, ErrInvalidBusinessRequest }
		_, resp, err := c.api.DefaultAPI.GetRolePermissions(ctx, id).Execute()
		return sdkResult(resp, err)
	case "identity.roles.preview":
		if !ValidSessionID(id) { return Result{}, ErrInvalidBusinessRequest }
		body, err := decodeBusinessBody[sdk.PreviewPermissionsBody](input)
		if err != nil { return Result{}, err }
		ctx = context.WithValue(ctx, businessBodyKey{}, input)
		_, resp, err := c.api.DefaultAPI.PreviewRolePermissions(ctx, id).PreviewPermissionsBody(body).Execute()
		return sdkResult(resp, err)
	default:
		return Result{}, ErrInvalidBusinessRequest
	}
}

// Generated models can normalize explicit null, omitted fields and empty
// arrays. Keep the original JSON on the wire for business requests until
// all nullable/PATCH models have exact round-trip parity. The typed operation
// and model still select the endpoint; Platform's strict decoder remains the
// authoritative validator, including for unknown fields.
type businessBodyKey struct{}

func decodeBusinessBody[T any](input json.RawMessage) (T, error) {
	var value T
	raw := bytes.TrimSpace(input)
	if len(raw) == 0 || raw[0] != '{' || !json.Valid(raw) || json.Unmarshal(raw, &value) != nil {
		return value, ErrInvalidBusinessRequest
	}
	return value, nil
}

func validIdempotencyKey(value string) bool {
	return len(value) <= 256 && !strings.ContainsAny(value, "\r\n\x00")
}

// Apply the original body before workload signing, never after it.
func preserveBusinessBody(req *http.Request) *http.Request {
	raw, ok := req.Context().Value(businessBodyKey{}).(json.RawMessage)
	if !ok { return req }
	clone := req.Clone(req.Context())
	body := append([]byte(nil), raw...)
	clone.Body = newBusinessBody(body)
	clone.ContentLength = int64(len(body))
	clone.GetBody = func() (io.ReadCloser, error) { return newBusinessBody(body), nil }
	clone.Header.Del("Content-Length")
	return clone
}

func newBusinessBody(raw []byte) io.ReadCloser { return io.NopCloser(bytes.NewReader(raw)) }
