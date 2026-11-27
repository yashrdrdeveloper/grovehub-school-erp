package connection

import (
    dto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/entity"
)

// @Mapper ConnectionMapper
type ConnectionMapper interface {
    ToResponse(connection entity.Connection) dto.Connection
    TestToResponse(test entity.ConnectionTest) dto.ConnectionTest
}
