package impl

import (
	"context"
	"errors"

	coreerrors "grovehub/libs/core/errors"
	coreidempotency "grovehub/libs/core/idempotency"
	foundationclient "grovehub/services/foundation/client"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
	identityresponse "grovehub/services/identity/dto/response"
	identitymappers "grovehub/services/identity/mappers"
	identityservices "grovehub/services/identity/services"
	servicehelpers "grovehub/services/identity/services/helpers"
	identityvalidators "grovehub/services/identity/validators"
)

type AuthService struct {
	foundationClient          foundationclient.Client
	accountService            identityservices.AccountService
	passwordCredentialService identityservices.PasswordCredentialService
	mapper                    identitymappers.AuthMapper
	validator                 identityvalidators.AuthValidator
}

func NewAuthService(
	foundationClient foundationclient.Client,
	accountService identityservices.AccountService,
	passwordCredentialService identityservices.PasswordCredentialService,
	mapper identitymappers.AuthMapper,
	validator identityvalidators.AuthValidator,
) identityservices.AuthService {
	return &AuthService{
		foundationClient:          foundationClient,
		accountService:            accountService,
		passwordCredentialService: passwordCredentialService,
		mapper:                    mapper,
		validator:                 validator,
	}
}

func (service *AuthService) Register(
	ctx context.Context,
	request identityrequest.Register,
) (identityresponse.Register, error) {
	idempotencyKey, err := coreidempotency.ValidateKey(request.IdempotencyKey)
	if err != nil {
		return identityresponse.Register{}, coreerrors.NewError(identityerror.InvalidRegistrationIdempotencyKey)
	}
	request.IdempotencyKey = idempotencyKey

	validated, err := service.validator.ValidateRegister(request)
	if err != nil {
		return identityresponse.Register{}, err
	}
	validated.IdempotencyKey = idempotencyKey
	requestHash, err := coreidempotency.Fingerprint(validated)
	if err != nil {
		return identityresponse.Register{}, coreerrors.Wrap(identityerror.RegistrationFailed, err)
	}

	if account, readErr := service.accountService.GetByRegistrationIdempotencyKey(ctx, idempotencyKey); readErr == nil {
		if account.GetRegistrationRequestHash() != requestHash {
			return identityresponse.Register{}, coreerrors.NewError(identityerror.RegistrationIdempotencyConflict)
		}
		if err := service.passwordCredentialService.Set(ctx, string(account.GetAccountId()), validated.Password); err != nil {
			return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(err)
		}
		return service.mapper.ToRegisterResponse(account), nil
	} else if !isAccountNotFound(readErr) {
		return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(readErr)
	}

	person, err := service.resolveRegisterPerson(ctx, validated)
	if err != nil {
		return identityresponse.Register{}, err
	}

	if _, readErr := service.accountService.GetByPersonGid(ctx, person.GID); readErr == nil {
		return identityresponse.Register{}, coreerrors.NewError(identityerror.RegistrationAlreadyExists)
	} else if !isAccountNotFound(readErr) {
		return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(readErr)
	}

	person, err = service.foundationClient.UpdatePerson(
		ctx,
		foundationrequest.PersonByGID{GID: person.GID},
		servicehelpers.BuildRegisterPersonUpdate(validated),
	)
	if err != nil {
		return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(err)
	}

	account, created, err := service.accountService.Create(ctx, identityservices.CreateAccountInput{
		PersonGid:                  person.GID,
		Email:                      validated.Email,
		Phone:                      validated.PhoneNumber,
		RegistrationIdempotencyKey: idempotencyKey,
		RegistrationRequestHash:    requestHash,
	})
	if err != nil {
		return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(err)
	}

	accountId := string(account.GetAccountId())
	if err := service.passwordCredentialService.Set(ctx, accountId, validated.Password); err != nil {
		if !created {
			return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(err)
		}
		cleanupErr := service.accountService.PermanentlyDelete(ctx, accountId)
		if cleanupErr != nil {
			return identityresponse.Register{}, coreerrors.Wrap(
				identityerror.RegistrationFailed,
				errors.Join(err, cleanupErr),
			)
		}
		return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(err)
	}

	return service.mapper.ToRegisterResponse(account), nil
}

func (service *AuthService) resolveRegisterPerson(
	ctx context.Context,
	request identityrequest.Register,
) (foundationresponse.Person, error) {
	if request.Gid != "" {
		person, err := service.foundationClient.GetPersonByGID(
			ctx,
			foundationrequest.PersonByGID{GID: request.Gid},
		)
		if err != nil {
			return foundationresponse.Person{}, servicehelpers.MapAuthDependencyError(err)
		}
		if !servicehelpers.PersonHasEmail(person, request.Email) {
			return foundationresponse.Person{}, coreerrors.NewError(identityerror.InvalidRegistration)
		}
		return person, nil
	}

	result, err := service.foundationClient.SearchPersons(
		ctx,
		servicehelpers.RegisterEmailSearch(request.Email),
	)
	if err != nil {
		return foundationresponse.Person{}, servicehelpers.MapAuthDependencyError(err)
	}
	if len(result.Items) > 1 {
		return foundationresponse.Person{}, coreerrors.NewError(identityerror.RegistrationFailed)
	}
	if len(result.Items) == 1 {
		return result.Items[0], nil
	}

	person, _, err := service.foundationClient.CreatePerson(
		ctx,
		servicehelpers.BuildRegisterPersonRequest(request),
	)
	if err != nil {
		return foundationresponse.Person{}, servicehelpers.MapAuthDependencyError(err)
	}
	if !servicehelpers.PersonHasEmail(person, request.Email) {
		return foundationresponse.Person{}, coreerrors.NewError(identityerror.RegistrationAlreadyExists)
	}
	return person, nil
}

func isAccountNotFound(err error) bool {
	var applicationError *coreerrors.Error
	return errors.As(err, &applicationError) &&
		applicationError.Definition().Code() == identityerror.AccountNotFound.Code()
}
