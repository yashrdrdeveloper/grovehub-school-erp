package search
import(
 "context"
 "time"
 coresearch "grovehub/libs/core/search"
 dto "grovehub/services/identity/internal/app/dto/account"
 mapper "grovehub/services/identity/internal/app/mappers/account"
 accountrepo "grovehub/services/identity/internal/app/repository/account"
 "grovehub/services/identity/internal/app/security/principal"
 validator "grovehub/services/identity/internal/app/validators/account"
)
type Operation struct{accounts accountrepo.Repository;mapper mapper.Mapper;validator validator.Validator}
func New(a accountrepo.Repository,m mapper.Mapper,v validator.Validator)*Operation{return &Operation{accounts:a,mapper:m,validator:v}}
func(o *Operation)Execute(ctx context.Context,request coresearch.Request)(coresearch.Result[dto.Account],error){
 p,err:=principal.RequireAdmin(ctx);if err!=nil{return coresearch.Result[dto.Account]{},err}
 query,err:=o.validator.Search(request);if err!=nil{return coresearch.Result[dto.Account]{},err}
 result,err:=o.accounts.Search(ctx,p.TenantID,query);if err!=nil{return coresearch.Result[dto.Account]{},err}
 out:=coresearch.Result[dto.Account]{Items:make([]dto.Account,0,len(result.Items))}
 for _,a:=range result.Items{out.Items=append(out.Items,o.mapper.ToResponse(a))}
 out.PageInfo=coresearch.PageResponse{Rows:len(out.Items),HasNext:result.HasNext}
 if result.HasNext&&len(result.Items)>0{
  last:=result.Items[len(result.Items)-1];value:=last.CreatedAt.UTC().Format(time.RFC3339Nano)
  switch query.Sort.Field{
  case "accountId":value=last.AccountID
  case "updatedAt":value=last.UpdatedAt.UTC().Format(time.RFC3339Nano)
  case "status":value=string(last.Status)
  case "username":if last.Username!=nil{value=*last.Username}else{value=""}
  }
  out.PageInfo.NextCursor,err=coresearch.EncodeCursor(coresearch.CursorPayload{Field:query.Sort.Field,Direction:query.Sort.Direction,Value:value,CursorId:last.AccountID})
 }
 return out,err
}
