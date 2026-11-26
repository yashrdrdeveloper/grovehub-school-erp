package impl

import (
	accountrolerepo "grovehub/services/identity/internal/app/repository/accountrole"
	accountroleimpl "grovehub/services/identity/internal/app/repository/accountrole/impl"
	approvalrepo "grovehub/services/identity/internal/app/repository/approval"
	approvalimpl "grovehub/services/identity/internal/app/repository/approval/impl"
	producerimpl "grovehub/services/identity/internal/app/events/producers/approval/impl"
	"grovehub/services/identity/internal/app/repository/authority"
	authorityimpl "grovehub/services/identity/internal/app/repository/authority/impl"
	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	authorizationrepo "grovehub/services/identity/internal/app/repository/authorization"
	authorizationimpl "grovehub/services/identity/internal/app/repository/authorization/impl"
	credentialrepo "grovehub/services/identity/internal/app/repository/credential"
	credentialimpl "grovehub/services/identity/internal/app/repository/credential/impl"
	invitationrepo "grovehub/services/identity/internal/app/repository/invitation"
	invitationimpl "grovehub/services/identity/internal/app/repository/invitation/impl"
	passwordtokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	passwordtokenimpl "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
	registryrepo "grovehub/services/identity/internal/app/repository/registry"
	registryimpl "grovehub/services/identity/internal/app/repository/registry/impl"
	rolerepo "grovehub/services/identity/internal/app/repository/role"
	roleimpl "grovehub/services/identity/internal/app/repository/role/impl"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	sessionimpl "grovehub/services/identity/internal/app/repository/session/impl"
)

type repositories struct {
	accountRoles accountrolerepo.Repository
	approval approvalrepo.Repository
	approvalReads approvalrepo.Reader
	authority authority.Locker
	accounts      accountrepo.Repository
	authorization authorizationrepo.Repository
	registry      registryrepo.Repository
	roles         rolerepo.Repository
	credentials   credentialrepo.Repository
	sessions      sessionrepo.Repository
	tokens        passwordtokenrepo.Repository
	mail          invitationrepo.Repository
	unit          transaction.UnitOfWork
}

func (b *builder) newRepository() repositories {
	return repositories{
		accountRoles: accountroleimpl.New(b.db),
		approval: approvalimpl.New(b.db, producerimpl.New(b.audit)),
		approvalReads: approvalimpl.NewReader(b.db),
		authority: authorityimpl.New(b.db),
		accounts:      accountimpl.New(b.db),
		authorization: authorizationimpl.New(b.db),
		registry:      registryimpl.New(b.db),
		roles:         roleimpl.New(b.db),
		credentials:   credentialimpl.New(b.db),
		sessions:      sessionimpl.New(b.db),
		tokens:        passwordtokenimpl.New(b.db),
		mail:          invitationimpl.New(b.db),
		unit:          transaction.NewUnitOfWork(b.db),
	}
}
