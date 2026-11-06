package providers

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "strings"
 "sort"
)

// Action is an exact, stable capability identifier.
// A provider only registers actions that it can actually execute.
type Action string

const (
 ActionMailProfile Action = "mail.profile"
 ActionMailList Action = "mail.list"
 ActionMailGet Action = "mail.get"
 ActionMailAttachment Action = "mail.attachment.get"
 ActionMailHistory Action = "mail.history"
 ActionMailDelta Action = "mail.delta"
 ActionMailSend Action = "mail.send"
)

var (
 ErrUnsupportedAction = errors.New("integration action is unsupported")
 ErrInvalidActionInput = errors.New("invalid integration action input")
 ErrProviderUnavailable = errors.New("integration provider is not configured")
)

// Session is produced by Integration after tenant ownership and connection status
// checks. Credentials never cross the public action API.
type Session struct {
 ConnectionID string
 ProviderID string
 AccessToken string
}

// ActionHandler implements one optional provider capability.
type ActionHandler interface {
 Execute(ctx context.Context, session Session, input json.RawMessage) (any, error)
}

// ActionRegistry is a dependency-injected strategy registry.
// No provider-specific selection happens inside the generic HTTP controller.
type ActionRegistry struct {
 handlers map[string]map[Action]ActionHandler
}

func NewActionRegistry() *ActionRegistry {
 return &ActionRegistry{handlers:make(map[string]map[Action]ActionHandler)}
}

func (r *ActionRegistry) Register(provider string, action Action, handler ActionHandler) error {
 provider = strings.ToUpper(strings.TrimSpace(provider))
 if provider=="" || strings.TrimSpace(string(action))=="" || handler==nil {
  return ErrInvalidActionInput
 }
 if r.handlers[provider]==nil {r.handlers[provider]=make(map[Action]ActionHandler)}
 if _,exists:=r.handlers[provider][action];exists {
  return fmt.Errorf("duplicate action registration: %s %s",provider,action)
 }
 r.handlers[provider][action]=handler
 return nil
}

func (r *ActionRegistry) Execute(ctx context.Context,session Session,action Action,input json.RawMessage)(any,error){
 items:=r.handlers[strings.ToUpper(session.ProviderID)]
 if items==nil{return nil,ErrProviderUnavailable}
 handler,ok:=items[action]
 if !ok{return nil,ErrUnsupportedAction}
 return handler.Execute(ctx,session,input)
}

func (r *ActionRegistry) Actions(provider string) []Action {
 handlers:=r.handlers[strings.ToUpper(provider)]
 result:=make([]Action,0,len(handlers))
 for action:=range handlers{result=append(result,action)}
 sort.Slice(result,func(i,j int)bool{return result[i]<result[j]})
 return result
}
