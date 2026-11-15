package impl

import (
	"grovehub/services/integration/internal/app/providers"
		"strings"
    contracts "grovehub/services/integration/internal/app/providers/contracts"
    dto "grovehub/services/integration/internal/app/providers/dto"
    helpers "grovehub/services/integration/internal/app/providers/helpers"
    resolve "grovehub/services/integration/internal/app/providers/operations/resolve"
    catalog "grovehub/services/integration/internal/app/providers/operations/catalog"
)

type registry struct{ items map[string]providers.Provider }

func NewRegistry(items ...providers.Provider) (contracts.Registry, error) {
	r := &registry{items: make(map[string]providers.Provider, len(items))}
	for _, item := range items {
		descriptor := item.Descriptor()
		id := strings.ToUpper(strings.TrimSpace(descriptor.ID))
		if id == "" {
			return nil, helpers.ErrUnknownProvider
		}
		if _, exists := r.items[id]; exists {
			return nil, helpers.ErrDuplicateProvider
		}
		r.items[id] = item
	}
	return r, nil
}

func (r *registry) Find(id string) (providers.Provider,error) {
    return resolve.Find(r.items,id)
}

func (r *registry) List() []dto.Descriptor {
    return catalog.List(r.items)
}
