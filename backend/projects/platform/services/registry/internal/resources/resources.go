// Package resources owns the Registry's embedded migrations.
package resources

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS
