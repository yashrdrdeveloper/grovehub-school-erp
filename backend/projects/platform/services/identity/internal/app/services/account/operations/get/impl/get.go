package impl

import (
	"context"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/services/security/principal"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	get "grovehub/services/identity/internal/app/services/account/operations/get"
)

type Service struct {
	accounts accountrepo.Reader
	mapper accountmapper.Mapper
}

var _ get.Operation = (*Service)(nil)

func New(accounts accountrepo.Reader, mapper accountmapper.Mapper) get.Operation {
	return &Service{accounts: accounts, mapper: mapper}
}

func (op *Service) Execute(ctx context.Context, accountID string) (accountdto.Account, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}
	accountEntity, err := op.accounts.FindByID(ctx, accountID)
	if err != nil {
		return accountdto.Account{}, err
	}
	if accountEntity.TenantID == nil || *accountEntity.TenantID != actor.TenantID {
		return accountdto.Account{}, accountrepo.ErrNotFound
	}
	return op.mapper.ToResponse(accountEntity), nil
}
