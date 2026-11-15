package helpers

import (
	connectiondto "grovehub/services/integration/internal/app/dto/connection"
	"grovehub/services/integration/internal/app/entity"
)

func TestToResponse(test entity.ConnectionTest) connectiondto.ConnectionTest {
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
