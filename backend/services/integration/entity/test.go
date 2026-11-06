package entity

import "time"

type ConnectionTest struct {
 ID string `gorm:"column:id;type:uuid;primaryKey"`
 ConnectionID string `gorm:"column:connection_id;type:uuid;not null;index"`
 Direction string `gorm:"column:direction;not null"`
 Status string `gorm:"column:status;not null"`
 Detail string `gorm:"column:detail;not null"`
 StartedAt time.Time `gorm:"column:started_at;not null"`
 CompletedAt time.Time `gorm:"column:completed_at;not null"`
}
func (ConnectionTest) TableName() string { return "integration.connection_tests" }
