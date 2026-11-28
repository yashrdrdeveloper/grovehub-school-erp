package impl

import (
	"context"

	dto "grovehub/services/identity/internal/app/dto/authentication"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/account/operations/setup/eligibility"
	"grovehub/services/identity/internal/app/services/authentication"
)

type Onboarding struct {
	authentication.Service
	proofs *eligibility.ProofGuard
}

func NewOnboarding(base authentication.Service, proofs *eligibility.ProofGuard) authentication.Service {
	return &Onboarding{Service: base, proofs: proofs}
}

func (service *Onboarding) Exchange(ctx context.Context, secret string, purpose string) (authentication.Grant, error) {
	if purpose != "invitation" {
		return service.Service.Exchange(ctx, secret, purpose)
	}
	var result authentication.Grant
	err := service.proofs.Run(ctx, secret, entity.AccountSetup, func(tx context.Context, _ entity.PasswordToken) error {
		var err error
		result, err = service.Service.Exchange(tx, secret, purpose)
		return err
	})
	if err != nil {
		return authentication.Grant{}, err
	}
	return result, nil
}

func (service *Onboarding) InspectProof(ctx context.Context, secret string, purpose string) (dto.GrantResult, error) {
	if purpose != "invitation" {
		return service.Service.InspectProof(ctx, secret, purpose)
	}
	var result dto.GrantResult
	err := service.proofs.Run(ctx, secret, entity.AccountSetup, func(tx context.Context, _ entity.PasswordToken) error {
		var err error
		result, err = service.Service.InspectProof(tx, secret, purpose)
		return err
	})
	return result, err
}
