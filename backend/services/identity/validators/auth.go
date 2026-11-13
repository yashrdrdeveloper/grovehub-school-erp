package validators

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"
	identityrequest "grovehub/services/identity/dto/request"
)

type AuthValidator interface {
	RegisterRules() annotationvalidation.Rules[identityrequest.RegisterRequest]
}
