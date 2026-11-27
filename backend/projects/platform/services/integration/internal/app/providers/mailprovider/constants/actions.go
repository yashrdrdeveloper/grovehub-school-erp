package constants

import dto "grovehub/services/integration/internal/app/providers/dto"

const (
    MailProfile dto.Action = "mail.profile"
    MailList dto.Action = "mail.list"
    MailGet dto.Action = "mail.get"
    MailAttachment dto.Action = "mail.attachment.get"
    MailHistory dto.Action = "mail.history"
    MailDelta dto.Action = "mail.delta"
    MailSend dto.Action = "mail.send"
)
