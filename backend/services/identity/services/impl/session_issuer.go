package impl

import (
    "context"
    "errors"
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


// Refresh atomically rotates both bearer credentials for the same session.
// The old refresh token can be successfully redeemed at most once.
func (issuer *SessionIssuer) Refresh(ctx context.Context, rawToken string) (identityresponse.Login, error) {
    oldBytes, err := base64.RawURLEncoding.DecodeString(rawToken)
    if err != nil || len(oldBytes) != 48 { return identityresponse.Login{}, errors.New("invalid refresh credentials") }
    accessBytes := make([]byte, 32)
    refreshBytes := make([]byte, 48)
    if _, err := rand.Read(accessBytes); err != nil { return identityresponse.Login{}, err }
    if _, err := rand.Read(refreshBytes); err != nil { return identityresponse.Login{}, err }
    oldDigest := sha256.Sum256(oldBytes)
    accessDigest := sha256.Sum256(accessBytes)
    refreshDigest := sha256.Sum256(refreshBytes)
    now := time.Now().UTC()
    err = issuer.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        var result struct { AccountID string `gorm:"column:account_id"` }
        updateErr := tx.Raw(`
            UPDATE identity.sessions SET access_hash = ?, refresh_hash = ?,
                access_expires_at = ?, updated_at = ?
            WHERE refresh_hash = ? AND revoked_at IS NULL AND refresh_expires_at > ?
            RETURNING account_id
        `, hex.EncodeToString(accessDigest[:]), hex.EncodeToString(refreshDigest[:]),
            now.Add(15*time.Minute), now, hex.EncodeToString(oldDigest[:]), now).Scan(&result).Error
        if updateErr != nil { return updateErr }
        if result.AccountID == "" { return errors.New("invalid refresh credentials") }
        var account entity.Account
        if err := tx.Where("account_id = ? AND status = ?", result.AccountID, entity.AccountStatusActive).Take(&account).Error; err != nil {
            return errors.New("invalid refresh credentials")
        }
        return nil
    })
    if err != nil { return identityresponse.Login{}, err }
    return identityresponse.Login{
        AccessToken: base64.RawURLEncoding.EncodeToString(accessBytes),
        RefreshToken: base64.RawURLEncoding.EncodeToString(refreshBytes),
        TokenType: "Bearer",
        ExpiresIn: 900,
    }, nil
}

func (issuer *SessionIssuer) RevokeAccess(ctx context.Context, accountID entity.AccountId, rawToken string) error {
    secret, err := base64.RawURLEncoding.DecodeString(rawToken)
    if err != nil || len(secret) != 32 { return errors.New("invalid access credential") }
    digest := sha256.Sum256(secret)
    now := time.Now().UTC()
    result := issuer.db.WithContext(ctx).Model(&entity.Session{}).
        Where("account_id = ? AND access_hash = ? AND revoked_at IS NULL", accountID, hex.EncodeToString(digest[:])).
        Updates(map[string]any{"revoked_at": now, "updated_at": now})
    if result.Error != nil { return result.Error }
    if result.RowsAffected == 0 { return errors.New("session is not active") }
    return nil
}
