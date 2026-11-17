package impl

import (
    "context"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "time"

    "gorm.io/gorm"
    coreid "grovehub/libs/core/id"
    identityresponse "grovehub/services/identity/dto/response"
    "grovehub/services/identity/entity"
    identityservices "grovehub/services/identity/services"
)

// SessionIssuer creates revocable opaque tokens. Do not interpret these as JWTs;
// protected middleware must validate access hashes against the sessions table.
type SessionIssuer struct { db *gorm.DB }

func NewSessionIssuer(db *gorm.DB) identityservices.SessionIssuer { return &SessionIssuer{db: db} }

func (issuer *SessionIssuer) IssueLogin(ctx context.Context, account entity.Account) (identityresponse.Login, error) {
    accessBytes := make([]byte, 32)
    refreshBytes := make([]byte, 48)
    if _, err := rand.Read(accessBytes); err != nil { return identityresponse.Login{}, err }
    if _, err := rand.Read(refreshBytes); err != nil { return identityresponse.Login{}, err }
    sessionID, err := coreid.NewUUIDv7()
    if err != nil { return identityresponse.Login{}, err }
    accessHash := sha256.Sum256(accessBytes)
    refreshHash := sha256.Sum256(refreshBytes)
    now := time.Now().UTC()
    session := entity.Session{
        SessionID: sessionID, AccountId: account.AccountId,
        AccessHash: hex.EncodeToString(accessHash[:]),
        RefreshHash: hex.EncodeToString(refreshHash[:]),
        AccessExpiresAt: now.Add(15 * time.Minute),
        RefreshExpiresAt: now.Add(30 * 24 * time.Hour),
        CreatedAt: now, UpdatedAt: now,
    }
    if err := issuer.db.WithContext(ctx).Create(&session).Error; err != nil {
        return identityresponse.Login{}, err
    }
    return identityresponse.Login{
        AccessToken: base64.RawURLEncoding.EncodeToString(accessBytes),
        RefreshToken: base64.RawURLEncoding.EncodeToString(refreshBytes),
        TokenType: "Bearer",
        ExpiresIn: 15 * 60,
    }, nil
}
