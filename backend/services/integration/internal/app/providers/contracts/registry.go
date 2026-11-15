package contracts

import (
    providers "grovehub/services/integration/internal/app/providers"
    dto "grovehub/services/integration/internal/app/providers/dto"
)

type Registry interface {
    Find(id string) (providers.Provider,error)
    List() []dto.Descriptor
}
