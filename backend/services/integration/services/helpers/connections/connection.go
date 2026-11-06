package connections

import (
 "strings"

 "github.com/google/uuid"

 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
)

// ValidID enforces the platform's UUIDv7 identity contract.
func ValidID(raw string)bool{
 parsed,err:=uuid.Parse(strings.TrimSpace(raw))
 return err==nil&&parsed.Version()==7
}

// ToResponse prevents GORM entities from leaking into transport responses.
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
