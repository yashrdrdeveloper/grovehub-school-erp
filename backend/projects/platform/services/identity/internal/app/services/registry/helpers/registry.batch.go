package helpers

import (
	"sort"

	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	registryrepo "grovehub/services/identity/internal/app/repository/registry"
	"grovehub/services/identity/internal/app/services/registry"
)

const maxBatchSize = 100

func PrepareBatch(request registry.Request) (registryrepo.Batch, error) {
	size := len(request.Features) + len(request.Permissions)
	if size == 0 || size > maxBatchSize {
		return registryrepo.Batch{}, registry.ErrInvalidBatch
	}
	batch := registryrepo.Batch{}
	features := make(map[feature.Code]bool)
	for _, code := range request.Features {
		definition, known := feature.Find(code)
		if !known || features[code] {
			return registryrepo.Batch{}, registry.ErrInvalidBatch
		}
		features[code] = true
		batch.Features = append(batch.Features, definition)
	}
	permissions := make(map[permission.Code]bool)
	for _, code := range request.Permissions {
		definition, known := permission.Find(code)
		if !known || permissions[code] {
			return registryrepo.Batch{}, registry.ErrInvalidBatch
		}
		permissions[code] = true
		batch.Permissions = append(batch.Permissions, definition)
	}
	sort.Slice(batch.Features, func(i, j int) bool { return batch.Features[i].Code < batch.Features[j].Code })
	sort.Slice(batch.Permissions, func(i, j int) bool { return batch.Permissions[i].Code < batch.Permissions[j].Code })
	return batch, nil
}
