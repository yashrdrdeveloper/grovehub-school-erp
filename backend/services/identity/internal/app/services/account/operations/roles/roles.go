package roles

import (
	"context"

	rolecatalog "grovehub/services/identity/internal/app/constants/role"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/security/principal"
)

// List exposes the server-owned assignable catalog to authenticated tenant members.
func List(ctx context.Context) ([]accountdto.Role, error) {
	if _, err := principal.RequireTenant(ctx); err != nil {
		return nil, err
	}
	return rolecatalog.Catalog(), nil
}
