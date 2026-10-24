package impl

import (
	"context"
	"errors"

	coreerrors "grovehub/libs/core/errors"
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

type RegistrationService struct {
	foundationClient          foundationclient.Client
	accountService            identityservices.AccountService
	passwordCredentialService identityservices.PasswordCredentialService
	mapper                    identitymappers.RegistrationMapper
	validator                 identityvalidators.RegistrationValidator
}

func NewRegistrationService(
	foundationClient foundationclient.Client,
	accountService identityservices.AccountService,
	passwordCredentialService identityservices.PasswordCredentialService,
	mapper identitymappers.RegistrationMapper,
	validator identityvalidators.RegistrationValidator,
) identityservices.RegistrationService {
	return &RegistrationService{
		foundationClient:          foundationClient,
		accountService:            accountService,
		passwordCredentialService: passwordCredentialService,
		mapper:                    mapper,
		validator:                 validator,
	}
}

func (service *RegistrationService) Register(
	ctx context.Context,
	request identityrequest.Register,
) (identityresponse.Registration, error) {
	validated, err := service.validator.Validate(request)
	if err != nil {
		return identityresponse.Registration{}, err
	}

	person, err := service.resolvePerson(ctx, validated)
	if err != nil {
		return identityresponse.Registration{}, err
	}

	account, err := service.accountService.Create(ctx, identityservices.CreateAccountInput{
		PersonGid: person.GID,
		Email:     validated.Email,
		Phone:     validated.PhoneNumber,
	})
	if err != nil {
		return identityresponse.Registration{}, servicehelpers.MapRegistrationDependencyError(err)
	}

	accountId := string(account.GetAccountId())
	if err := service.passwordCredentialService.Set(ctx, accountId, validated.Password); err != nil {
		cleanupErr := service.accountService.PermanentlyDelete(ctx, accountId)
		if cleanupErr != nil {
			return identityresponse.Registration{}, coreerrors.Wrap(
				identityerror.RegistrationFailed,
				errors.Join(err, cleanupErr),
			)
		}
		return identityresponse.Registration{}, servicehelpers.MapRegistrationDependencyError(err)
	}

	return service.mapper.ToRegistrationResponse(account), nil
}

func (service *RegistrationService) resolvePerson(
	ctx context.Context,
	request identityrequest.Register,
) (foundationresponse.Person, error) {
	if request.Gid != "" {
		person, err := service.foundationClient.GetPersonByGID(
			ctx,
			foundationrequest.PersonByGID{GID: request.Gid},
		)
		if err != nil {
			return foundationresponse.Person{}, servicehelpers.MapRegistrationDependencyError(err)
		}
		if !servicehelpers.PersonHasEmail(person, request.Email) {
			return foundationresponse.Person{}, coreerrors.NewError(identityerror.InvalidRegistration)
		}
		return person, nil
	}

	result, err := service.foundationClient.SearchPersons(
		ctx,
		servicehelpers.RegistrationEmailSearch(request.Email),
	)
	if err != nil {
		return foundationresponse.Person{}, servicehelpers.MapRegistrationDependencyError(err)
	}
	if len(result.Items) > 1 {
		return foundationresponse.Person{}, coreerrors.NewError(identityerror.RegistrationFailed)
	}
	if len(result.Items) == 1 {
		return result.Items[0], nil
	}

	person, _, err := service.foundationClient.CreatePerson(
		ctx,
		servicehelpers.BuildRegistrationPersonRequest(request),
	)
	if err != nil {
		return foundationresponse.Person{}, servicehelpers.MapRegistrationDependencyError(err)
	}
	if !servicehelpers.PersonHasEmail(person, request.Email) {
		return foundationresponse.Person{}, coreerrors.NewError(identityerror.RegistrationAlreadyExists)
	}
	return person, nil
}
