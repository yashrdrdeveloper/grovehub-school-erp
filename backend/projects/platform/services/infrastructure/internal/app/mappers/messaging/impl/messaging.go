package impl

import (
    "encoding/json"
    "time"

    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
    "grovehub/services/infrastructure/internal/app/entity"
    mapper "grovehub/services/infrastructure/internal/app/mappers/messaging"
)

type messagingMapper struct{}

var _ mapper.Mapper = (*messagingMapper)(nil)

func New() mapper.Mapper { return &messagingMapper{} }

func (mapper *messagingMapper) ToEntity(event dto.Event, now time.Time) (entity.Messaging,error) {
    metadata,err:=json.Marshal(event.Metadata)
    if err!=nil{return entity.Messaging{},err}
    return entity.Messaging{
        EventID:event.EventID,
        Topic:event.Topic,
        Payload:event.Payload,
        Metadata:json.RawMessage(metadata),
        Status:entity.MessagingPending,
        NextAttemptAt:now,
    },nil
}
