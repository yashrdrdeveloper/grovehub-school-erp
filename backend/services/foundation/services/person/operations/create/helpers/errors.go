package helpers

import (
    servicepipeline "grovehub/libs/core/pipeline/service"
    coreerrors "grovehub/libs/core/errors"
    personerror "grovehub/services/foundation/constants/error"
    persondto "grovehub/services/foundation/dto/person"
)

// FailRead maps a repository read failure for the Person Create workflow.
func FailRead(err error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonReadFailed, err))
}

// FailCreate maps a creation failure for the Person Create workflow.
func FailCreate(err error) servicepipeline.Outcome[persondto.Person] {
    return servicepipeline.Fail[persondto.Person](coreerrors.Wrap(personerror.PersonCreateFailed, err))
}
