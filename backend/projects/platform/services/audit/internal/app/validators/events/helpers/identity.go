package helpers

import (
	"encoding/json"

	"github.com/google/uuid"

	"grovehub/services/audit/internal/app/dto/producer"
)

const NamespaceID = "bca86488-355c-4c40-a2cf-41121835ed83"
const VersionedSource = "versioned-source"

func VersionEventID(data producer.Activity) string {
	key, _ := json.Marshal([]any{VersionedSource, data.TenantID, data.Source, data.AssetType, data.AssetID, data.Operation, data.AssetVersion})
	return uuid.NewSHA1(uuid.MustParse(NamespaceID), key).String()
}
