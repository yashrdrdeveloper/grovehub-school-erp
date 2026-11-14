// Package resources contains the assets embedded into the Foundation binary.
package resources

import "embed"

// Migrations is bundled into every production binary importing Foundation.
//go:embed migrations/*.sql
var Migrations embed.FS

// OpenAPI embeds the Foundation HTTP specification for use without source files.
//go:embed api/openapi/v1/*.yaml
var OpenAPI embed.FS
