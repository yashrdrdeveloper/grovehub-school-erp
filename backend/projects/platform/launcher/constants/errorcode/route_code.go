package errorcode

import coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"

var (
	RouteNotFound = coreerrors.NewDefinition(
		"PLATFORM_ROUTE_NOT_FOUND",
		"API route not found",
	)
)
