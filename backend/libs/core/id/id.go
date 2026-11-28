// Package id exposes domain-independent identifier generation.
package id

import "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/id/internal/src/operations/uuid"

// NewUUIDv7 generates a time-ordered UUID version 7 string.
func NewUUIDv7() (string, error) {
	return uuid.NewUUIDv7()
}
