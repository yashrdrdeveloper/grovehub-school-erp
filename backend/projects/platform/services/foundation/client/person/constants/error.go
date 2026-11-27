package constants

import personerrors "grovehub/services/foundation/internal/app/constants/error"

// Public Person client error contracts. The definitions remain owned by Foundation.
var (
    InvalidPersonGID = personerrors.InvalidPersonGID
    InvalidPersonRequest = personerrors.InvalidPersonRequest
    InvalidPersonContact = personerrors.InvalidPersonContact
    PersonNotFound = personerrors.PersonNotFound
    PersonContactConflict = personerrors.PersonContactConflict
)
