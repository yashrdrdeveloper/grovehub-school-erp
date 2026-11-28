package activity

import (
	"context"
	"errors"

	"grovehub/services/audit/internal/app/events/consumers/activity/constants"
	"grovehub/services/audit/internal/app/events/consumers/activity/helpers"
	"grovehub/services/audit/internal/app/repository/activity"
	"grovehub/services/audit/internal/app/repository/rejection"
	"grovehub/services/infrastructure/client/messaging"
)

type Consumer struct {
	repository activity.Repository
	rejected   rejection.Repository
}

func New(repository activity.Repository, rejected rejection.Repository) *Consumer {
	return &Consumer{repository: repository, rejected: rejected}
}

func (consumer *Consumer) Handle(ctx context.Context, eventType string, message messaging.Message) error {
	_, err := consumer.Process(ctx, eventType, message)
	return err
}

// Process distinguishes persisted activities from quarantined events.
func (consumer *Consumer) Process(ctx context.Context, eventType string, message messaging.Message) (string, error) {
	entry, err := helpers.Project(ctx, eventType, message.Payload)
	if errors.Is(err, helpers.ErrInvalidEvent) {
		return constants.InvalidEvent, consumer.reject(ctx, eventType, message.Payload, constants.InvalidEventReason, err)
	}
	if err != nil {
		return constants.Retry, err
	}

	err = consumer.repository.Insert(ctx, entry)
	if errors.Is(err, activity.ErrConflict) {
		return constants.IdentityConflict, consumer.reject(ctx, eventType, message.Payload, constants.IdentityConflictReason, err)
	}

	return constants.Accepted, err
}

func (consumer *Consumer) reject(ctx context.Context, eventType string, payload []byte, reason string, original error) error {
	if consumer.rejected == nil {
		return original
	}

	return consumer.rejected.Record(ctx, eventType, payload, reason)
}
