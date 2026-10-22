package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	RouteNotFound = coreerrors.NewDefinition(
		"PLATFORM_ROUTE_NOT_FOUND",
		"API route not found",
	)
)
