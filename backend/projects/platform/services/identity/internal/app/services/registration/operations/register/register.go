package register

import (
	"context"
	"errors"
	"strings"
	"time"

	coredatabase "grovehub/libs/core/database"
	headercontext "grovehub/libs/core/headercontext"
	coreid "grovehub/libs/core/id"
	idempotency "grovehub/libs/core/idempotency"
	corepipeline "grovehub/libs/core/pipeline"
	coresearch "grovehub/libs/core/search"
	foundation "grovehub/services/foundation/client/person"
	accountdto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	accountproducer "grovehub/services/identity/internal/app/events/producers/account"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accounthelpers "grovehub/services/identity/internal/app/services/account/helpers"
	"grovehub/services/identity/internal/app/services/security/principal"
	"grovehub/services/identity/internal/app/services/registration/operations/register/flows"
	registration "grovehub/services/identity/internal/app/services/registration"
	validator "grovehub/services/identity/internal/app/validators/account"
)

type Operation interface {
	Execute(ctx context.Context, request accountdto.Create, key string) (accountdto.CreateResult, error)
}

type RegisterOperation struct {
	accounts accountrepo.Repository
	foundation foundation.Client
	validator validator.Validator
	invitations registration.InvitationSender
	unit coredatabase.UnitOfWork
	events accountproducer.ChangedProducer
}

func New(accounts accountrepo.Repository, foundationClient foundation.Client, validator validator.Validator, invitations registration.InvitationSender, unit coredatabase.UnitOfWork, events accountproducer.ChangedProducer) Operation {
	return &RegisterOperation{
		accounts: accounts,
		foundation: foundationClient,
		validator: validator,
		invitations: invitations,
		unit: unit,
		events: events,
	}
}

// @ChainOfResponsibility
func (operation *RegisterOperation) Execute(ctx context.Context, request accountdto.Create, key string) (accountdto.CreateResult, error) {
	flow := &flows.Register{Request: request, IdempotencyKey: key}

	wrap := func(step func(context.Context, *flows.Register) error) corepipeline.Handler[flows.Register, accountdto.CreateResult] {
		return func(ctx context.Context, flow *flows.Register) corepipeline.Outcome[accountdto.CreateResult] {
			if err := step(ctx, flow); err != nil {
				return corepipeline.Fail[accountdto.CreateResult](err)
			}
			return corepipeline.Next[accountdto.CreateResult]()
		}
	}

	var response accountdto.CreateResult
	err := operation.unit.Run(ctx, func(transactionCtx context.Context) error {
		var runErr error
		response, _, runErr = corepipeline.RunService(
			transactionCtx,
			flow,
			wrap(operation.authorize),
			wrap(operation.validate),
			func(ctx context.Context, flow *flows.Register) corepipeline.Outcome[accountdto.CreateResult] {
				if err := operation.replay(ctx, flow); err != nil {
					return corepipeline.Fail[accountdto.CreateResult](err)
				}
				if flow.Done {
					return corepipeline.Return(flow.Response, false)
				}
				return corepipeline.Next[accountdto.CreateResult]()
			},
			wrap(operation.resolve),
			func(ctx context.Context, flow *flows.Register) corepipeline.Outcome[accountdto.CreateResult] {
				if err := operation.persist(ctx, flow); err != nil {
					return corepipeline.Fail[accountdto.CreateResult](err)
				}
				return corepipeline.Return(flow.Response, true)
			},
		)
		return runErr
	})
	return response, err
}

func (operation *RegisterOperation) authorize(ctx context.Context, flow *flows.Register) error {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return err
	}
	if operation.invitations == nil {
		return errors.New("invitation service unavailable")
	}
	flow.Actor = actor
	return nil
}

func (operation *RegisterOperation) validate(ctx context.Context, flow *flows.Register) error {
	request, err := operation.validator.NormalizeCreate(flow.Request)
	if err != nil {
		return err
	}
	flow.Request = request

	if _, err := idempotency.ValidateKey(flow.IdempotencyKey); err != nil {
		return err
	}
	flow.RequestHash, err = idempotency.Fingerprint(flow.Request)
	return err
}

func (operation *RegisterOperation) replay(ctx context.Context, flow *flows.Register) error {
	current, err := operation.accounts.FindByIdempotencyKey(ctx, flow.IdempotencyKey)
	if errors.Is(err, accountrepo.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	if current.RegistrationRequestHash == nil ||
		*current.RegistrationRequestHash != flow.RequestHash {
		return errors.New("idempotency key conflict")
	}
	if current.TenantID == nil || *current.TenantID != flow.Actor.TenantID {
		return errors.New("forbidden tenant")
	}

	flow.Account = current
	flow.Response = mapCreated(current)
	flow.Done = true
	return nil
}

func (operation *RegisterOperation) resolve(ctx context.Context, flow *flows.Register) error {
	request := flow.Request
	searchRequest := coresearch.Request{
		Filters: []coresearch.Filter{
			{
				Field: "contacts.email",
				Operator: coresearch.OperatorEquals,
				Values: []string{request.Email},
			},
		},
		PageInfo: coresearch.PageRequest{Rows: 2},
	}

	result, err := operation.foundation.SearchPersons(ctx, foundation.SearchPersonsRequest{
		RequestBody: searchRequest,
	})
	if err != nil {
		return err
	}
	if len(result.Items) > 1 {
		return errors.New("ambiguous person")
	}
	if len(result.Items) == 1 {
		found := false
		for _, email := range result.Items[0].Contacts.Emails {
			if strings.EqualFold(email.Value, request.Email) {
				found = true
				break
			}
		}
		if !found {
			return errors.New("Foundation contact mismatch")
		}
		flow.Person = result.Items[0]
		return nil
	}

	middleName := ""
	if request.Person.MiddleName != nil {
		middleName = *request.Person.MiddleName
	}
	newPerson := foundation.CreatePerson{
		PersonNames: foundation.PersonNames{
			FirstName: request.Person.FirstName,
			MiddleName: middleName,
			LastName: request.Person.LastName,
		},
		Contacts: foundation.Contacts{
			Emails: []foundation.ContactValue{
				{Value: request.Email, IsPrimary: true},
			},
		},
	}
	if request.PhoneNumber != "" {
		newPerson.Contacts.Phones = []foundation.ContactValue{
			{Value: request.PhoneNumber, IsPrimary: true},
		}
	}

	flow.Person, _, err = operation.foundation.CreatePerson(
		headercontext.WithoutIdempotencyKey(ctx),
		foundation.CreatePersonRequest{RequestBody: newPerson},
	)
	return err
}

func (operation *RegisterOperation) persist(ctx context.Context, flow *flows.Register) error {
	_, err := operation.accounts.FindByPersonGID(ctx, flow.Person.GID)
	if err == nil {
		return accountrepo.ErrConflict
	}
	if !errors.Is(err, accountrepo.ErrNotFound) {
		return err
	}

	accountID, err := coreid.NewUUIDv7()
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	accountEntity := entity.Account{
		AccountID: accountID,
		PersonGID: flow.Person.GID,
		TenantID: &flow.Actor.TenantID,
		Email: flow.Request.Email,
		Status: entity.PendingSetup,
		Version: 1,
		CreatedAt: now,
		UpdatedAt: now,
		RegistrationIdempotencyKey: &flow.IdempotencyKey,
		RegistrationRequestHash: &flow.RequestHash,
	}
	if flow.Request.PhoneNumber != "" {
		accountEntity.Phone = &flow.Request.PhoneNumber
	}
	if flow.Request.RoleAssignment != nil {
		if err := accounthelpers.ApplyRoleAssignment(&accountEntity, *flow.Request.RoleAssignment); err != nil {
			return err
		}
	}
	accountEntity.SearchDetails = strings.ToLower(accountEntity.Email)

	if err := operation.accounts.Create(ctx, accountEntity); err != nil {
		return err
	}
	if err := operation.invitations.Queue(ctx, accountEntity, entity.AccountSetup, flow.Request.Person.FirstName); err != nil {
		return err
	}
	if err := operation.events.Enqueue(ctx, "CREATED", nil, &accountEntity); err != nil {
		return err
	}

	flow.Account = accountEntity
	flow.Response = mapCreated(accountEntity)
	return nil
}

func mapCreated(account entity.Account) accountdto.CreateResult {
	return accountdto.CreateResult{
		AccountID: account.AccountID,
		PersonGID: account.PersonGID,
		Email: account.Email,
		Status: string(account.Status),
		Invitation: accountdto.InvitationStatus{Status: "QUEUED"},
	}
}
