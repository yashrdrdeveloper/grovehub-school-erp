package search

import (
	"context"
	"time"

	coresearch "grovehub/libs/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	mapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/security/principal"
	validator "grovehub/services/identity/internal/app/validators/account"
)

type Operation struct {
	accounts accountrepo.Repository
	mapper mapper.Mapper
	validator validator.Validator
}

func New(accounts accountrepo.Repository, accountMapper mapper.Mapper, accountValidator validator.Validator) *Operation {
	return &Operation{
		accounts: accounts,
		mapper: accountMapper,
		validator: accountValidator,
	}
}

func (operation *Operation) Execute(ctx context.Context, request coresearch.Request) (coresearch.Result[accountdto.Account], error) {
	actor, err := principal.RequireAdmin(ctx)
	if err != nil {
		return coresearch.Result[accountdto.Account]{}, err
	}
	query, err := operation.validator.Search(request)
	if err != nil {
		return coresearch.Result[accountdto.Account]{}, err
	}
	result, err := operation.accounts.Search(ctx, actor.TenantID, query)
	if err != nil {
		return coresearch.Result[accountdto.Account]{}, err
	}

	response := coresearch.Result[accountdto.Account]{
		Items: make([]accountdto.Account, 0, len(result.Items)),
	}
	for _, account := range result.Items {
		response.Items = append(response.Items, operation.mapper.ToResponse(account))
	}
	response.PageInfo = coresearch.PageResponse{
		Rows: len(response.Items),
		HasNext: result.HasNext,
	}

	if result.HasNext && len(result.Items) != 0 {
		last := result.Items[len(result.Items)-1]
		value := last.CreatedAt.UTC().Format(time.RFC3339Nano)
		switch query.Sort.Field {
		case "accountId":
			value = last.AccountID
		case "updatedAt":
			value = last.UpdatedAt.UTC().Format(time.RFC3339Nano)
		case "status":
			value = string(last.Status)
		case "username":
			if last.Username != nil {
				value = *last.Username
			} else {
				value = ""
			}
		}
		response.PageInfo.NextCursor, err = coresearch.EncodeCursor(coresearch.CursorPayload{
			Field: query.Sort.Field,
			Direction: query.Sort.Direction,
			Value: value,
			CursorId: last.AccountID,
		})
	}
	return response, err
}
