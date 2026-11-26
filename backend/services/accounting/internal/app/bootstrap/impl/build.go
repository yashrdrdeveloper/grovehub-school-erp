package impl

import (
 "errors"
 "log/slog"

 "gorm.io/gorm"
 "grovehub/services/accounting/internal/app/bootstrap"
 "grovehub/services/infrastructure/client/messaging"
)

type builder struct {
 logger *slog.Logger
 db *gorm.DB
 events messaging.EventStore
}

var _ bootstrap.Bootstrap = (*builder)(nil)

func New(logger *slog.Logger, db *gorm.DB, events messaging.EventStore) bootstrap.Bootstrap {
 return &builder{logger: logger, db: db, events: events}
}

func (b *builder) Build() (bootstrap.Runtime, error) {
 if b.db == nil || b.events == nil {
  return nil, errors.New("accounting requires PostgreSQL and Infrastructure")
 }
 if b.logger == nil { b.logger = slog.Default() }
 repositories := b.newRepository()
 return &runtime{services: b.newServices(repositories)}, nil
}
