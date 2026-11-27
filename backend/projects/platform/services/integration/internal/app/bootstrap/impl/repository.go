package impl

import (
    "grovehub/services/audit"
    "grovehub/services/infrastructure"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    connectionimpl "grovehub/services/integration/internal/app/repository/connection/impl"
    watchrepo "grovehub/services/integration/internal/app/repository/watch"
    watchimpl "grovehub/services/integration/internal/app/repository/watch/impl"
)

func (b *builder) newRepository() (connectionrepo.ConnectionRepository, watchrepo.WatchRepository) {
    return connectionimpl.NewConnectionRepository(b.db,audit.NewProducer(infrastructure.NewEventStore(b.db))), watchimpl.NewWatchRepository(b.db)
}
