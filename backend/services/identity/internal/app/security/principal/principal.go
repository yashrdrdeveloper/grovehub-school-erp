package principal
import ("context";"errors";"time")
type key struct{}
type Principal struct{AccountID,PersonGID,TenantID,SessionID string;Admin bool;AuthenticatedAt time.Time}
func With(ctx context.Context,p Principal)context.Context{return context.WithValue(ctx,key{},p)}
func From(ctx context.Context)(Principal,bool){p,ok:=ctx.Value(key{}).(Principal);return p,ok&&p.AccountID!=""}
func Require(ctx context.Context)(Principal,error){p,ok:=From(ctx);if !ok{return p,errors.New("authentication required")};return p,nil}
func RequireAdmin(ctx context.Context)(Principal,error){p,err:=Require(ctx);if err!=nil{return p,err};if !p.Admin||p.TenantID==""{return p,errors.New("admin access forbidden")};return p,nil}
func Recent(ctx context.Context,maxAge time.Duration)error{p,err:=Require(ctx);if err!=nil{return err};if p.AuthenticatedAt.IsZero()||time.Since(p.AuthenticatedAt)>maxAge{return errors.New("recent authentication required")};return nil}
