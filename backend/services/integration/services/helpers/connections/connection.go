package connections

import (
 corenormalize "grovehub/libs/core/normalize"
 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
)

// ValidID reuses the platform's canonical UUIDv7 validator.
func ValidID(raw string)bool{
 _,err:=corenormalize.UUIDv7(raw)
 return err==nil
}

func ToResponse(connection entity.Connection)integrationresponse.Connection{
 return integrationresponse.Connection{
  ID:connection.ID,
  TenantID:connection.TenantID,
  OwnerGID:connection.OwnerGID,
  Provider:connection.Provider,
  Status:connection.Status,
  ExternalID:connection.ExternalID,
  LastTestedAt:connection.LastTestedAt,
  LastTestStatus:connection.LastTestStatus,
  CreatedAt:connection.CreatedAt,
  UpdatedAt:connection.UpdatedAt,
 }
}
