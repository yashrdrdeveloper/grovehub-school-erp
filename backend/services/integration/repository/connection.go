package repository

import (
 "context"
 "errors"

 "grovehub/services/integration/entity"
)

var ErrConnectionNotFound = errors.New("integration connection not found")

type ConnectionRepository interface {
 Create(ctx context.Context, connection entity.Connection) error
 Get(ctx context.Context, tenantID, id string) (entity.Connection, error)
 List(ctx context.Context, tenantID string) ([]entity.Connection, error)
 Disconnect(ctx context.Context, tenantID, id string) (entity.Connection, error)
}
