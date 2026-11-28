package personprofile

import (
 "context"
 "encoding/json"
 "errors"
 "strings"
 "testing"
 schoolgateway "grovehub/libs/platform/session"
 foundation "grovehub/services/foundation/client/person"
 dto "grovehub/services/identity/internal/app/dto/account"
 repo "grovehub/services/identity/internal/app/repository/account"
 "grovehub/services/identity/internal/app/services/authorization"
)

const accountID="e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
const personID="019abcde-1234-7890-8abc-123456789abc"

type scopedAccounts struct{calls int;foreign bool}
func(a *scopedAccounts)Get(ctx context.Context,id string)(dto.Account,error){
 a.calls++
 actor,ok:=schoolgateway.PrincipalFrom(ctx)
 if !ok||actor.TenantID!="school"||a.foreign{return dto.Account{},repo.ErrNotFound}
 return dto.Account{AccountID:id,PersonGID:personID},nil
}
type globalPersons struct{calls int;wrong bool;after func()}
func(p *globalPersons)GetPersonByGID(_ context.Context,request foundation.GetPersonByGIDRequest)(foundation.Person,error){
 p.calls++
 if p.after!=nil{p.after()}
 gid:=request.Params.GID
 if p.wrong{gid=accountID}
 return foundation.Person{GID:gid,FirstName:"A",LastName:"B",Contacts:foundation.ResponseContacts{Emails:[]foundation.ResponseContactValue{{Value:"private@example.test"}}}},nil
}
type policyFunc func(context.Context,string)error
func(f policyFunc)Authorize(ctx context.Context,pattern string)error{return f(ctx,pattern)}

func TestMembershipProjection(t *testing.T){
 for _,tc:=range []struct{name string;foreign,disabled,wrong,disableDuringRead bool;want error;accounts,persons int}{
  {name:"authorized",accounts:1,persons:1},
  {name:"foreign",foreign:true,want:repo.ErrNotFound,accounts:1},
  {name:"disabled",disabled:true,want:authorization.ErrForbidden},
  {name:"wrong person response",wrong:true,want:repo.ErrNotFound,accounts:1,persons:1},
  {name:"disabled during lookup",disableDuringRead:true,want:authorization.ErrForbidden,accounts:1,persons:1},
 }{t.Run(tc.name,func(t *testing.T){
  accounts:=&scopedAccounts{foreign:tc.foreign}
  persons:=&globalPersons{wrong:tc.wrong}
  active:=!tc.disabled
  if tc.disableDuringRead{persons.after=func(){active=false}}
  policy:=policyFunc(func(_ context.Context,pattern string)error{if pattern!=route||!active{return authorization.ErrForbidden};return nil})
  operation:=New(accounts,persons,policy)
  ctx:=schoolgateway.WithPrincipal(context.Background(),schoolgateway.Principal{AccountID:"actor",SessionID:"session",TenantID:"school"})
  got,err:=operation.Execute(ctx,accountID)
  if !errors.Is(err,tc.want){t.Fatal("unexpected error",err)}
  if accounts.calls!=tc.accounts||persons.calls!=tc.persons{t.Fatal("unexpected calls",accounts.calls,persons.calls)}
  if tc.want==nil{raw,err:=json.Marshal(got);if err!=nil||got.PersonGID!=personID||got.FirstName!="A"||strings.Contains(string(raw),"private")||strings.Contains(string(raw),"contacts"){t.Fatal("unexpected profile",string(raw),err)}}
 })}
}
