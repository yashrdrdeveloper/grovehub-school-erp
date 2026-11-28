package impl

import (
	"context"

	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"
)

func NewOperatorBootstrap(migrate func(context.Context, *gorm.DB) error) platformbootstrap.Service[*gorm.DB] {
	return platformbootstrap.Registration[*gorm.DB]{
		ServiceName: "foundation",
		MigrationPlan: platformbootstrap.Migration[*gorm.DB]{
			Target: platformbootstrap.BusinessDatabase,
			Run:    migrate,
		},
	}
}
