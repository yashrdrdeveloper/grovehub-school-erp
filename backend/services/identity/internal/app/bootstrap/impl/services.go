package impl

import (
	"grovehub/services/identity/internal/app/config"
	mapper "grovehub/services/identity/internal/app/mappers/account/impl"
	account "grovehub/services/identity/internal/app/services/account"
	accountimpl "grovehub/services/identity/internal/app/services/account/impl"
	auth "grovehub/services/identity/internal/app/services/authentication"
	authimpl "grovehub/services/identity/internal/app/services/authentication/impl"
	authorization "grovehub/services/identity/internal/app/services/authorization"
	authorizationimpl "grovehub/services/identity/internal/app/services/authorization/impl"
	credentialimpl "grovehub/services/identity/internal/app/services/credential/impl"
	registrationimpl "grovehub/services/identity/internal/app/services/registration/impl"
	invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
	registry "grovehub/services/identity/internal/app/services/registry"
	registryimpl "grovehub/services/identity/internal/app/services/registry/impl"
	validator "grovehub/services/identity/internal/app/validators/account/impl"
)

type serviceGraph struct {
	account       account.Service
	authorization authorization.Service
	registry      registry.Service
	auth          auth.Service
	worker        *invitation.Worker
}

func (b *builder) newServices(repositories repositories, cfg config.Config) serviceGraph {
	credentials := credentialimpl.New(repositories.credentials)
	invitations := invitation.New(repositories.tokens, repositories.mail, cfg)
	accountEvents := b.newEvents()
	registrations := registrationimpl.New(repositories.accounts, b.foundation, validator.New(), invitations, repositories.unit, accountEvents)

	accounts := accountimpl.New(
		repositories.accounts,
		repositories.tokens,
		repositories.sessions,
		credentials,
		registrations,
		mapper.New(),
		validator.New(),
		repositories.unit,
		accountEvents,
	)
	authentication := authimpl.New(
		repositories.accounts,
		repositories.sessions,
		repositories.tokens,
		credentials,
		invitations,
		mapper.New(),
		repositories.unit,
	)

	var worker *invitation.Worker
	if b.integration != nil && len(cfg.EncryptionKey) == 32 {
		worker = invitation.NewWorker(repositories.mail, repositories.accounts, b.integration, cfg)
	}

	return serviceGraph{account: accounts, auth: authentication, worker: worker, authorization: authorizationimpl.New(repositories.authorization), registry: registryimpl.New(repositories.registry)}
}
