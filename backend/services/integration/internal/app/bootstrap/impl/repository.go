package impl

import (
    repository "grovehub/services/integration/internal/app/repository"
    repositoryimpl "grovehub/services/integration/internal/app/repository/impl"
)

func (b *builder) newRepository() (repository.ConnectionRepository,repository.WatchRepository) {
    return repositoryimpl.NewConnectionRepository(b.db),repositoryimpl.NewWatchRepository(b.db)
}
