package impl

import (
    "context"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "net/http"
    "strings"
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
    "grovehub/services/identity/entity"
    identityservices "grovehub/services/identity/services"
    "errors"
)

var ErrRegistrationNotAuthorized = errors.New("registration requires an authorized administrator")

type accountContextKey struct{}

// AccountIDFromContext reads only identity verified by SessionAuthenticator.
// Never promote untrusted X-User-ID, X-Role or similar headers into this context.
func AccountIDFromContext(ctx context.Context) (entity.AccountId, bool) {
    id, ok := ctx.Value(accountContextKey{}).(entity.AccountId)
    return id, ok && id != ""
}

// SessionAuthenticator validates opaque access tokens and live account state.
type SessionAuthenticator struct{ db *gorm.DB }

func NewSessionAuthenticator(db *gorm.DB) *SessionAuthenticator { return &SessionAuthenticator{db: db} }

func (authenticator *SessionAuthenticator) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
        authorization := request.Header.Get("Authorization")
        if !strings.HasPrefix(authorization, "Bearer ") {
            http.Error(writer, "unauthorized", http.StatusUnauthorized)
            return
        }
        bearer := strings.TrimPrefix(authorization, "Bearer ")
        secret, err := base64.RawURLEncoding.DecodeString(bearer)
        if err != nil || len(secret) != 32 {
            http.Error(writer, "unauthorized", http.StatusUnauthorized)
            return
        }
        digest := sha256.Sum256(secret)
        now := time.Now().UTC()
        var session entity.Session
        err = authenticator.db.WithContext(request.Context()).
            Where("access_hash = ? AND revoked_at IS NULL AND access_expires_at > ?", hex.EncodeToString(digest[:]), now).
            Take(&session).Error
        if err != nil {
            http.Error(writer, "unauthorized", http.StatusUnauthorized)
            return
        }
        var account entity.Account
        err = authenticator.db.WithContext(request.Context()).
            Where("account_id = ? AND status = ?", session.AccountId, entity.AccountStatusActive).
            Take(&account).Error
        if err != nil {
            http.Error(writer, "unauthorized", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(request.Context(), accountContextKey{}, account.AccountId)
        next.ServeHTTP(writer, request.WithContext(ctx))
    })
}

// AdminRegistrationAuthorizer is a temporary explicit bootstrap policy.
// A configured allowlist is used until Identity has a real RBAC/PDP integration.
// No configured admin IDs means fail closed; only verified access sessions count.
type AdminRegistrationAuthorizer struct{ adminIDs map[entity.AccountId]struct{} }

func NewAdminRegistrationAuthorizer(idsCSV string) identityservices.RegistrationAuthorizer {
    ids := make(map[entity.AccountId]struct{})
    for _, value := range strings.Split(idsCSV, ",") {
        value = strings.TrimSpace(value)
        if value == "" { continue }
        parsed, err := uuid.Parse(value)
        if err == nil { ids[entity.AccountId(parsed.String())] = struct{}{} }
    }
    return &AdminRegistrationAuthorizer{adminIDs: ids}
}

func (authorizer *AdminRegistrationAuthorizer) AuthorizeRegistration(ctx context.Context) error {
    accountID, ok := AccountIDFromContext(ctx)
    if !ok { return ErrRegistrationNotAuthorized }
    if _, allowed := authorizer.adminIDs[accountID]; !allowed { return ErrRegistrationNotAuthorized }
    return nil
}
