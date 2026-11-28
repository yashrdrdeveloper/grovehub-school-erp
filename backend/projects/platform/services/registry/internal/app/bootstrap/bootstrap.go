package bootstrap

import (
	"gorm.io/gorm"

	platformbootstrap "grovehub/libs/platform/bootstrap"
)

// Service exposes the Registry's currently available bootstrap capabilities.
// The HTTP runtime will be added with the authenticated Registry API.
type Service = platformbootstrap.Service[*gorm.DB]
