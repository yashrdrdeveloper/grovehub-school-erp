// Package logging is a compatibility facade.
// Deprecated: use the root core/http package for new code.
package logging

import ("log/slog"; corehttp "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http")
func LevelForStatus(status int) slog.Level { return corehttp.LevelForStatus(status) }
