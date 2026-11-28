// Package personprofile implements a read-only membership projection. Identity
// owns the account link; Foundation remains the only owner of person records.
package personprofile

import (
 "context"
 "github.com/google/uuid"
 foundation "grovehub/services/foundation/client/person"
 dto "grovehub/services/identity/internal/app/dto/account"
 repo "grovehub/services/identity/internal/app/repository/account"
 "grovehub/services/identity/internal/app/services/authorization"
 schoolgateway "grovehub/libs/platform/session"
)

const route = "GET /identity/protected/accounts/{accountId}/person"

type AccountReader interface {
 Get(context.Context,string)(dto.Account,error)
}
type PersonReader interface {
 GetPersonByGID(context.Context,foundation.GetPersonByGIDRequest)(foundation.Person,error)
}
type Policy interface {
 Authorize(context.Context,string)error
}

type Operation struct {
 accounts AccountReader
 persons PersonReader
 policy Policy
}

func New(accounts AccountReader,persons PersonReader,policy Policy)*Operation {
 return &Operation{accounts:accounts,persons:persons,policy:policy}
}

func(o *Operation) Execute(ctx context.Context,id string)(dto.PersonProfile,error){
 if o==nil||o.accounts==nil||o.persons==nil||o.policy==nil{return dto.PersonProfile{},authorization.ErrUnavailable}
 if err:=o.policy.Authorize(ctx,route);err!=nil{return dto.PersonProfile{},err}
 actor,ok:=schoolgateway.PrincipalFrom(ctx)
 if !ok||actor.TenantID==""{return dto.PersonProfile{},authorization.ErrForbidden}
 parsed,err:=uuid.Parse(id)
 if err!=nil||parsed==uuid.Nil||parsed.String()!=id{return dto.PersonProfile{},repo.ErrNotFound}
 // The production reader is the guarded account service: it derives the
 // tenant from the verified principal and denies a foreign account before
 // invoking Foundation. The caller never supplies a global person ID.
 account,err:=o.accounts.Get(ctx,id)
 if err!=nil{return dto.PersonProfile{},err}
 if account.AccountID!=id{return dto.PersonProfile{},repo.ErrNotFound}
 gid,err:=uuid.Parse(account.PersonGID)
 if err!=nil||gid==uuid.Nil||gid.String()!=account.PersonGID{return dto.PersonProfile{},repo.ErrNotFound}
 result,err:=o.persons.GetPersonByGID(ctx,foundation.GetPersonByGIDRequest{Params:foundation.PersonByGID{GID:account.PersonGID}})
 if err!=nil{return dto.PersonProfile{},authorization.ErrUnavailable}
 if result.GID!=account.PersonGID{return dto.PersonProfile{},repo.ErrNotFound}
 // Recheck configured authorization after the cross-module call, and expose
 // only the fields required for displaying the account holder's name.
 if err:=o.policy.Authorize(ctx,route);err!=nil{return dto.PersonProfile{},err}
 return dto.PersonProfile{PersonGID:result.GID,FirstName:result.FirstName,MiddleName:result.MiddleName,LastName:result.LastName},nil
}
