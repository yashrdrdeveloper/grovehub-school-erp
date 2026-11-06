package impl

import (
 "context"
 "encoding/json"
 "fmt"
 "sort"
 "strings"
 "grovehub/services/integration/providers"
)

type actionRegistry struct{handlers map[string]map[providers.Action]providers.ActionHandler}

func NewActionRegistry()providers.ActionRegistry{
 return &actionRegistry{handlers:make(map[string]map[providers.Action]providers.ActionHandler)}
}
func(r *actionRegistry)Register(provider string,action providers.Action,handler providers.ActionHandler)error{
 provider=strings.ToUpper(strings.TrimSpace(provider))
 if provider==""||strings.TrimSpace(string(action))==""||handler==nil{return providers.ErrInvalidActionInput}
 if r.handlers[provider]==nil{r.handlers[provider]=make(map[providers.Action]providers.ActionHandler)}
 if _,exists:=r.handlers[provider][action];exists{
  return fmt.Errorf("%w: %s %s",providers.ErrDuplicateProvider,provider,action)
 }
 r.handlers[provider][action]=handler
 return nil
}
func(r *actionRegistry)Execute(ctx context.Context,session providers.Session,action providers.Action,input json.RawMessage)(any,error){
 items:=r.handlers[strings.ToUpper(session.ProviderID)]
 if items==nil{return nil,providers.ErrProviderUnavailable}
 handler,ok:=items[action]
 if !ok{return nil,providers.ErrUnsupportedAction}
 return handler.Execute(ctx,session,input)
}
func(r *actionRegistry)Actions(provider string)[]providers.Action{
 handlers:=r.handlers[strings.ToUpper(provider)]
 result:=make([]providers.Action,0,len(handlers))
 for action:=range handlers{result=append(result,action)}
 sort.Slice(result,func(i,j int)bool{return result[i]<result[j]})
 return result
}
