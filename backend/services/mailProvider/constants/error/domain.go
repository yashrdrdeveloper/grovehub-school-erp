package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
	InvalidDomain = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_INVALID",
		"Mail provider domain is invalid",
	)
	DomainAlreadyExists = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_ALREADY_EXISTS",
		"Mail provider domain already exists",
	)
	DomainNotFound = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_NOT_FOUND",
		"Mail provider domain was not found",
	)
	DomainNotVerified = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_NOT_VERIFIED",
		"Mail provider domain is not verified",
	)
	DomainCreateFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_CREATE_FAILED",
		"Unable to create mail provider domain",
	)
	DomainReadFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_READ_FAILED",
		"Unable to read mail provider domain",
	)
	DomainUpdateFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_UPDATE_FAILED",
		"Unable to update mail provider domain",
	)
	DomainVerificationFailed = coreerrors.NewDefinition(
		"MAIL_PROVIDER_DOMAIN_VERIFICATION_FAILED",
		"Unable to verify mail provider domain",
	)
)
