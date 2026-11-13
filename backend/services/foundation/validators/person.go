package validators

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"
	foundationrequest "grovehub/services/foundation/dto/request"
)

type PersonValidator interface {
	CreateRules() annotationvalidation.Rules[foundationrequest.CreatePersonRequest]
	UpdateRules() annotationvalidation.Rules[foundationrequest.UpdatePersonRequest]
	SearchRules() annotationvalidation.Rules[foundationrequest.SearchPersonsRequest]
}
