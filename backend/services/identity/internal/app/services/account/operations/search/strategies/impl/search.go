package impl

import (
	"context"
	"time"
	"strings"
	"grovehub/services/identity/internal/app/entity"

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
		TotalRecords: state.Accounts.TotalRecords,
		Rows: len(result.Items),
		HasNext: state.Accounts.HasNext,
	}

	if state.Accounts.HasNext && len(state.Accounts.Items) > 0 {
		last := state.Accounts.Items[len(state.Accounts.Items)-1]
		payload := coresearch.CursorPayload{
			Field: state.Query.Sort.Field,
			Direction: state.Query.Sort.Direction,
			Value: accountSortValue(last, state.Query.Sort.Field),
			CursorId: last.AccountID,
		}
		for _, rule := range state.Query.SecondarySorts {
			payload.SecondarySorts = append(payload.SecondarySorts, coresearch.Sort{Field: rule.Field, Direction: rule.Direction})
			payload.SecondaryValues = append(payload.SecondaryValues, accountSortValue(last, rule.Field))
		}
		cursor, err := coresearch.EncodeCursor(payload)
		if err != nil {
			return chain.Fail[coresearch.Result[accountdto.Account]](err)
		}
		result.PageInfo.NextCursor = cursor
	}
	return chain.Return(result, false)
}

func accountSortValue(account entity.Account, field string) string {
	switch field {
	case accountdto.FieldAccountID:
		return account.AccountID
	case accountdto.FieldUpdatedAt:
		return account.UpdatedAt.UTC().Format(time.RFC3339Nano)
	case accountdto.FieldStatus:
		return string(account.Status)
	case accountdto.FieldUsername:
		if account.Username != nil { return strings.ToLower(*account.Username) }
		return ""
	default:
		return account.CreatedAt.UTC().Format(time.RFC3339Nano)
	}
}
