package impl

import (
    dto "grovehub/services/integration/internal/app/dto/connection"
    errorcode "grovehub/services/integration/internal/app/services/connection/helpers"
    "grovehub/services/integration/internal/app/validators/connection"
    "grovehub/services/integration/internal/app/validators/connection/helpers"
)

type connectionValidator struct {}

var _ connection.ConnectionValidator = (*connectionValidator)(nil)

func NewConnectionValidator() connection.ConnectionValidator {return &connectionValidator{}}

func (v *connectionValidator) ValidatePair(tenantID string, resourceID string) error {
    if !helpers.ValidTenantID(tenantID)||!helpers.ValidID(resourceID) {return errorcode.ErrInvalidConnection}
    return nil
}

func (v *connectionValidator) ValidateCreate(tenantID string, request dto.CreateConnection) error {
    return v.ValidatePair(tenantID,request.OwnerGID)
}
