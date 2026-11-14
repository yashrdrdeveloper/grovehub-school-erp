package person

import (
	annotationvalidation "grovehub/libs/core/annotation/validation"
	persondto "grovehub/services/foundation/internal/app/dto/person"
)

type PersonValidator interface {
	CreateRules() annotationvalidation.Rules[persondto.CreatePersonRequest]
	UpdateRules() annotationvalidation.Rules[persondto.UpdatePersonRequest]
	SearchRules() annotationvalidation.Rules[persondto.SearchPersonsRequest]
}
