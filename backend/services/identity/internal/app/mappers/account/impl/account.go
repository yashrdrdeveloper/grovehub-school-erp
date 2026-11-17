package impl

import (
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	accountmapper "grovehub/services/identity/internal/app/mappers/account"
)

type Mapper struct{}

var _ accountmapper.Mapper = (*Mapper)(nil)

func New() accountmapper.Mapper {
	return &Mapper{}
}

func (mapper *Mapper) ToResponse(account entity.Account) accountdto.Account {
	return accountdto.Account{
		AccountID: account.AccountID,
		PersonGID: account.PersonGID,
		Username: account.Username,
		Email: account.Email,
		Status: string(account.Status),
		Version: account.Version,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
