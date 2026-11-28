package entity

import (
	"time"
	"grovehub/services/registry/internal/app/constants/status"
)

// @Entity registry.projects
type Project struct {
	ProjectID string        `gorm:"column:project_id;primaryKey;type:uuid"`
	Code      string        `gorm:"column:code;not null"`
	Name      string        `gorm:"column:name;not null"`
	Status    status.Status `gorm:"column:status;not null"`
	Version   int64         `gorm:"column:version;not null"`
	CreatedAt time.Time     `gorm:"column:created_at;not null"`
	UpdatedAt time.Time     `gorm:"column:updated_at;not null"`
}
