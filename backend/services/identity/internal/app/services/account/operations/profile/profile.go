package profile

import (
	"context"

	coresearch "grovehub/libs/core/search"
	foundation "grovehub/services/foundation/client/person"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	accountservice "grovehub/services/identity/internal/app/services/account"
	"grovehub/services/identity/internal/app/services/account/operations/profile/helpers"
)

type Service interface {
	Get(ctx context.Context, accountID string) (accountdto.Profile, error)
	Update(ctx context.Context, accountID string, input accountdto.UpdateProfile) (accountdto.Profile, error)
	Enrich(ctx context.Context, accounts []accountdto.Account) ([]accountdto.Account, error)
}

type Operation struct {
	accounts   accountservice.Service
	foundation foundation.Client
}

func New(accounts accountservice.Service, client foundation.Client) Service {
	return &Operation{accounts: accounts, foundation: client}
}

func (operation *Operation) Get(ctx context.Context, accountID string) (accountdto.Profile, error) {
	account, err := operation.accounts.Get(ctx, accountID)
	if err != nil {
		return accountdto.Profile{}, err
	}
	person, err := operation.foundation.GetPersonByGID(ctx, foundation.GetPersonByGIDRequest{
		Params: foundation.PersonByGID{GID: account.PersonGID},
	})
	if err != nil {
		return accountdto.Profile{}, err
	}
	return helpers.MapPerson(person), nil
}

func (operation *Operation) Update(ctx context.Context, accountID string, input accountdto.UpdateProfile) (accountdto.Profile, error) {
	account, err := operation.accounts.Get(ctx, accountID)
	if err != nil {
		return accountdto.Profile{}, err
	}
	patch, err := helpers.BuildPatch(input)
	if err != nil {
		return accountdto.Profile{}, err
	}
	person, err := operation.foundation.UpdatePerson(ctx, foundation.UpdatePersonRequest{
		Params: foundation.PersonByGID{GID: account.PersonGID},
		RequestBody: patch,
	})
	if err != nil {
		return accountdto.Profile{}, err
	}
	return helpers.MapPerson(person), nil
}

// Enrich receives only accounts already returned by tenant-scoped account search.
// It batches Foundation lookups using those trusted person GIDs, never caller-provided GIDs.
func (operation *Operation) Enrich(ctx context.Context, accounts []accountdto.Account) ([]accountdto.Account, error) {
	profiles := make(map[string]accountdto.Profile, len(accounts))
	for start := 0; start < len(accounts); start += 25 {
		end := start + 25
		if end > len(accounts) {
			end = len(accounts)
		}
		gids := make([]string, 0, end-start)
		for _, account := range accounts[start:end] {
			gids = append(gids, account.PersonGID)
		}
		result, err := operation.foundation.SearchPersons(ctx, foundation.SearchPersonsRequest{
			RequestBody: coresearch.Request{
				Filters: []coresearch.Filter{{Field: "gid", Operator: coresearch.OperatorIn, Values: gids}},
				PageInfo: coresearch.PageRequest{Rows: len(gids)},
			},
		})
		if err != nil {
			return nil, err
		}
		for _, person := range result.Items {
			profiles[person.GID] = helpers.MapPerson(person)
		}
	}
	for index := range accounts {
		if person, ok := profiles[accounts[index].PersonGID]; ok {
			accounts[index].Profile = &person
		}
	}
	return accounts, nil
}
