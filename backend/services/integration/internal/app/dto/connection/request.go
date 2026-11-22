package connection

type CreateConnection struct {
	Provider string `json:"provider"`
	OwnerGID string `json:"ownerGid"`
}

type ConfigureSystemMailSender struct {
	ConnectionID string `json:"connectionId"`
}
