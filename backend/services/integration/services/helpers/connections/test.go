package connections

import (
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
)

func TestToResponse(test entity.ConnectionTest)integrationresponse.ConnectionTest{
 return integrationresponse.ConnectionTest{
  ID:test.ID,
  ConnectionID:test.ConnectionID,
  Direction:test.Direction,
  Status:test.Status,
  Detail:test.Detail,
  StartedAt:test.StartedAt,
  CompletedAt:test.CompletedAt,
 }
}
