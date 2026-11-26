package register

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	foundation "grovehub/services/foundation/client/person"
	"grovehub/services/identity/internal/app/constants/permission"
	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	producer "grovehub/services/identity/internal/app/events/producers/registration"
	repo "grovehub/services/identity/internal/app/repository/account"
	"grovehub/services/identity/internal/app/repository/authority"
	"grovehub/services/identity/internal/app/repository/rbacmigration"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/registration/operations/register/flows"
	"grovehub/services/identity/internal/app/services/registration/operations/register/helpers"
	"grovehub/services/identity/internal/app/services/security/principal"
	validator "grovehub/services/identity/internal/app/validators/account"
)

// Normalized is staged until invitation issuance and setup are authority-aware.
// It never calls the legacy account writer or assigns privileges.
type Normalized struct {
	writer repo.RegistrationWriter
	foundation foundation.Client
	validator validator.Validator
	locker authority.Locker
	readiness rbacmigration.Readiness
	authorizer authorization.Service
	events producer.Producer
}

var _ Operation = (*Normalized)(nil)

func NewNormalized(writer repo.RegistrationWriter, foundationClient foundation.Client, validator validator.Validator, locker authority.Locker, readiness rbacmigration.Readiness, authorizer authorization.Service, events producer.Producer) Operation {
	return &Normalized{writer: writer, foundation: foundationClient, validator: validator, locker: locker, readiness: readiness, authorizer: authorizer, events: events}
}

func (operation *Normalized) Execute(ctx context.Context, request dto.Create, key string) (dto.CreateResult, error) {
	actor, err := principal.RequireTenant(ctx)
	if err != nil {
		return dto.CreateResult{}, err
	}
	if operation.writer == nil || operation.foundation == nil || operation.validator == nil || operation.locker == nil || operation.readiness == nil || operation.authorizer == nil || operation.events == nil {
		return dto.CreateResult{}, authorization.ErrUnavailable
	}
	request, key, digest, err := helpers.NormalizePending(request, key, operation.validator)
	if err != nil {
		return dto.CreateResult{}, err
	}
	var response dto.CreateResult
	err = operation.locker.Run(ctx, actor.TenantID, func(tx context.Context) error {
		if err := operation.authorizer.Require(tx, permission.AccountCreate); err != nil {
			return err
		}
		if err := operation.readiness.Require(tx, actor.TenantID); err != nil {
			return err
		}
		receipt, err := operation.writer.FindRegistration(tx, actor.TenantID, actor.AccountID, key)
		if err == nil {
			if receipt.Digest != digest {
				return repo.ErrConflict
			}
			response = pendingResult(receipt.AccountID, receipt.PersonGID, receipt.Email)
			return nil
		}
		if !errors.Is(err, repo.ErrNotFound) {
			return err
		}
		flow := &flows.Register{Request: request}
		resolver := &RegisterOperation{foundation: operation.foundation}
		if err := resolver.resolve(tx, flow); err != nil {
			return err
		}
		personID, err := uuid.Parse(flow.Person.GID)
		if err != nil || personID == uuid.Nil {
			return repo.ErrConflict
		}
		registration := repo.Registration{TenantID: actor.TenantID, ActorID: actor.AccountID, AccountID: uuid.NewString(), PersonGID: personID.String(), Email: request.Email, Key: key, Digest: digest, EventID: uuid.NewString(), At: time.Now().UTC()}
		if request.PhoneNumber != "" {
			registration.Phone = &request.PhoneNumber
		}
		if err := operation.writer.CreatePending(tx, registration); err != nil {
			return err
		}
		if err := operation.events.Enqueue(tx, registration); err != nil {
			return err
		}
		response = pendingResult(registration.AccountID, registration.PersonGID, registration.Email)
		return nil
	})
	if err != nil {
		return dto.CreateResult{}, err
	}
	return response, nil
}

// A replay returns the creation receipt, not current lifecycle or mail state.
func pendingResult(accountID string, personID string, email string) dto.CreateResult {
	return dto.CreateResult{AccountID: accountID, PersonGID: personID, Email: email, Status: string(entity.PendingSetup), Invitation: dto.InvitationStatus{Status: "AWAITING_ROLE_ASSIGNMENT"}}
}
