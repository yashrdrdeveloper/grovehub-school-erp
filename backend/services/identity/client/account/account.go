package account
import(
 "context"
 coresearch "grovehub/libs/core/search"
 dto "grovehub/services/identity/internal/app/dto/account"
)
type(
 Create=dto.Create
 CreateResult=dto.CreateResult
 Account=dto.Account
 Update=dto.Update
)
// @Client identityAccountClient
type Client interface{
 // @Proxy accounts.Get
 Get(ctx context.Context,id string)(dto.Account,error)
 // @Proxy accounts.Search
 Search(ctx context.Context,query coresearch.Request)(coresearch.Result[dto.Account],error)
}
