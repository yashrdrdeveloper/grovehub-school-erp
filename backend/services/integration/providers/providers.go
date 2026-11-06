package providers

import (
 "context"
 "encoding/json"
 actionconstants "grovehub/services/integration/constants/actions"
 errorcode "grovehub/services/integration/constants/error"
)

// Public provider contracts; implementation and registration belong in impl.
type Action = actionconstants.Action
type Capability = actionconstants.Capability

const (
 ActionMailProfile = actionconstants.MailProfile
 ActionMailList = actionconstants.MailList
 ActionMailGet = actionconstants.MailGet
 ActionMailAttachment = actionconstants.MailAttachment
 ActionMailHistory = actionconstants.MailHistory
 ActionMailDelta = actionconstants.MailDelta
 ActionMailSend = actionconstants.MailSend
 CapabilityMailSend = actionconstants.CapabilityMailSend
 CapabilityMailRead = actionconstants.CapabilityMailRead
 CapabilityMailAttachment = actionconstants.CapabilityMailAttachment
 CapabilityMailHistory = actionconstants.CapabilityMailHistory
 CapabilityMailDelta = actionconstants.CapabilityMailDelta
 CapabilityMailWatch = actionconstants.CapabilityMailWatch
)

var (
 ErrUnsupportedAction = errorcode.ErrUnsupportedAction
 ErrInvalidActionInput = errorcode.ErrInvalidActionInput
 ErrProviderUnavailable = errorcode.ErrProviderUnavailable
 ErrUnknownProvider = errorcode.ErrUnknownProvider
 ErrDuplicateProvider = errorcode.ErrDuplicateProvider
)

type Descriptor struct {
 ID string `json:"id"`
 DisplayName string `json:"displayName"`
 Capabilities []Capability `json:"capabilities"`
 AuthType string `json:"authType"`
}
type Provider interface{Descriptor()Descriptor}
type Session struct {
 ConnectionID string
 ProviderID string
 AccessToken string
}
type ActionHandler interface{
 Execute(ctx context.Context,session Session,input json.RawMessage)(any,error)
}
type ActionRegistry interface {
 Register(provider string,action Action,handler ActionHandler)error
 Execute(ctx context.Context,session Session,action Action,input json.RawMessage)(any,error)
 Actions(provider string)[]Action
}
type Registry interface{
 Find(id string)(Provider,error)
 List()[]Descriptor
}
