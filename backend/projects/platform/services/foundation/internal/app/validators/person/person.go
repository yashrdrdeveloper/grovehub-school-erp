package person

import (
	coreannotation "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation"
	persondto "grovehub/services/foundation/internal/app/dto/person"
)

type PersonValidator interface {
	CreateRules() coreannotation.ValidationRules[persondto.CreatePersonRequest]
	UpdateRules() coreannotation.ValidationRules[persondto.UpdatePersonRequest]
	SearchRules() coreannotation.ValidationRules[persondto.SearchPersonsRequest]
}
