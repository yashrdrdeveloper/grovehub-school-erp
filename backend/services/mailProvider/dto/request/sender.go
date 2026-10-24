package request

type RegisterSender struct {
	DomainId    string `json:"domainId"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}
