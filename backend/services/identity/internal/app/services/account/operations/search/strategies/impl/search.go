package impl

import (
	"context"
	"time"

	chain "grovehub/libs/core/pipeline/service"
	coresearch "grovehub/libs/core/search"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	mapper "grovehub/services/identity/internal/app/mappers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/services/security/principal"
	"grovehub/services/identity/internal/app/services/account/operations/search/flows"
	strategies "grovehub/services/identity/internal/app/services/account/operations/search/strategies"
	validator "grovehub/services/identity/internal/app/validators/account"
)

type Search struct {
	accounts accountrepo.Repository
	mapper mapper.Mapper
	validator validator.Validator
}

var _ strategies.SearchStrategy = (*Search)(nil)

func New(accounts accountrepo.Repository, accountMapper mapper.Mapper, accountValidator validator.Validator) strategies.SearchStrategy {
	return &Search{accounts: accounts, mapper: accountMapper, validator: accountValidator}
}

func (strategy *Search) Prepare(ctx context.Context, state *flows.Search) chain.Outcome[coresearch.Result[accountdto.Account]] {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return chain.Fail[coresearch.Result[accountdto.Account]](err)
	}
	query, err := strategy.validator.Search(state.Request)
	if err != nil {
		return chain.Fail[coresearch.Result[accountdto.Account]](err)
	}
	state.TenantID = actor.TenantID
	state.Query = query
	return chain.Next[coresearch.Result[accountdto.Account]]()
}

func (strategy *Search) Query(ctx context.Context, state *flows.Search) chain.Outcome[coresearch.Result[accountdto.Account]] {
	result, err := strategy.accounts.Search(ctx, state.TenantID, state.Query)
	if err != nil {
		return chain.Fail[coresearch.Result[accountdto.Account]](err)
	}
	state.Accounts = result
	return chain.Next[coresearch.Result[accountdto.Account]]()
}

func (strategy *Search) Hydrate(ctx context.Context, state *flows.Search) chain.Outcome[coresearch.Result[accountdto.Account]] {
	result := coresearch.Result[accountdto.Account]{Items: make([]accountdto.Account, 0, len(state.Accounts.Items))}
	for _, account := range state.Accounts.Items {
		result.Items = append(result.Items, strategy.mapper.ToResponse(account))
	}
	result.PageInfo = coresearch.PageResponse{
		Rows: len(result.Items),
		HasNext: state.Accounts.HasNext,
	}

	if state.Accounts.HasNext && len(state.Accounts.Items) > 0 {
		last := state.Accounts.Items[len(state.Accounts.Items)-1]
		value := last.CreatedAt.UTC().Format(time.RFC3339Nano)
		switch state.Query.Sort.Field {
		case accountdto.FieldAccountID:
			value = last.AccountID
		case accountdto.FieldUpdatedAt:
			value = last.UpdatedAt.UTC().Format(time.RFC3339Nano)
		case accountdto.FieldStatus:
			value = string(last.Status)
		case accountdto.FieldUsername:
			if last.Username != nil {
				value = *last.Username
			} else {
				value = ""
			}
		}
		cursor, err := coresearch.EncodeCursor(coresearch.CursorPayload{
			Field: state.Query.Sort.Field,
			Direction: state.Query.Sort.Direction,
			Value: value,
			CursorId: last.AccountID,
		})
		if err != nil {
			return chain.Fail[coresearch.Result[accountdto.Account]](err)
		}
		result.PageInfo.NextCursor = cursor
	}
	return chain.Return(result, false)
}
