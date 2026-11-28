package bootstrap

import (
	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"
	"grovehub/services/registry/client/authorization"
)

type Dependencies struct {
	DB       *gorm.DB
	Audience string
}

type Runtime interface {
	Authorization() authorization.Client
}

type Service = platformbootstrap.RuntimeService[*gorm.DB, Dependencies, Runtime]
