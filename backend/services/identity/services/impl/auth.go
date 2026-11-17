package impl

import (
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "time"
	"context"
	"errors"

 "gorm.io/gorm"

 coretransaction "grovehub/libs/core/database/transaction"

	coreerrors "grovehub/libs/core/errors"
	annotationvalidation "grovehub/libs/core/annotation/validation"
	coreidempotency "grovehub/libs/core/idempotency"
    headercontext "grovehub/libs/core/headercontext"
    coreevents "grovehub/libs/core/events"
	foundationclient "grovehub/services/foundation/client/person"
		identityerror "grovehub/services/identity/constants/error"
	identityrequest "grovehub/services/identity/dto/request"
	identityresponse "grovehub/services/identity/dto/response"
	identitymappers "grovehub/services/identity/mappers"
	identityservices "grovehub/services/identity/services"
	servicehelpers "grovehub/services/identity/services/helpers"
    "grovehub/services/identity/entity"
    validatorhelpers "grovehub/services/identity/validators/helpers"
	identityvalidators "grovehub/services/identity/validators"
)

type AuthService struct {
	foundationClient          foundationclient.Client
	accountService            identityservices.AccountService
	passwordCredentialService identityservices.PasswordCredentialService
	mapper                    identitymappers.AuthMapper
	validator                 identityvalidators.AuthValidator
 db *gorm.DB
    dependencies identityservices.AuthDependencies
}

func NewAuthService(
	foundationClient foundationclient.Client,
	accountService identityservices.AccountService,
	passwordCredentialService identityservices.PasswordCredentialService,
	mapper identitymappers.AuthMapper,
	validator identityvalidators.AuthValidator,
 db *gorm.DB,
    dependencies ...identityservices.AuthDependencies,
) identityservices.AuthService {
	var configured identityservices.AuthDependencies
    if len(dependencies) > 0 { configured = dependencies[0] }
	return &AuthService{
        dependencies: configured,
		foundationClient:          foundationClient,
		accountService:            accountService,
		passwordCredentialService: passwordCredentialService,
		mapper:                    mapper,
		validator:                 validator,
 db: db,
	}
}

func (service *AuthService) Register(ctx context.Context, request identityrequest.RegisterRequest) (identityresponse.Register, error) {
    if service.dependencies.RegistrationAuthorizer == nil { return identityresponse.Register{}, coreerrors.NewError(identityerror.RegistrationUnauthorized) }
    if err := service.dependencies.RegistrationAuthorizer.AuthorizeRegistration(ctx); err != nil { return identityresponse.Register{}, coreerrors.NewError(identityerror.RegistrationUnauthorized) }
    if service.dependencies.InvitationQueue == nil { return identityresponse.Register{}, coreerrors.NewError(identityerror.RegistrationUnavailable) }
	idempotencyKey, err := coreidempotency.ValidateKey(headercontext.IdempotencyKey(ctx))
	if err != nil {
		return identityresponse.Register{}, coreerrors.NewError(identityerror.InvalidRegistrationIdempotencyKey)
	}
    ctx = coreevents.WithCausationID(ctx, idempotencyKey)

	validated, err := annotationvalidation.Process(request, service.validator.RegisterRules())
	if err != nil {
		return identityresponse.Register{}, err
	}
	requestHash, err := coreidempotency.Fingerprint(validated.RequestBody)
	if err != nil {
		return identityresponse.Register{}, coreerrors.Wrap(identityerror.RegistrationFailed, err)
	}

    var registration identityresponse.Register
    err = coretransaction.Run(ctx, service.db, func(txCtx context.Context) error {
        var operationErr error
        registration, operationErr = service.register(txCtx, validated.RequestBody, idempotencyKey, requestHash)
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
        if err := service.dependencies.InvitationQueue.QueuePasswordSetup(ctx, account); err != nil { return identityresponse.Register{}, coreerrors.Wrap(identityerror.RegistrationFailed, err) }
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

	person, err = service.foundationClient.UpdatePerson(ctx, foundationclient.UpdatePersonRequest{
		Params: foundationclient.PersonByGID{GID: person.GID},
		RequestBody: servicehelpers.BuildRegisterPersonUpdate(validated),
	})
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

    if err := service.dependencies.InvitationQueue.QueuePasswordSetup(ctx, account); err != nil { return identityresponse.Register{}, coreerrors.Wrap(identityerror.RegistrationFailed, err) }

	return service.mapper.ToRegisterResponse(account), nil
}

func (service *AuthService) resolveRegisterPerson(
	ctx context.Context,
	request identityrequest.Register,
) (foundationclient.Person, error) {
	result, err := service.foundationClient.SearchPersons(ctx, foundationclient.SearchPersonsRequest{
		RequestBody: servicehelpers.RegisterEmailSearch(request.Email),
	})
	if err != nil {
		return foundationclient.Person{}, servicehelpers.MapAuthDependencyError(err)
	}
	if len(result.Items) > 1 {
		return foundationclient.Person{}, coreerrors.NewError(identityerror.RegistrationFailed)
	}
	if len(result.Items) == 1 {
		return result.Items[0], nil
	}

	person, _, err := service.foundationClient.CreatePerson(headercontext.WithoutIdempotencyKey(ctx), foundationclient.CreatePersonRequest{
		RequestBody: servicehelpers.BuildRegisterPersonRequest(request),
	})
	if err != nil {
		return foundationclient.Person{}, servicehelpers.MapAuthDependencyError(err)
	}
	if !servicehelpers.PersonHasEmail(person, request.Email) {
		return foundationclient.Person{}, coreerrors.NewError(identityerror.RegistrationAlreadyExists)
	}
	return person, nil
}

func isAccountNotFound(err error) bool {
	var applicationError *coreerrors.Error
	return errors.As(err, &applicationError) &&
		applicationError.Definition().Code() == identityerror.AccountNotFound.Code()
}


// Login performs credential verification and delegates token/session issuance to an
// explicitly configured, revocable session issuer. Missing dependencies fail closed.
func (service *AuthService) Login(ctx context.Context, request identityrequest.LoginRequest) (identityresponse.Login, error) {
    if service.dependencies.SessionIssuer == nil { return identityresponse.Login{}, coreerrors.NewError(identityerror.AuthenticationUnavailable) }
    email, err := validatorhelpers.NormalizeEmail(request.RequestBody.Email)
    if err != nil || request.RequestBody.Password == "" { return identityresponse.Login{}, coreerrors.NewError(identityerror.InvalidCredentials) }
    account, err := service.accountService.GetByEmail(ctx, email)
    if err != nil { return identityresponse.Login{}, coreerrors.NewError(identityerror.InvalidCredentials) }
    if account.Status != entity.AccountStatusActive { return identityresponse.Login{}, coreerrors.NewError(identityerror.InvalidCredentials) }
    valid, err := service.passwordCredentialService.Verify(ctx, string(account.AccountId), request.RequestBody.Password)
    if err != nil || !valid { return identityresponse.Login{}, coreerrors.NewError(identityerror.InvalidCredentials) }
    result, err := service.dependencies.SessionIssuer.IssueLogin(ctx, account)
    if err != nil { return identityresponse.Login{}, coreerrors.Wrap(identityerror.AuthenticationUnavailable, err) }
    return result, nil
}


// SetPassword redeems an emailed, single-use setup token. Consuming the token,
// inserting the credential and activating the account are one Identity transaction.
func (service *AuthService) SetPassword(ctx context.Context, request identityrequest.SetPasswordRequest) (identityresponse.SetPassword, error) {
    if service.dependencies.PasswordSetupTokens == nil { return identityresponse.SetPassword{}, coreerrors.NewError(identityerror.AuthenticationUnavailable) }
    if err := validatorhelpers.ValidatePassword(request.RequestBody.Password); err != nil {
        return identityresponse.SetPassword{}, coreerrors.NewError(identityerror.InvalidPasswordSetup)
    }
    tokenBytes, err := base64.RawURLEncoding.DecodeString(request.RequestBody.Token)
    if err != nil || len(tokenBytes) != 32 {
        return identityresponse.SetPassword{}, coreerrors.NewError(identityerror.InvalidPasswordSetupToken)
    }
    digest := sha256.Sum256(tokenBytes)
    tokenHash := hex.EncodeToString(digest[:])
    err = coretransaction.Run(ctx, service.db, func(txCtx context.Context) error {
        accountID, consumeErr := service.dependencies.PasswordSetupTokens.Consume(txCtx, tokenHash, time.Now().UTC())
        if consumeErr != nil { return coreerrors.NewError(identityerror.InvalidPasswordSetupToken) }
        if credentialErr := service.passwordCredentialService.Set(txCtx, string(accountID), request.RequestBody.Password); credentialErr != nil {
            return coreerrors.Wrap(identityerror.PasswordSetupFailed, credentialErr)
        }
        _, activateErr := service.accountService.Activate(txCtx, string(accountID))
        if activateErr != nil { return coreerrors.Wrap(identityerror.PasswordSetupFailed, activateErr) }
        return nil
    })
    if err != nil { return identityresponse.SetPassword{}, err }
    return identityresponse.SetPassword{Status: string(entity.AccountStatusActive)}, nil
}


// Refresh rotates a valid refresh credential. Failure responses do not disclose
// whether a particular token ever existed.
func (service *AuthService) Refresh(ctx context.Context, request identityrequest.RefreshRequest) (identityresponse.Login, error) {
    if service.dependencies.SessionIssuer == nil { return identityresponse.Login{}, coreerrors.NewError(identityerror.AuthenticationUnavailable) }
    result, err := service.dependencies.SessionIssuer.Refresh(ctx, request.RequestBody.RefreshToken)
    if err != nil { return identityresponse.Login{}, coreerrors.NewError(identityerror.InvalidCredentials) }
    return result, nil
}

func (service *AuthService) Logout(ctx context.Context, accessToken string) error {
    if service.dependencies.SessionIssuer == nil { return coreerrors.NewError(identityerror.AuthenticationUnavailable) }
    accountID, ok := AccountIDFromContext(ctx)
    if !ok { return coreerrors.NewError(identityerror.InvalidCredentials) }
    if err := service.dependencies.SessionIssuer.RevokeAccess(ctx, accountID, accessToken); err != nil {
        return coreerrors.NewError(identityerror.InvalidCredentials)
    }
    return nil
}
