package actions

// Mail actions are provider-neutral identifiers exposed by mail providers.
type Action string
type Capability string

const (
	MailProfile    Action = "mail.profile"
	MailList       Action = "mail.list"
	MailGet        Action = "mail.get"
	MailAttachment Action = "mail.attachment.get"
	MailHistory    Action = "mail.history"
	MailDelta      Action = "mail.delta"
	MailSend       Action = "mail.send"
)

const (
	CapabilityMailSend       Capability = "mail.send"
	CapabilityMailRead       Capability = "mail.read"
	CapabilityMailAttachment Capability = "mail.attachment.get"
	CapabilityMailHistory    Capability = "mail.history"
	CapabilityMailDelta      Capability = "mail.delta"
	CapabilityMailWatch      Capability = "mail.watch"
)
