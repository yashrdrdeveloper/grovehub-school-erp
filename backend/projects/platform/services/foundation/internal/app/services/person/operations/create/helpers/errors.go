package helpers

import (
    corepipeline "grovehub/libs/core/pipeline"
    coreerrors "grovehub/libs/core/errors"
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
