package credential
import "context"
type Service interface {
 Set(context.Context,string,string)error
 Verify(context.Context,string,string)(bool,error)
 Replace(context.Context,string,string)error
}
