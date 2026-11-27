// Package resources embeds Infrastructure migration SQL as deployment assets.
package resources

import "embed"

//go:embed migrations/*.sql
var Migrations embed.FS
