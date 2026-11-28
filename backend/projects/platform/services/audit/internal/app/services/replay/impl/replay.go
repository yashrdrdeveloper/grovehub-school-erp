package impl

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	coredatabase "grovehub/libs/core/database"
	"grovehub/services/audit/internal/app/constants/events"
	"grovehub/services/audit/internal/app/events/consumers/activity/helpers"
	"grovehub/services/audit/internal/app/repository/activity"
	"grovehub/services/audit/internal/app/services/replay"
	"grovehub/services/infrastructure/client/messaging"
)

var errInvalid = errors.New("invalid replay request or event outside scope")

type Service struct {
	db         *gorm.DB
	source     messaging.ReplaySource
	repository activity.Repository
}

func New(db *gorm.DB, source messaging.ReplaySource, repository activity.Repository) replay.Service {
	return &Service{db: db, source: source, repository: repository}
}

func (service *Service) Replay(ctx context.Context, input replay.Request) (replay.Result, error) {
	id, err := uuid.Parse(input.EventID)
	if err != nil || id == uuid.Nil || id.String() != input.EventID || input.TenantID == "" || len(input.TenantID) > 128 || strings.TrimSpace(input.TenantID) != input.TenantID || (input.Apply && (strings.TrimSpace(input.OperatorLabel) == "" || len(input.OperatorLabel) > 128 || strings.TrimSpace(input.OperatorLabel) != input.OperatorLabel)) || service.db == nil || service.source == nil || service.repository == nil || coredatabase.TransactionActive(ctx) {
		return replay.Result{}, errInvalid
	}
	event, err := service.source.Published(ctx, input.EventID)
	if err != nil {
		return replay.Result{}, err
	}
	eventType, known := events.Topics()[event.Topic]
	if !known || event.EventID != input.EventID {
		return replay.Result{}, errInvalid
	}
	entry, err := helpers.Project(ctx, eventType, event.Payload)
	if err != nil || entry.TenantID != input.TenantID || entry.SourceEventID != input.EventID {
		return replay.Result{}, errInvalid
	}
	result := replay.Result{EventID: input.EventID, AuditID: entry.AuditID, Operation: entry.Operation}
	if !input.Apply {
		return result, nil
	}
	replayID := uuid.NewString()
	err = coredatabase.RunTransaction(ctx, service.db, func(tx context.Context) error {
		if err := service.repository.Insert(tx, entry); err != nil {
			return err
		}
		return coredatabase.TransactionDB(tx, service.db).Exec("INSERT INTO audit.replay_runs(replay_id,tenant_id,source_event_id,audit_id,operator_label) VALUES (?,?,?,?,?)", replayID, input.TenantID, input.EventID, entry.AuditID, input.OperatorLabel).Error
	})
	if err != nil {
		return replay.Result{}, err
	}
	result.Applied = true
	result.ReplayID = replayID
	return result, nil
}
