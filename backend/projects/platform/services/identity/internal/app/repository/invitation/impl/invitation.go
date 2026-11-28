package impl

import (
	"context"
	"errors"
	auditclient "grovehub/services/audit/client/producer"
	"time"

	"gorm.io/gorm"

	coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/internal/app/entity"
	invitationrepo "grovehub/services/identity/internal/app/repository/invitation"
)

type Repository struct {
	db *gorm.DB
	audit auditclient.Client
}

var _ invitationrepo.Repository = (*Repository)(nil)

func New(db *gorm.DB, audit auditclient.Client) invitationrepo.Repository {
	return &Repository{db: db,audit:audit}
}

func (repository *Repository) Queue(ctx context.Context, invitation entity.Invitation) error {
	return coredatabase.RunTransaction(ctx,repository.db,func(tx context.Context) error {
		if err:=coredatabase.TransactionDB(tx,repository.db).Create(&invitation).Error;err!=nil { return err }
		return repository.record(tx,invitation,"QUEUED","",nil,invitation.CreatedAt)
	})
}

func (repository *Repository) Claim(ctx context.Context, now time.Time) (entity.Invitation, bool, error) {
	var invitation entity.Invitation
	err := coredatabase.RunTransaction(ctx,repository.db,func(tx context.Context) error {
	result := coredatabase.TransactionDB(tx, repository.db).Raw(`
		UPDATE identity.invitations
		SET lease_until = ?, attempts = attempts + 1
		WHERE id = (
			SELECT id
			FROM identity.invitations
			WHERE sent_at IS NULL AND suppressed_at IS NULL
			  AND next_attempt_at <= ?
			  AND (lease_until IS NULL OR lease_until < ?)
			ORDER BY next_attempt_at, id
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING *
	`, now.Add(45*time.Second), now, now).Scan(&invitation)

	if result.Error!=nil || invitation.ID=="" { return result.Error }
		return repository.record(tx,invitation,"CLAIMED","",nil,now)
	})
	if err!=nil { return entity.Invitation{},false,err }
	return invitation,invitation.ID!="",nil
}

func (repository *Repository) Delivered(ctx context.Context, invitationID string) error {
	now:=time.Now().UTC()
	return repository.transition(ctx,invitationID,0,"DELIVERED","",nil,now,
		"UPDATE identity.invitations SET sent_at = ?, lease_until = NULL WHERE id = ? AND sent_at IS NULL AND suppressed_at IS NULL AND lease_until IS NOT NULL RETURNING *",now,invitationID)
}

func (repository *Repository) Retry(ctx context.Context, invitationID string, nextAttempt time.Time) error {
	return repository.transition(ctx,invitationID,0,"RETRY_SCHEDULED","DELIVERY_FAILED",&nextAttempt,time.Now().UTC(),
		"UPDATE identity.invitations SET next_attempt_at = ?, lease_until = NULL WHERE id = ? AND sent_at IS NULL AND suppressed_at IS NULL AND lease_until IS NOT NULL RETURNING *",nextAttempt,invitationID)
}

// LockClaim fences a worker whose lease expired while waiting for tenant authority.
func (repository *Repository) LockClaim(ctx context.Context, invitation entity.Invitation, now time.Time) (bool, error) {
	if !coredatabase.TransactionActive(ctx) {
		return false, coredatabase.ErrMissingDatabase
	}
	var row struct { ID string }
	result := coredatabase.TransactionDB(ctx, repository.db).Raw(`
		SELECT id FROM identity.invitations
		WHERE id = ? AND attempts = ? AND lease_until > ?
		  AND sent_at IS NULL AND suppressed_at IS NULL
		FOR UPDATE
	`, invitation.ID, invitation.Attempts, now).Scan(&row)
	return row.ID != "", result.Error
}

func (repository *Repository) Suppress(ctx context.Context, invitation entity.Invitation, reason string) error {
	now:=time.Now().UTC()
	return repository.transition(ctx,invitation.ID,invitation.Attempts,"SUPPRESSED",reason,nil,now,
		"UPDATE identity.invitations SET suppressed_at = ?, suppression_reason = ?, lease_until = NULL WHERE id = ? AND attempts = ? AND sent_at IS NULL AND suppressed_at IS NULL RETURNING *",now,reason,invitation.ID,invitation.Attempts)
}

func (repository *Repository) transition(ctx context.Context,id string,attempt int,operation string,reason string,next *time.Time,at time.Time,statement string,args ...any) error {
	return coredatabase.RunTransaction(ctx,repository.db,func(tx context.Context) error {
		var entry entity.Invitation
		result:=coredatabase.TransactionDB(tx,repository.db).Raw(statement,args...).Scan(&entry)
		if result.Error!=nil { return result.Error }
		if entry.ID=="" { return errors.New("invitation transition lost ownership") }
		if entry.ID!=id || (attempt>0 && entry.Attempts!=attempt) { return errors.New("invitation transition identity mismatch") }
		return repository.record(tx,entry,operation,reason,next,at)
	})
}
