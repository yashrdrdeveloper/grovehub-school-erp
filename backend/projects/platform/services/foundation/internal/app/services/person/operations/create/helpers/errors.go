package helpers

import (
    servicepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/service"
    coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
    personerror "grovehub/services/foundation/internal/app/constants/error"
    persondto "grovehub/services/foundation/internal/app/dto/person"
)

// FailRead maps a repository read failure for the Person Create workflow.
func FailRead(err error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonReadFailed, err))
}

// FailCreate maps a creation failure for the Person Create workflow.
func FailCreate(err error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonCreateFailed, err))
}
