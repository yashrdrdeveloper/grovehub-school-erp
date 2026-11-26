package impl

import (
	"context"
	"time"

	"gorm.io/gorm"

	transaction "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
	"grovehub/services/identity/internal/app/entity"
	invitationrepo "grovehub/services/identity/internal/app/repository/invitation"
)

type Repository struct {
	db *gorm.DB
}

var _ invitationrepo.Repository = (*Repository)(nil)

func New(db *gorm.DB) invitationrepo.Repository {
	return &Repository{db: db}
}

func (repository *Repository) Queue(ctx context.Context, invitation entity.Invitation) error {
	return transaction.Database(ctx, repository.db).Create(&invitation).Error
}

func (repository *Repository) Claim(ctx context.Context, now time.Time) (entity.Invitation, bool, error) {
	var invitation entity.Invitation
	result := transaction.Database(ctx, repository.db).Raw(`
		UPDATE identity.invitations
		SET lease_until = ?, attempts = attempts + 1
		WHERE id = (
			SELECT id
			FROM identity.invitations
			WHERE sent_at IS NULL AND suppressed_at IS NULL
			  AND next_attempt_at <= ?
			  AND (lease_until IS NULL OR lease_until < ?)
			ORDER BY next_attempt_at
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING *
	`, now.Add(45*time.Second), now, now).Scan(&invitation)

	return invitation, invitation.ID != "", result.Error
}

func (repository *Repository) Delivered(ctx context.Context, invitationID string) error {
	return transaction.Database(ctx, repository.db).
		Model(&entity.Invitation{}).
		Where("id = ?", invitationID).
		Updates(map[string]any{
			"sent_at":     time.Now().UTC(),
			"lease_until": nil,
		}).Error
}

func (repository *Repository) Retry(ctx context.Context, invitationID string, nextAttempt time.Time) error {
	return transaction.Database(ctx, repository.db).
		Model(&entity.Invitation{}).
		Where("id = ?", invitationID).
		Updates(map[string]any{
			"next_attempt_at": nextAttempt,
			"lease_until":     nil,
		}).Error
}

// LockClaim fences a worker whose lease expired while waiting for tenant authority.
func (repository *Repository) LockClaim(ctx context.Context, invitation entity.Invitation, now time.Time) (bool, error) {
	if !transaction.Active(ctx) {
		return false, transaction.ErrMissingDatabase
	}
	var row struct { ID string }
	result := transaction.Database(ctx, repository.db).Raw(`
		SELECT id FROM identity.invitations
		WHERE id = ? AND attempts = ? AND lease_until > ?
		  AND sent_at IS NULL AND suppressed_at IS NULL
		FOR UPDATE
	`, invitation.ID, invitation.Attempts, now).Scan(&row)
	return row.ID != "", result.Error
}

func (repository *Repository) Suppress(ctx context.Context, invitation entity.Invitation, reason string) error {
	return transaction.Database(ctx, repository.db).Model(&entity.Invitation{}).
		Where("id = ? AND attempts = ? AND sent_at IS NULL AND suppressed_at IS NULL", invitation.ID, invitation.Attempts).
		Updates(map[string]any{"suppressed_at": time.Now().UTC(), "suppression_reason": reason, "lease_until": nil}).Error
}
