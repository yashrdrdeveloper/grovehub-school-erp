package person

import (
	annotationvalidation "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/validation"
	persondto "grovehub/services/foundation/internal/app/dto/person"
)

type PersonValidator interface {
	CreateRules() annotationvalidation.Rules[persondto.CreatePersonRequest]
	UpdateRules() annotationvalidation.Rules[persondto.UpdatePersonRequest]
	SearchRules() annotationvalidation.Rules[persondto.SearchPersonsRequest]
}
