package person

import (
	coreannotation "grovehub/libs/core/annotation"
	persondto "grovehub/services/foundation/internal/app/dto/person"
)

type PersonValidator interface {
	CreateRules() coreannotation.ValidationRules[persondto.CreatePersonRequest]
	UpdateRules() coreannotation.ValidationRules[persondto.UpdatePersonRequest]
	SearchRules() coreannotation.ValidationRules[persondto.SearchPersonsRequest]
}
