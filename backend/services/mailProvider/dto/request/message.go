package request

type SendMessage struct {
	GID      string   `json:"gid"`
	SenderId string   `json:"senderId"`
	To       []string `json:"to"`
	Cc       []string `json:"cc,omitempty"`
	Bcc      []string `json:"bcc,omitempty"`
	Subject  string   `json:"subject"`
	TextBody string   `json:"textBody,omitempty"`
	HTMLBody string   `json:"htmlBody,omitempty"`
}

type ReplyMessage struct {
	TextBody string `json:"textBody,omitempty"`
	HTMLBody string `json:"htmlBody,omitempty"`
}
