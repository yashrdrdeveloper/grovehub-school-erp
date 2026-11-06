package request

type CreateConnection struct {
 Provider string `json:"provider"`
 OwnerGID string `json:"ownerGid"`
}
