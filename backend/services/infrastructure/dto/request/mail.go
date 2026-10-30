package request

type MailAddress struct {
	Email       string
	DisplayName string
}

type MailMessage struct {
	From                MailAddress
	To                  []string
	Cc                  []string
	Bcc                 []string
	Subject             string
	TextBody            string
	HTMLBody            string
	InternetMessageID   string
	InReplyTo           string
	ReferenceMessageIDs []string
}
