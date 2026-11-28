package entity

import "grovehub/libs/platform/scope"

// @Entity registry.capabilities
type Capability struct {
	Code  string      `gorm:"column:code;primaryKey"`
	Scope scope.Level `gorm:"column:scope;primaryKey"`
}
