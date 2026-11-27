package roles

import (
	"context"

	rolecatalog "grovehub/services/identity/internal/app/constants/role"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/security/principal"
	"grovehub/services/identity/internal/app/repository/role"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/security/operations/routepolicy"
)

// List exposes the server-owned assignable catalog to authenticated tenant members.
func List(ctx context.Context) ([]accountdto.Role, error) {
	if _, err := principal.RequireTenant(ctx); err != nil {
		return nil, err
	}
	return rolecatalog.Catalog(), nil
}

type Operation struct {
	reader role.CatalogReader
	policy *routepolicy.Policy
	legacy bool
}

func New(reader role.CatalogReader, policy *routepolicy.Policy, legacy bool) *Operation {
	return &Operation{reader: reader, policy: policy, legacy: legacy}
}

func (operation *Operation) Execute(ctx context.Context) ([]accountdto.Role, error) {
	if operation.policy == nil {
		return nil, authorization.ErrUnavailable
	}
	if err := operation.policy.Authorize(ctx, "GET /identity/protected/accounts/roles"); err != nil {
		return nil, err
	}
	if operation.legacy {
		return List(ctx)
	}
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return nil, err
	}
	if operation.reader == nil {
		return nil, authorization.ErrUnavailable
	}
	return operation.reader.ListActive(ctx, actor.TenantID)
}
