package helpers

import (
    "errors"

    coreerrors "grovehub/libs/core/errors"
    personerror "grovehub/services/foundation/constants/error"
    foundationmappers "grovehub/services/foundation/mappers/person"
    repositoryerrors "grovehub/services/foundation/repository/constants/errors"
)

// MapUpdateError preserves Person mutation error codes across operations.
func MapUpdateError(gid string, err error, mapper foundationmappers.PersonMapper) error {
	mapped := map[error]func() error{
		repositoryerrors.ErrPersonNotFound: func() error {
			return coreerrors.NewErrorWithMetadata(
				personerror.PersonNotFound,
				mapper.ToPersonErrorMetadata(gid),
			)
		},
		repositoryerrors.ErrContactAlreadyExists: func() error {
			return coreerrors.NewError(personerror.PersonContactConflict)
		},
		repositoryerrors.ErrPersonVersionConflict: func() error {
			return coreerrors.NewErrorWithMetadata(
				personerror.PersonConcurrentModification,
				mapper.ToPersonErrorMetadata(gid),
			)
		},
	}
	for repositoryError, factory := range mapped {
		if errors.Is(err, repositoryError) {
			return factory()
		}
	}
	return coreerrors.WrapWithMetadata(
		personerror.PersonUpdateFailed,
		err,
		mapper.ToPersonErrorMetadata(gid),
	)
}
