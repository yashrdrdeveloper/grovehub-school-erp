package response

import "time"

type Connection struct {
 ID string `json:"id"`
 TenantID string `json:"tenantId"`
 OwnerGID string `json:"ownerGid"`
 Provider string `json:"provider"`
 Status string `json:"status"`
 CreatedAt time.Time `json:"createdAt"`
 UpdatedAt time.Time `json:"updatedAt"`
}
