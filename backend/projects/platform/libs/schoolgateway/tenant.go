package schoolgateway

import (
 "context"
 "errors"
)

type TenantCheck func(context.Context,string)error
type tenantCheckKey struct{}
var ErrTenantDenied=errors.New("tenant access denied")

func WithTenantCheck(ctx context.Context,check TenantCheck)context.Context{
 if check==nil{return ctx}
 return context.WithValue(ctx,tenantCheckKey{},check)
}

// CheckTenant is called only after deriving tenant identity from a trusted
// account or session. Nil policy preserves isolated module/test use.
func CheckTenant(ctx context.Context,tenant string)error{
 check,_:=ctx.Value(tenantCheckKey{}).(TenantCheck)
 if check==nil{return nil}
 if tenant==""{return ErrTenantDenied}
 return check(ctx,tenant)
}
