package impl

import (
	"context"
	"errors"

 "gorm.io/gorm"

 coretransaction "grovehub/libs/core/database/transaction"

	coreerrors "grovehub/libs/core/errors"
	coreidempotency "grovehub/libs/core/idempotency"
    coreevents "grovehub/libs/core/events"
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
 db *gorm.DB
}

func NewAuthService(
	foundationClient foundationclient.Client,
	accountService identityservices.AccountService,
	passwordCredentialService identityservices.PasswordCredentialService,
	mapper identitymappers.AuthMapper,
	validator identityvalidators.AuthValidator,
 db *gorm.DB,
) identityservices.AuthService {
	return &AuthService{
		foundationClient:          foundationClient,
		accountService:            accountService,
		passwordCredentialService: passwordCredentialService,
		mapper:                    mapper,
		validator:                 validator,
 db: db,
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
    ctx = coreevents.WithCausationID(ctx, idempotencyKey)

	validated, err := service.validator.ValidateRegister(request)
	if err != nil {
		return identityresponse.Register{}, err
	}
	validated.IdempotencyKey = idempotencyKey
	requestHash, err := coreidempotency.Fingerprint(validated)
	if err != nil {
		return identityresponse.Register{}, coreerrors.Wrap(identityerror.RegistrationFailed, err)
	}

    var registration identityresponse.Register
    err = coretransaction.Run(ctx, service.db, func(txCtx context.Context) error {
        var operationErr error
        registration, operationErr = service.register(txCtx, validated, idempotencyKey, requestHash)
        return operationErr
    })
    if err != nil { return identityresponse.Register{}, err }
    return registration, nil
}

func (service *AuthService) register(
    ctx context.Context,
    validated identityrequest.Register,
    idempotencyKey string,
    requestHash string,
) (identityresponse.Register, error) {
	if account, readErr := service.accountService.GetByRegistrationIdempotencyKey(ctx, idempotencyKey); readErr == nil {
		if servicehelpers.NullableString(account.RegistrationRequestHash) != requestHash {
			return identityresponse.Register{}, coreerrors.NewError(identityerror.RegistrationIdempotencyConflict)
		}
		if err := service.passwordCredentialService.Set(ctx, string(account.AccountId), validated.Password); err != nil {
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

	account, _, err := service.accountService.Create(ctx, identityservices.CreateAccountInput{
		PersonGid:                  person.GID,
		Email:                      validated.Email,
		Phone:                      validated.PhoneNumber,
		RegistrationIdempotencyKey: idempotencyKey,
		RegistrationRequestHash:    requestHash,
	})
	if err != nil {
		return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(err)
	}

	accountId := string(account.AccountId)
	if err := service.passwordCredentialService.Set(ctx, accountId, validated.Password); err != nil {
		return identityresponse.Register{}, servicehelpers.MapAuthDependencyError(err)
	}

	return service.mapper.ToRegisterResponse(account), nil
}

func (service *AuthService) resolveRegisterPerson(
	ctx context.Context,
	request identityrequest.Register,
) (foundationresponse.Person, error) {
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
