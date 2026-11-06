package providers

import (
 "context"
 "encoding/json"
 "grovehub/services/integration/constants"
 errorcode "grovehub/services/integration/constants/error"
)

// Public provider contracts; implementation and registration belong in impl.
type Action = constants.Action
type Capability = constants.Capability

const (
 ActionMailProfile = constants.ActionMailProfile
 ActionMailList = constants.ActionMailList
 ActionMailGet = constants.ActionMailGet
 ActionMailAttachment = constants.ActionMailAttachment
 ActionMailHistory = constants.ActionMailHistory
 ActionMailDelta = constants.ActionMailDelta
 ActionMailSend = constants.ActionMailSend
 CapabilityMailSend = constants.CapabilityMailSend
 CapabilityMailRead = constants.CapabilityMailRead
 CapabilityMailAttachment = constants.CapabilityMailAttachment
 CapabilityMailHistory = constants.CapabilityMailHistory
 CapabilityMailDelta = constants.CapabilityMailDelta
 CapabilityMailWatch = constants.CapabilityMailWatch
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
