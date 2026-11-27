package impl

import (
	"context"
	"errors"

	"gorm.io/gorm"
	client "grovehub/services/infrastructure/client/messaging"
	"grovehub/services/infrastructure/internal/app/entity"
)

type replaySource struct { db *gorm.DB }

func NewReplaySource(db *gorm.DB) client.ReplaySource {
	return &replaySource{db:db}
}

func (source *replaySource) Published(ctx context.Context,eventID string) (client.Event,error) {
	if source.db==nil { return client.Event{},client.ErrReplayUnavailable }
	var row entity.Messaging
	err:=source.db.WithContext(ctx).Select("event_id","topic","payload").Where("event_id = ? AND status = ?",eventID,entity.MessagingSent).Take(&row).Error
	if errors.Is(err,gorm.ErrRecordNotFound) { return client.Event{},client.ErrReplayUnavailable }
	if err!=nil { return client.Event{},err }
	return client.Event{EventID:row.EventID,Topic:row.Topic,Payload:row.Payload},nil
}
