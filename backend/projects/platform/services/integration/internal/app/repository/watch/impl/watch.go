package impl

import (
	"context"
 "fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	connectionconstants "grovehub/services/integration/internal/app/constants/connections"
	providerconstants "grovehub/services/integration/internal/app/constants/providers"
	"grovehub/services/integration/internal/app/entity"
)

// WatchRepository uses GORM's model, association, and upsert APIs.
// There are no manual JOIN/SELECT/SCAN expressions in this repository.
type WatchRepository struct{ db *gorm.DB }

func NewWatchRepository(db *gorm.DB) *WatchRepository {
	return &WatchRepository{db: db}
}

func (r *WatchRepository) SaveWatch(ctx context.Context, connectionID, historyID string, expires time.Time) error {
	entry := entity.WebhookSubscription{
		ID:           uuid.NewString(),
		ConnectionID: connectionID,
		Resource:     connectionconstants.GmailMailboxResource,
		ExternalID:   historyID,
		Status:       connectionconstants.Active,
		ExpiresAt:    &expires,
	}
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "connection_id"}, {Name: "resource"}},
		DoUpdates: clause.AssignmentColumns([]string{"external_id", "status", "expires_at"}),
	}).Create(&entry).Error
}

func (r *WatchRepository) GetGmailByMailbox(ctx context.Context, email string) ([]entity.Connection, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	rows := make([]entity.Connection, 0)
	// Provider account addresses are normalized to lowercase during OAuth.
	criteria := entity.Connection{Provider: providerconstants.Gmail, Status: entity.StatusActive, ExternalID: &email}
	err := r.db.WithContext(ctx).Model(&entity.Connection{}).Where(&criteria).Find(&rows).Error
	return rows, err
}

func (r *WatchRepository) ExpiringWatches(ctx context.Context, before time.Time) ([]entity.Connection, error) {
	subscriptions := make([]entity.WebhookSubscription, 0)
	err := r.db.WithContext(ctx).Model(&entity.WebhookSubscription{}).
		Where(&entity.WebhookSubscription{Resource: connectionconstants.GmailMailboxResource, Status: connectionconstants.Active}).
		Where("expires_at < ?", before).
		Preload("Connection").
		Limit(500).
		Find(&subscriptions).Error
	if err != nil {
		return nil, err
	}

	result := make([]entity.Connection, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		connection := subscription.Connection
		if connection.ID != "" && connection.Provider == providerconstants.Gmail && connection.Status == entity.StatusActive {
			result = append(result, connection)
		}
	}
	return result, nil
}

// ExpiringWatchesAfter scans active Gmail connections in stable keyset order.
// Filter before LIMIT so disconnected rows cannot starve eligible watches.
func(r *WatchRepository) ExpiringWatchesAfter(ctx context.Context,before time.Time,after string,limit int)([]entity.Connection,error){
 if limit<1||limit>500{return nil,fmt.Errorf("invalid watch page size")}
 rows:=make([]entity.Connection,0,limit)
 query:=r.db.WithContext(ctx).Table("integration.connections AS c").Select("c.*").
  Joins("JOIN integration.webhook_subscriptions AS w ON w.connection_id=c.id").
  Where("c.provider=? AND c.status=? AND w.resource=? AND w.status=? AND w.expires_at<?",providerconstants.Gmail,entity.StatusActive,connectionconstants.GmailMailboxResource,connectionconstants.Active,before)
 if after!=""{query=query.Where("c.id > ?",after)}
 err:=query.Order("c.id").Limit(limit).Find(&rows).Error
 return rows,err
}
