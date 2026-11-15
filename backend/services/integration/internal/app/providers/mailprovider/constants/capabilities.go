package constants

import dto "grovehub/services/integration/internal/app/providers/dto"

const (
    CapabilityMailSend dto.Capability = "mail.send"
    CapabilityMailRead dto.Capability = "mail.read"
    CapabilityMailAttachment dto.Capability = "mail.attachment.get"
    CapabilityMailHistory dto.Capability = "mail.history"
    CapabilityMailDelta dto.Capability = "mail.delta"
    CapabilityMailWatch dto.Capability = "mail.watch"
)
