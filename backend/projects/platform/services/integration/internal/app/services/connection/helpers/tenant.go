package helpers

import (
 "context"
 schoolgateway "grovehub/launcher/session"
)

type tenantStateKey struct{}

// WithTenantState installs the configured lookup, not caller-provided headers.
func WithTenantState(ctx context.Context, lookup schoolgateway.TenantStateLookup) context.Context {
 if lookup==nil{return ctx}
 ctx=schoolgateway.WithTenantCheck(ctx,schoolgateway.ActiveTenantCheck(lookup))
 return context.WithValue(ctx,tenantStateKey{},lookup)
}

func TenantActive(ctx context.Context, tenant string)(bool,error){
 if lookup,ok:=ctx.Value(tenantStateKey{}).(schoolgateway.TenantStateLookup);ok{
  if tenant==""{return false,schoolgateway.ErrTenantDenied}
  state,err:=lookup(ctx,tenant)
  if err!=nil{return false,err}
  if state.Generation<1{return false,schoolgateway.ErrTenantDenied}
  return state.Active,nil
 }
 _,err:=schoolgateway.TenantGeneration(ctx,tenant)
 return err==nil,err
}

// RequireTenantGeneration enforces both live registry state and any verified
// browser principal present. Workload-only calls have no browser principal.
func RequireTenantGeneration(ctx context.Context, tenant string)(int64,error){
 generation,err:=schoolgateway.TenantGeneration(ctx,tenant)
 if err!=nil{return 0,ErrConnectionInactive}
 if actor,ok:=schoolgateway.PrincipalFrom(ctx);ok{
  if actor.TenantID!=tenant||(generation!=0&&actor.SessionGeneration!=generation){return 0,ErrConnectionInactive}
 }
 return generation,nil
}

func RequireSameGeneration(ctx context.Context,tenant string,issued int64)error{
 generation,err:=RequireTenantGeneration(ctx,tenant)
 if err!=nil{return err}
 if generation!=issued{return ErrConnectionInactive}
 return nil
}
