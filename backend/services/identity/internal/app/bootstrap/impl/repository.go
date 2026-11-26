package impl

import (
	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	accountrepo "grovehub/services/identity/internal/app/repository/account"
	accountimpl "grovehub/services/identity/internal/app/repository/account/impl"
	credentialrepo "grovehub/services/identity/internal/app/repository/credential"
	credentialimpl "grovehub/services/identity/internal/app/repository/credential/impl"
	invitationrepo "grovehub/services/identity/internal/app/repository/invitation"
	invitationimpl "grovehub/services/identity/internal/app/repository/invitation/impl"
	passwordtokenrepo "grovehub/services/identity/internal/app/repository/passwordtoken"
	passwordtokenimpl "grovehub/services/identity/internal/app/repository/passwordtoken/impl"
	sessionrepo "grovehub/services/identity/internal/app/repository/session"
	sessionimpl "grovehub/services/identity/internal/app/repository/session/impl"
)

type repositories struct {
	accounts accountrepo.Repository
	credentials credentialrepo.Repository
	sessions sessionrepo.Repository
	tokens passwordtokenrepo.Repository
	mail invitationrepo.Repository
	unit transaction.UnitOfWork
}

func (b *builder) newRepository() repositories {
	return repositories{
		accounts: accountimpl.New(b.db),
		credentials: credentialimpl.New(b.db),
		sessions: sessionimpl.New(b.db),
		tokens: passwordtokenimpl.New(b.db),
		mail: invitationimpl.New(b.db),
		unit: transaction.NewUnitOfWork(b.db),
	}
}
