package changed

import resource "grovehub/services/identity/internal/app/events/producers/resource/operations/changed"

type Changed struct {
	TenantID string `json:"tenantId"`
	ActorID string `json:"actorId"`
	ResourceType string `json:"resourceType"`
	Operation string `json:"operation"`
	Old *resource.Snapshot `json:"old"`
	New resource.Snapshot `json:"new"`
}
