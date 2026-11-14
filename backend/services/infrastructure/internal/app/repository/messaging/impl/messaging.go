// Package impl owns all Infrastructure messaging SQL; libs/core stays SQL-free.
package impl

import (
    "context"
    "encoding/json"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    coretransaction "grovehub/libs/core/database/transaction"
    errorcode "grovehub/services/infrastructure/internal/app/constants/error"
    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
    "grovehub/services/infrastructure/internal/app/entity"
    repository "grovehub/services/infrastructure/internal/app/repository/messaging"
    validator "grovehub/services/infrastructure/internal/app/validators/messaging"
)

type messagingRepository struct {
    db *gorm.DB
    validator validator.Validator
}

var _ repository.Repository = (*messagingRepository)(nil)

func New(db *gorm.DB, eventValidator validator.Validator) repository.Repository {
    return &messagingRepository{db:db, validator:eventValidator}
}

// Enqueue participates in the caller's GORM transaction (same PostgreSQL DB).
func (store *messagingRepository) Enqueue(ctx context.Context, event dto.Event) error {
    if store.db == nil || store.validator == nil { return errorcode.ErrInvalidStore }
    if err := store.validator.ValidateEvent(event); err != nil { return err }
    metadata,err := json.Marshal(event.Metadata)
    if err != nil { return err }
    entry := entity.Messaging{
        EventID:event.EventID, Topic:event.Topic, Payload:event.Payload,
        Metadata:json.RawMessage(metadata), Status:entity.MessagingPending,
        NextAttemptAt:time.Now().UTC(),
    }
    return coretransaction.Database(ctx, store.db).Table(entry.TableName()).Create(&entry).Error
}

func (store *messagingRepository) Claim(ctx context.Context, limit int, lease time.Duration) ([]entity.Messaging, error) {
    if store.db == nil { return nil, errorcode.ErrInvalidStore }
    if limit <= 0 { return nil, nil }
    if lease <= 0 { lease = 30*time.Second }
    now := time.Now().UTC()
    until := now.Add(lease)
    var entries []entity.Messaging
    err := store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        table := entity.Messaging{}.TableName()
        if err := tx.Table(table).Clauses(clause.Locking{Strength:"UPDATE",Options:"SKIP LOCKED"}).
            Where("(status = ? AND next_attempt_at <= ?) OR (status = ? AND lease_until <= ?)",entity.MessagingPending,now,entity.MessagingSending,now).
            Order("created_at, event_id").Limit(limit).Find(&entries).Error; err != nil { return err }
        for index := range entries {
            if err := tx.Table(table).Where("event_id = ?",entries[index].EventID).
                Updates(map[string]any{"status":entity.MessagingSending,"lease_until":until,"attempts":gorm.Expr("attempts + 1")}).Error;err != nil {return err}
            entries[index].Attempts++
        }
        return nil
    })
    return entries,err
}

func (store *messagingRepository) MarkSent(ctx context.Context,eventID string) error {
    if store.db == nil { return errorcode.ErrInvalidStore }
    return store.db.WithContext(ctx).Table(entity.Messaging{}.TableName()).
        Where("event_id = ? AND status = ?",eventID,entity.MessagingSending).
        Updates(map[string]any{"status":entity.MessagingSent,"lease_until":nil,"published_at":time.Now().UTC()}).Error
}

func (store *messagingRepository) MarkRetry(ctx context.Context,eventID string,attempt int) error {
    if store.db == nil { return errorcode.ErrInvalidStore }
    seconds:=1
    for index:=1;index<attempt && seconds<64;index++{seconds*=2}
    next:=time.Now().UTC().Add(time.Duration(seconds)*time.Second)
    return store.db.WithContext(ctx).Table(entity.Messaging{}.TableName()).
        Where("event_id = ? AND status = ?",eventID,entity.MessagingSending).
        Updates(map[string]any{"status":entity.MessagingPending,"lease_until":nil,"next_attempt_at":next}).Error
}

// PruneSent deletes only successfully published records older than the
// configured retention time, in bounded batches to avoid long table locks.
func (store *messagingRepository) PruneSent(ctx context.Context,olderThan time.Time,limit int) (int64,error) {
    if store.db == nil { return 0,errorcode.ErrInvalidStore }
    if limit<=0 { return 0,nil }
    table:=entity.Messaging{}.TableName()
    ids:=store.db.WithContext(ctx).Table(table).Select("event_id").
        Where("status = ? AND published_at < ?",entity.MessagingSent,olderThan).
        Order("published_at, event_id").Limit(limit)
    result:=store.db.WithContext(ctx).Table(table).Where("event_id IN (?)",ids).
        Delete(&entity.Messaging{})
    return result.RowsAffected,result.Error
}
