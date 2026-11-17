package impl

import (
	"context"
	"grovehub/services/identity/internal/app/security/principal"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
	me "grovehub/services/identity/internal/app/services/authentication/operations/me"
)

type Service struct {
	accounts accountrepo.Repository
	mapper accountmapper.Mapper
}

var _ me.Operation = (*Service)(nil)

func New(accounts accountrepo.Repository, mapper accountmapper.Mapper) me.Operation {
	return &Service{accounts: accounts, mapper: mapper}
}

func (op *Service) Execute(ctx context.Context) (accountdto.Account, error) {
	actor, err := principal.Require(ctx)
	if err != nil {
		return accountdto.Account{}, err
	}
	account, err := op.accounts.FindByID(ctx, actor.AccountID)
	if err != nil {
		return accountdto.Account{}, err
	}
	return op.mapper.ToResponse(account), nil
}
