package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidDKIMIdentity = coreerrors.NewDefinition(
		"INFRASTRUCTURE_DKIM_INVALID",
		"Infrastructure DKIM identity is invalid",
	)
	DKIMSelectorConflict = coreerrors.NewDefinition(
		"INFRASTRUCTURE_DKIM_SELECTOR_CONFLICT",
		"Infrastructure DKIM selector conflicts with the persisted identity",
	)
	DKIMPrepareFailed = coreerrors.NewDefinition(
		"INFRASTRUCTURE_DKIM_PREPARE_FAILED",
		"Unable to prepare infrastructure DKIM identity",
	)
)
