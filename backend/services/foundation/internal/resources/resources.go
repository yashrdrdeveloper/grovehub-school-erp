// Package resources contains the assets embedded into the Foundation binary.
package resources

import "embed"

// Migrations is bundled into every production binary importing Foundation.
//go:embed migrations/*.sql
var Migrations embed.FS
