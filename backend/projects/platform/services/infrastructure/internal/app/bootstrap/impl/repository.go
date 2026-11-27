package impl

import (
    mapperimpl "grovehub/services/infrastructure/internal/app/mappers/messaging/impl"
    repository "grovehub/services/infrastructure/internal/app/repository/messaging"
    repositoryimpl "grovehub/services/infrastructure/internal/app/repository/messaging/impl"
    validatorimpl "grovehub/services/infrastructure/internal/app/validators/messaging/impl"
)

func (b *builder) newRepository() repository.Repository {
    return repositoryimpl.New(b.db, validatorimpl.New(), mapperimpl.New())
}
