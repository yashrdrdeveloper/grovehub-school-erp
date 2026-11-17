package services

import (
    "context"
    identityresponse "grovehub/services/identity/dto/response"
    "grovehub/services/identity/entity"
    identityrepository "grovehub/services/identity/repository"
)

// RegistrationAuthorizer must authenticate the actor and enforce the register permission.
// Never derive administrative authority from request JSON or unverified headers.
type RegistrationAuthorizer interface {
    AuthorizeRegistration(ctx context.Context) error
}

// InvitationQueue must persist the one-time setup intent in the caller's Identity
// transaction. It must deduplicate by account ID and purpose and never email a password.
// Delivery is performed by a worker after commit, not inside the transaction.
type InvitationQueue interface {
    QueuePasswordSetup(ctx context.Context, account entity.Account) error
}

// SessionIssuer must persist revocable sessions, securely issue access/refresh tokens,
// rotate refresh credentials and prevent issuance for inactive accounts.
type SessionIssuer interface {
    IssueLogin(ctx context.Context, account entity.Account) (identityresponse.Login, error)
    Refresh(ctx context.Context, token string) (identityresponse.Login, error)
}

// AuthDependencies are mandatory for the corresponding live operations.
// Nil dependencies fail closed until bootstrap integrates real implementations.
type AuthDependencies struct {
    RegistrationAuthorizer RegistrationAuthorizer
    InvitationQueue InvitationQueue
    SessionIssuer SessionIssuer
    PasswordSetupTokens identityrepository.PasswordSetupTokenRepository
}
