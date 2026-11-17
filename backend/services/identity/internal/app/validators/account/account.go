package account
import ("grovehub/libs/core/search";"grovehub/services/identity/internal/app/dto/account")
type Validator interface{NormalizeCreate(account.Create)(account.Create,error);NormalizeUsername(string)(string,error);Search(search.Request)(search.Query[string],error)}
