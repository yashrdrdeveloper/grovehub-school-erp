package helpers

import (
    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
)

// FailRead maps a repository read failure for the Person Create workflow.
func FailRead(err error) corepipeline.Outcome[persondto.Person] {
    return corepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonReadFailed, err))
}

// FailCreate maps a creation failure for the Person Create workflow.
func FailCreate(err error) corepipeline.Outcome[persondto.Person] {
    return corepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonCreateFailed, err))
}
