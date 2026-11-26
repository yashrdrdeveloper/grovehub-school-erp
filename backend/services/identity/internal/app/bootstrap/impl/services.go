package impl

import (
	metadataaudit "grovehub/services/identity/internal/app/repository/rolemetadataaudit/impl"
	metadataproducer "grovehub/services/identity/internal/app/events/producers/rolemetadata/impl"
	capabilityrepo "grovehub/services/identity/internal/app/repository/capability/impl"
	"grovehub/services/identity/internal/app/services/capability"
	capabilityimpl "grovehub/services/identity/internal/app/services/capability/impl"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	"grovehub/services/identity/internal/app/services/accountrole"
	accountroleimpl "grovehub/services/identity/internal/app/services/accountrole/impl"
	"grovehub/services/identity/internal/app/services/approval"
	approvalimpl "grovehub/services/identity/internal/app/services/approval/impl"
	"grovehub/services/identity/internal/app/config"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountread "grovehub/services/identity/internal/app/repository/account/impl"
	roleread "grovehub/services/identity/internal/app/repository/role/impl"
	catalog "grovehub/services/identity/internal/app/services/account/operations/roles"
	"grovehub/services/identity/internal/app/services/security/operations/routepolicy"
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
	role "grovehub/services/identity/internal/app/services/role"
	roleimpl "grovehub/services/identity/internal/app/services/role/impl"
	validator "grovehub/services/identity/internal/app/validators/account/impl"
)

type serviceGraph struct {
	capabilities capability.Service
	accountRoles accountrole.Service
	approvalReads approval.Reader
	approvalMutations approval.Service
	approvalExecutor approval.Executor
	account       account.Service
	authorization authorization.Service
	policy *routepolicy.Policy
	catalog *catalog.Operation
	registry      registry.Service
	roles         role.Service
	auth          auth.Service
	worker        *invitation.Worker
}

func (b *builder) newServices(repositories repositories, cfg config.Config) serviceGraph {
	credentials := credentialimpl.New(repositories.credentials)
	invitations := invitation.New(repositories.tokens, repositories.mail, cfg)
	accountEvents := b.newEvents()
	registrations := registrationimpl.New(repositories.accounts, b.foundation, validator.New(), invitations, repositories.unit, accountEvents)

	var reader accountrepo.Reader = repositories.accounts
	accountValidator := validator.New()
	if !cfg.LegacyRBACDevelopment {
		reader = accountread.NewNormalizedReader(b.db)
		accountValidator = validator.NewNormalized()
	}

	accounts := accountimpl.New(
		repositories.accounts,
		repositories.tokens,
		repositories.sessions,
		credentials,
		registrations,
		mapper.New(),
		accountValidator,
		repositories.unit,
		accountEvents,
		reader,
	)
	authentication := authimpl.New(
		repositories.accounts,
		repositories.sessions,
		repositories.tokens,
		credentials,
		invitations,
		mapper.New(),
		repositories.unit,
		reader,
	)

	var worker *invitation.Worker
	if b.integration != nil && len(cfg.EncryptionKey) == 32 {
		worker = invitation.NewWorker(repositories.mail, repositories.accounts, b.integration, cfg)
	}

	authorizer := authorizationimpl.New(repositories.authorization)
	gate := cutover.New(cfg.ApprovalExecutionTenant)
	policy := routepolicy.NewWithExecution(authorizer, cfg.LegacyRBACDevelopment, gate)
	return serviceGraph{
		capabilities: capabilityimpl.New(capabilityrepo.New(b.db), repositories.authority, authorizer, cfg.ApprovalExecutionTenant),
		approvalExecutor: b.newApprovalExecutor(repositories, cfg, reader, authorizer, gate),
		accountRoles: accountroleimpl.New(repositories.accountRoles, authorizer),
		approvalReads: approvalimpl.NewReader(repositories.approvalReads, repositories.approval, repositories.authority, repositories.authorization),
		approvalMutations: approvalimpl.New(repositories.approval, repositories.authority, authorizer, repositories.authorization, reader, repositories.roles, repositories.accountRoles),
		account: accountimpl.NewGuarded(accounts, policy), auth: authentication, worker: worker,
		policy: policy, catalog: catalog.New(roleread.NewCatalogReader(b.db), policy, cfg.LegacyRBACDevelopment),
		authorization: authorizer, registry: registryimpl.New(repositories.registry),
		roles: roleimpl.NewAudited(repositories.roles, authorizer, repositories.authority, roleread.NewDraftInspector(b.db), metadataproducer.New(metadataaudit.New(b.db), b.events)),
	}
}
