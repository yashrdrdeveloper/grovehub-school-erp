package impl

import (
	"encoding/json"
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
	secondary := make([]string, 0)
	if len(account.SecondaryRoleIDs) > 0 {
		_ = json.Unmarshal(account.SecondaryRoleIDs, &secondary)
	}
	if secondary == nil {
		secondary = make([]string, 0)
	}
	return accountdto.Account{
		AccountID: account.AccountID,
		PersonGID: account.PersonGID,
		Username: account.Username,
		PrimaryRoleID: account.PrimaryRoleID,
		SecondaryRoleIDs: secondary,
		Email: account.Email,
		Status: string(account.Status),
		Version: account.Version,
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}
}
