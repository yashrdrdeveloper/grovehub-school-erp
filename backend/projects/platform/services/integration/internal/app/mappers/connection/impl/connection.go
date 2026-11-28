package impl

import (
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/entity"
    mapper "grovehub/services/integration/internal/app/mappers/connection"
)

type connectionMapper struct{}

var _ mapper.ConnectionMapper = (*connectionMapper)(nil)

func NewConnectionMapper() mapper.ConnectionMapper {return &connectionMapper{}}

func (m *connectionMapper) ToResponse(connection entity.Connection) connectiondto.Connection {
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

func (m *connectionMapper) TestToResponse(test entity.ConnectionTest) connectiondto.ConnectionTest {
	return connectiondto.ConnectionTest{
		ID:           test.ID,
		ConnectionID: test.ConnectionID,
		Direction:    test.Direction,
		Status:       test.Status,
		Detail:       test.Detail,
		StartedAt:    test.StartedAt,
		CompletedAt:  test.CompletedAt,
	}
}
