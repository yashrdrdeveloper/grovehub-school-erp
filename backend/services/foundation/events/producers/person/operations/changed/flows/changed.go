// Package flows contains event pipeline inputs that do not depend on transport.
package flows

import (
    foundationevents "grovehub/services/foundation/constants/events"
    persondto "grovehub/services/foundation/dto/person"
)

type Change struct {
    Operation foundationevents.Operation
    Old *persondto.Person
    New *persondto.Person
}
