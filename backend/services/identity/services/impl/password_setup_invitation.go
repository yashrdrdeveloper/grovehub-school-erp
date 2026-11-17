package impl

import (
    "context"
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "crypto/sha256"
    "encoding/base64"
    "encoding/hex"
    "errors"
    "net/url"
    "time"

    "gorm.io/gorm"
    coretransaction "grovehub/libs/core/database/transaction"
    "grovehub/services/identity/entity"
    identityrepository "grovehub/services/identity/repository"
    identityservices "grovehub/services/identity/services"
)

var ErrInvalidInvitationConfiguration = errors.New("identity invitation encryption key and HTTPS setup URL are required")

// PasswordSetupInvitationQueue implements durable invitation intent.
// It stores the secret only encrypted using a 256-bit, independently managed key.
// QueuePasswordSetup MUST run inside the registration transaction.
type PasswordSetupInvitationQueue struct {
    db *gorm.DB
    tokens identityrepository.PasswordSetupTokenRepository
    encryptionKey []byte
    setupURL *url.URL
}

func NewPasswordSetupInvitationQueue(db *gorm.DB, tokens identityrepository.PasswordSetupTokenRepository, encryptionKey []byte, rawSetupURL string) (identityservices.InvitationQueue, error) {
    parsed, err := url.Parse(rawSetupURL)
    if err != nil || parsed.Scheme != "https" || parsed.Host == "" || len(encryptionKey) != 32 || db == nil || tokens == nil {
        return nil, ErrInvalidInvitationConfiguration
    }
    copied := append([]byte(nil), encryptionKey...)
    return &PasswordSetupInvitationQueue{db: db, tokens: tokens, encryptionKey: copied, setupURL: parsed}, nil
}

func (queue *PasswordSetupInvitationQueue) QueuePasswordSetup(ctx context.Context, account entity.Account) error {
    tx := coretransaction.Database(ctx, queue.db)
    var existing entity.PasswordSetupMail
    if err := tx.Where("account_id = ?", account.AccountId).Take(&existing).Error; err == nil {
        return nil
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        return err
    }

    secret := make([]byte, 32)
    if _, err := rand.Read(secret); err != nil { return err }
    digest := sha256.Sum256(secret)
    now := time.Now().UTC()
    token := entity.PasswordSetupToken{
        TokenHash: hex.EncodeToString(digest[:]),
        AccountId: account.AccountId,
        ExpiresAt: now.Add(30 * time.Minute),
        CreatedAt: now,
    }
    if err := queue.tokens.Issue(ctx, token); err != nil { return err }

    link := *queue.setupURL
    values := link.Query()
    values.Set("token", base64.RawURLEncoding.EncodeToString(secret))
    link.RawQuery = values.Encode()

    block, err := aes.NewCipher(queue.encryptionKey)
    if err != nil { return err }
    aead, err := cipher.NewGCM(block)
    if err != nil { return err }
    nonce := make([]byte, aead.NonceSize())
    if _, err := rand.Read(nonce); err != nil { return err }
    encrypted := aead.Seal(nil, nonce, []byte(link.String()), []byte(account.AccountId))
    mail := entity.PasswordSetupMail{
        AccountId: account.AccountId, Email: account.Email,
        Ciphertext: encrypted, Nonce: nonce,
        NextAttemptAt: now, CreatedAt: now,
    }
    return tx.Create(&mail).Error
}

// SystemMailSender is a privileged, system-owned mail adapter. It must use a
// configured transactional mailbox and must never accept caller OAuth tokens.
type SystemMailSender interface {
    SendPasswordSetup(ctx context.Context, recipient string, setupLink string) error
}

// PasswordSetupMailWorker claims queued invitations with a short DB lease.
// At-least-once delivery is intentional; duplicate emails contain the same link.
type PasswordSetupMailWorker struct {
    db *gorm.DB
    encryptionKey []byte
    sender SystemMailSender
}

func NewPasswordSetupMailWorker(db *gorm.DB, encryptionKey []byte, sender SystemMailSender) (*PasswordSetupMailWorker, error) {
    if db == nil || len(encryptionKey) != 32 || sender == nil { return nil, ErrInvalidInvitationConfiguration }
    return &PasswordSetupMailWorker{db: db, encryptionKey: append([]byte(nil), encryptionKey...), sender: sender}, nil
}

func (worker *PasswordSetupMailWorker) RunOnce(ctx context.Context) (bool, error) {
    var job entity.PasswordSetupMail
    now := time.Now().UTC()
    err := worker.db.WithContext(ctx).Raw(`
        UPDATE identity.password_setup_mail SET
            lease_until = ?, attempts = attempts + 1
        WHERE account_id = (
            SELECT account_id FROM identity.password_setup_mail
            WHERE sent_at IS NULL AND next_attempt_at <= ?
              AND (lease_until IS NULL OR lease_until < ?)
            ORDER BY next_attempt_at
            FOR UPDATE SKIP LOCKED LIMIT 1
        ) RETURNING *
    `, now.Add(2*time.Minute), now, now).Scan(&job).Error
    if err != nil { return false, err }
    if job.AccountId == "" { return false, nil }

    block, err := aes.NewCipher(worker.encryptionKey)
    if err != nil { return true, err }
    aead, err := cipher.NewGCM(block)
    if err != nil { return true, err }
    plaintext, err := aead.Open(nil, job.Nonce, job.Ciphertext, []byte(job.AccountId))
    if err != nil { return true, err }

    sendErr := worker.sender.SendPasswordSetup(ctx, job.Email, string(plaintext))
    if sendErr != nil {
        // Bounded retries; no plaintext link, mailbox credentials or provider errors in DB.
        backoff := time.Duration(1<<min(job.Attempts, 8)) * time.Second
        updateErr := worker.db.WithContext(ctx).Model(&entity.PasswordSetupMail{}).
            Where("account_id = ? AND sent_at IS NULL", job.AccountId).
            Updates(map[string]any{"next_attempt_at": time.Now().UTC().Add(backoff), "lease_until": nil}).Error
        if updateErr != nil { return true, updateErr }
        return true, sendErr
    }
    err = worker.db.WithContext(ctx).Model(&entity.PasswordSetupMail{}).
        Where("account_id = ? AND sent_at IS NULL", job.AccountId).
        Updates(map[string]any{"sent_at": time.Now().UTC(), "lease_until": nil}).Error
    return true, err
}
