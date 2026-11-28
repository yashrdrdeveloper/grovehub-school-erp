package session_test

import (
 "context"
 "errors"
 "testing"
 . "grovehub/launcher/services/gateway/session"
)

func TestTenantSessionGeneration(t *testing.T) {
 generation := int64(1)
 active := true
 ctx := WithTenantCheck(context.Background(), func(_ context.Context, tenant string) (int64, error) {
  if tenant != "school" || !active { return 0, ErrTenantDenied }
  return generation, nil
 })
 issued, err := TenantGeneration(ctx, "school")
 if err != nil || issued != 1 { t.Fatal(issued, err) }
 if err := CheckTenantSession(ctx, "school", issued); err != nil { t.Fatal(err) }
 active = false
 generation++
 if err := CheckTenantSession(ctx, "school", issued); !errors.Is(err, ErrTenantDenied) { t.Fatal("disabled school accepted", err) }
 active = true
 for _, old := range []int64{0, issued, generation+1} {
  if err := CheckTenantSession(ctx, "school", old); !errors.Is(err, ErrTenantDenied) { t.Fatal("incorrect generation accepted", old, err) }
 }
 if err := CheckTenantSession(ctx, "school", generation); err != nil { t.Fatal("new session denied", err) }
 if err := CheckTenantSession(ctx, "other", generation); !errors.Is(err, ErrTenantDenied) { t.Fatal("wrong tenant accepted", err) }
}

func TestCheckAccountTenant(t *testing.T) {
 tenant := "stored-school"
 ctx := WithTenantCheck(context.Background(), func(_ context.Context, got string) (int64,error) {
  if got != tenant { t.Fatal("incorrect tenant",got) }
  return 0, ErrTenantDenied
 })
 if !errors.Is(CheckAccountTenant(ctx, &tenant), ErrTenantDenied) { t.Fatal("disabled account accepted") }
 if !errors.Is(CheckAccountTenant(ctx, nil), ErrTenantDenied) { t.Fatal("missing tenant accepted") }
}

func TestPrincipalContext(t *testing.T) {
 ctx := context.Background()
 if _, ok := PrincipalFrom(ctx); ok { t.Fatal("anonymous context accepted") }
 if _, ok := PrincipalFrom(WithPrincipal(ctx, Principal{AccountID:"a"})); ok { t.Fatal("missing session accepted") }
 p := Principal{AccountID:"a", SessionID:"s", TenantID:"school", SessionGeneration:1}
 if got, ok := PrincipalFrom(WithPrincipal(ctx,p)); !ok || got != p { t.Fatal(got,ok) }
}

func TestTenantLookupAndMissingPolicy(t *testing.T) {
 ctx := context.Background()
 if ActiveTenantCheck(nil)!=nil || WithTenantCheck(ctx,nil)!=ctx { t.Fatal("nil policy changed") }
 if n,err:=TenantGeneration(ctx,"");n!=0||err!=nil { t.Fatal(n,err) }
 if err:=CheckTenantSession(ctx,"school",3);err!=nil { t.Fatal(err) }
 check:=ActiveTenantCheck(func(_ context.Context,id string)(TenantState,error){
  switch id {
  case "active": return TenantState{Active:true,Generation:3},nil
  case "disabled": return TenantState{Generation:3},nil
  case "zero": return TenantState{Active:true},nil
  default: return TenantState{},errors.New("unavailable")
  }
 })
 scoped:=WithTenantCheck(ctx,check)
 for _,id:=range []string{"","disabled","zero","missing"} {
  if !errors.Is(CheckTenant(scoped,id),ErrTenantDenied) { t.Fatal(id) }
 }
 id:="active"
 if err:=CheckAccountTenant(scoped,&id);err!=nil {t.Fatal(err)}
 if err:=CheckAccountGeneration(scoped,&id,3);err!=nil {t.Fatal(err)}
 if !errors.Is(CheckAccountGeneration(scoped,nil,3),ErrTenantDenied) {t.Fatal("missing tenant accepted")}
}
