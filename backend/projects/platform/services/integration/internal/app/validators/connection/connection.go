package connection

import dto "grovehub/services/integration/internal/app/dto/connection"

// ConnectionValidator validates immutable request IDs and creation input.
type ConnectionValidator interface {
    ValidateCreate(tenantID string, request dto.CreateConnection) error
    ValidatePair(tenantID string, resourceID string) error
}
