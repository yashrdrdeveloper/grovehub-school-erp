package helpers

import (
	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	"grovehub/services/integration/internal/app/entity"
)


func ToResponse(connection entity.Connection) connectiondto.Connection {
	return connectiondto.Connection{
		ID:             connection.ID,
		TenantID:       connection.TenantID,
		OwnerGID:       connection.OwnerGID,
		Provider:       connection.Provider,
		Status:         connection.Status,
		ExternalID:     connection.ExternalID,
		LastTestedAt:   connection.LastTestedAt,
		LastTestStatus: connection.LastTestStatus,
		CreatedAt:      connection.CreatedAt,
		UpdatedAt:      connection.UpdatedAt,
	}
}
