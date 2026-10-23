package mappers

import (
	identityresponse "grovehub/services/identity/dto/response"
	"grovehub/services/identity/entity"
)

func ToBootstrapAccountResponse(
	account entity.Account,
	identifier entity.LoginIdentifier,
) identityresponse.BootstrapAccount {
	return identityresponse.BootstrapAccount{
		AccountId: string(account.GetAccountId()),
		GID:       string(account.GetPersonGid()),
		Status:    string(account.GetStatus()),
		Identifier: identityresponse.LoginIdentifier{
			Type:       string(identifier.GetType()),
			Value:      identifier.GetValue(),
			IsVerified: identifier.GetIsVerified(),
		},
		CreatedAt: account.GetCreatedAt(),
	}
}
