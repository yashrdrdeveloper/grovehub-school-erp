// Package outbox persists event deliveries within an existing GORM transaction.
package outbox

import (
    "context"
    "encoding/json"
    "errors"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    coretransaction "grovehub/libs/core/database/transaction"
)

const (
    Pending = "PENDING"
    Sending = "SENDING"
    Sent = "SENT"
)

var ErrInvalidMessage = errors.New("outbox message requires event ID, topic and payload")
var ErrInvalidStore = errors.New("outbox store requires a database and table name")

// Message is an immutable event prepared by the domain producer.
type Message struct {
    EventID string
    Topic string
    Payload []byte
    Metadata map[string]string
}

// Entry represents one durable, individually retriable delivery.
type Entry struct {
    EventID string `gorm:"column:event_id;primaryKey"`
    Topic string `gorm:"column:topic"`
    Payload []byte `gorm:"column:payload"`
    Metadata []byte `gorm:"column:metadata"`
    Status string `gorm:"column:status"`
    Attempts int `gorm:"column:attempts"`
    NextAttemptAt time.Time `gorm:"column:next_attempt_at"`
    LeaseUntil *time.Time `gorm:"column:lease_until"`
    CreatedAt time.Time `gorm:"column:created_at"`
    PublishedAt *time.Time `gorm:"column:published_at"`
}

type Store interface {
    Enqueue(ctx context.Context, message Message) error
    Claim(ctx context.Context, limit int, lease time.Duration) ([]Entry, error)
    MarkSent(ctx context.Context, eventID string) error
    MarkRetry(ctx context.Context, eventID string, attempt int) error
}

type GormStore struct {
    db *gorm.DB
    table string
}

func NewGormStore(db *gorm.DB, table string) *GormStore { return &GormStore{db: db, table: table} }

// Enqueue uses transaction.Database so a failed insert rolls back Person writes.
func (store *GormStore) Enqueue(ctx context.Context, message Message) error {
    if store.db == nil || store.table == "" { return ErrInvalidStore }
    if message.EventID == "" || message.Topic == "" || len(message.Payload) == 0 { return ErrInvalidMessage }
    metadata, err := json.Marshal(message.Metadata)
    if err != nil { return err }
    return coretransaction.Database(ctx, store.db).Table(store.table).Create(&Entry{
        EventID: message.EventID, Topic: message.Topic, Payload: message.Payload,
        Metadata: metadata, Status: Pending, NextAttemptAt: time.Now().UTC(),
    }).Error
}

// Claim atomically leases pending or expired in-flight events for multiple dispatchers.
func (store *GormStore) Claim(ctx context.Context, limit int, lease time.Duration) ([]Entry, error) {
    if store.db == nil || store.table == "" { return nil, ErrInvalidStore }
    if limit <= 0 { return nil, nil }
    if lease <= 0 { lease = 30*time.Second }
    now := time.Now().UTC()
    until := now.Add(lease)
    var entries []Entry
    err := store.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        if err := tx.Table(store.table).Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
            Where("(status = ? AND next_attempt_at <= ?) OR (status = ? AND lease_until <= ?)", Pending, now, Sending, now).
            Order("created_at, event_id").Limit(limit).Find(&entries).Error; err != nil { return err }
        for index := range entries {
            if err := tx.Table(store.table).Where("event_id = ?", entries[index].EventID).
                Updates(map[string]any{"status": Sending, "lease_until": until, "attempts": gorm.Expr("attempts + 1")}).Error; err != nil { return err }
            entries[index].Attempts++
        }
        return nil
    })
    return entries, err
}

func (store *GormStore) MarkSent(ctx context.Context, eventID string) error {
    now := time.Now().UTC()
    return store.db.WithContext(ctx).Table(store.table).Where("event_id = ? AND status = ?", eventID, Sending).
        Updates(map[string]any{"status": Sent, "lease_until": nil, "published_at": now}).Error
}

func (store *GormStore) MarkRetry(ctx context.Context, eventID string, attempt int) error {
    seconds := 1
    for index := 1; index < attempt && seconds < 64; index++ { seconds *= 2 }
    next := time.Now().UTC().Add(time.Duration(seconds)*time.Second)
    return store.db.WithContext(ctx).Table(store.table).Where("event_id = ? AND status = ?", eventID, Sending).
        Updates(map[string]any{"status": Pending, "lease_until": nil, "next_attempt_at": next}).Error
}
