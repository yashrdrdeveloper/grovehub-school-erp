package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	ServiceDown = coreerrors.NewDefinition(
		"IDENTITY_SERVICE_DOWN",
		"Identity service is unavailable",
	)
)
