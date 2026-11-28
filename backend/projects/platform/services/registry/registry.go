// Package registry exposes the generic project registry bootstrap.
package registry

import (
	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/registry/internal/app/bootstrap/impl"
)

func Bootstrap() platformbootstrap.Service[*gorm.DB] {
	return impl.New()
}
