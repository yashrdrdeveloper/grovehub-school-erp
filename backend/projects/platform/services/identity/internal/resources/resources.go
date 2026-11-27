package resources
import "embed"
//go:embed migrations/*.sql
var Migrations embed.FS
//go:embed api/openapi/v1/*.yaml
var OpenAPI embed.FS
//go:embed templates/registration/*.html
var Templates embed.FS
