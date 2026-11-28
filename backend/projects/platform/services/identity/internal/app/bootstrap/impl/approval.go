package impl

import (
	"grovehub/services/identity/internal/app/services/registration/operations/invitation"
	"grovehub/services/identity/internal/app/config"
	resourceproducer "grovehub/services/identity/internal/app/events/producers/resource/impl"
	mapper "grovehub/services/identity/internal/app/mappers/account/impl"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	assignmentimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	policyimpl "grovehub/services/identity/internal/app/repository/approvalpolicy/impl"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	"grovehub/services/identity/internal/app/services/account/operations/transition"
	"grovehub/services/identity/internal/app/services/accountrole/operations/replace"
	"grovehub/services/identity/internal/app/services/approval"
	approvalimpl "grovehub/services/identity/internal/app/services/approval/impl"
	"grovehub/services/identity/internal/app/services/approval/operations/cutover"
	"grovehub/services/identity/internal/app/services/authorization"
	"grovehub/services/identity/internal/app/services/authorization/operations/protect"
	"grovehub/services/identity/internal/app/services/role/operations/lifecycle"
)

func (b *builder) newApprovalExecutor(repositories repositories, cfg config.Config, reader accountrepo.Reader, authorizer authorization.Service, gate *cutover.Gate) approval.Executor {
	if cfg.LegacyRBACDevelopment || cfg.ApprovalExecutionTenant == "" {
		return nil
	}
	guard := protect.New(repositories.authority, authorityimpl.NewAdministratorReader(b.db))
	readiness := policyimpl.NewReadiness(b.db)
	effects := approvalimpl.Effects{
		PendingInvitation: invitation.NewApproved(reader, repositories.accountRoles, repositories.authorization, invitation.New(repositories.tokens, repositories.mail, cfg)),
		Audit: resourceproducer.New(b.audit),
		AccountLifecycle: transition.NewNormalized(guard, readiness, authorizer, reader, accountimpl.NewStatusWriter(b.db), repositories.accountRoles, repositories.sessions, mapper.New(), b.newEvents()),
		RoleLifecycle: lifecycle.New(guard, roleimpl.NewLifecycleWriter(b.db)),
		Permissions: roleimpl.NewDraftWriter(b.db),
		Assignment: replace.New(guard, assignmentimpl.NewWriter(b.db)),
	}
	executor := approvalimpl.NewExecutor(repositories.approval, repositories.authority, guard, readiness, authorizer, repositories.authorization, reader, repositories.roles, repositories.accountRoles, effects)
	return cutover.Wrap(gate, executor)
}
