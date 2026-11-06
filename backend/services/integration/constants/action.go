package constants

type Action string
type Capability string

const (
 ActionMailProfile Action = "mail.profile"
 ActionMailList Action = "mail.list"
 ActionMailGet Action = "mail.get"
 ActionMailAttachment Action = "mail.attachment.get"
 ActionMailHistory Action = "mail.history"
 ActionMailDelta Action = "mail.delta"
 ActionMailSend Action = "mail.send"

 CapabilityMailSend Capability = "mail.send"
 CapabilityMailRead Capability = "mail.read"
 CapabilityMailAttachment Capability = "mail.attachment.get"
 CapabilityMailHistory Capability = "mail.history"
 CapabilityMailDelta Capability = "mail.delta"
 CapabilityMailWatch Capability = "mail.watch"

 ProviderGmail = "GMAIL"
 ProviderOutlook = "OUTLOOK"
 ResourceGmailMailbox = "GMAIL_MAILBOX"
 StatusActive = "ACTIVE"
)
