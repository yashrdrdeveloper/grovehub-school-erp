// Package resources embeds Integration SQL and its OpenAPI description.
package resources

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS

//go:embed api/openapi/v1/*.yaml
var OpenAPI embed.FS
