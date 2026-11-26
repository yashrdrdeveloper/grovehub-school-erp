package impl

import (
	"context"
	"time"

	dto "grovehub/services/identity/internal/app/dto/account"
	"grovehub/services/identity/internal/app/entity"
	tokens "grovehub/services/identity/internal/app/repository/passwordtoken"
	"grovehub/services/identity/internal/app/services/account"
	"grovehub/services/identity/internal/app/services/account/operations/setup"
	"grovehub/services/identity/internal/app/services/account/operations/setup/eligibility"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	"grovehub/services/identity/internal/app/services/registration/operations/register"
	"grovehub/services/identity/internal/app/services/registration/operations/resend"
)

type Onboarding struct {
	account.Service
	gate *cutover.Gate
	registration register.Operation
	resend resend.Operation
	setup setup.Operation
	proofs *eligibility.ProofGuard
}

var _ account.Service = (*Onboarding)(nil)

func NewOnboarding(base account.Service, gate *cutover.Gate, registration register.Operation, resend resend.Operation, setup setup.Operation, proofs *eligibility.ProofGuard) account.Service {
	return &Onboarding{Service: base, gate: gate, registration: registration, resend: resend, setup: setup, proofs: proofs}
}

func (service *Onboarding) Create(ctx context.Context, request dto.Create, key string) (dto.CreateResult, error) {
	if err := service.gate.Require(ctx); err != nil {
		return dto.CreateResult{}, err
	}
	return service.registration.Execute(ctx, request, key)
}

func (service *Onboarding) Resend(ctx context.Context, id string) error {
	if err := service.gate.Require(ctx); err != nil {
		return err
	}
	return service.resend.Execute(ctx, id)
}

func (service *Onboarding) Setup(ctx context.Context, grant string, input dto.Setup) (dto.SetupResult, error) {
	// Release the read lock before reserving the independently committed attempt.
	if err := service.proofs.Run(ctx, grant, entity.SetupGrant, func(context.Context, entity.PasswordToken) error { return nil }); err != nil {
		return dto.SetupResult{}, err
	}
	return service.setup.Execute(ctx, grant, input)
}

func (service *Onboarding) SetupStatus(ctx context.Context, grant string) (dto.SetupStatusResult, error) {
	var result dto.SetupStatusResult
	err := service.proofs.Run(ctx, grant, entity.SetupGrant, func(_ context.Context, proof entity.PasswordToken) error {
		remaining := int(time.Until(proof.ExpiresAt).Seconds())
		if remaining < 1 {
			return tokens.ErrInvalidToken
		}
		result = dto.SetupStatusResult{Verified: true, ExpiresIn: remaining}
		return nil
	})
	return result, err
}
