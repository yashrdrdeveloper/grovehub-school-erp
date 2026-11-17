package impl
import(
 "errors"
 "strings"
 "unicode"
 "grovehub/libs/core/search"
 dto "grovehub/services/identity/internal/app/dto/account"
 contract "grovehub/services/identity/internal/app/validators/account"
)
type Validator struct{}
func New()contract.Validator{return Validator{}}
func (Validator)NormalizeCreate(r dto.Create)(dto.Create,error){
 r.Email=strings.ToLower(strings.TrimSpace(r.Email))
 r.PhoneNumber=strings.TrimSpace(r.PhoneNumber)
 r.Person.FirstName=strings.TrimSpace(r.Person.FirstName);r.Person.LastName=strings.TrimSpace(r.Person.LastName)
 if r.Person.MiddleName!=nil{v:=strings.TrimSpace(*r.Person.MiddleName);r.Person.MiddleName=&v}
 if !strings.Contains(r.Email,"@")||len(r.Email)>254||r.Person.FirstName==""||r.Person.LastName==""{return r,errors.New("invalid registration")}
 return r,nil
}
func (Validator)NormalizeUsername(s string)(string,error){
 s=strings.ToLower(strings.TrimSpace(s));if len(s)<3||len(s)>40{return "",errors.New("invalid username")}
 for _,c:=range s{if !(unicode.IsLower(c)||unicode.IsDigit(c)||c=='.'||c=='_'||c=='-'){return "",errors.New("invalid username")}}
 return s,nil
}
func (Validator)Search(q search.Request)(search.Query[string],error){
 bad:=errors.New("invalid account search")
 terms,err:=search.NormalizeSearchTerms(q.SearchText,256,10,bad);if err!=nil{return search.Query[string]{},err}
 rows,err:=search.ResolveRows(q.PageInfo.Rows,bad);if err!=nil{return search.Query[string]{},err}
 sortFields:=map[string]search.CursorValueValidator{"createdAt":func(v string)error{return search.ValidateTimeCursor(v,bad)},"updatedAt":func(v string)error{return search.ValidateTimeCursor(v,bad)},"accountId":func(v string)error{return search.ValidateUUIDv7Cursor(v,bad)},"status":func(v string)error{return search.ValidateNonEmptyCursor(v,bad)},"username":func(v string)error{return search.ValidateNonEmptyCursor(v,bad)}}
 sort,err:=search.NormalizeSort(q.Sort,search.QuerySort{Field:"createdAt",Direction:search.DirectionDescending},sortFields,bad);if err!=nil{return search.Query[string]{},err}
 fields:=map[string]search.FieldSpec{}
 _=fields
 // Validate untrusted filter expressions against the resource-owned allowlist.
 allowed:=map[string]bool{"accountId":true,"personGid":true,"username":true,"email":true,"status":true,"createdAt":true,"updatedAt":true}
 for _,f:=range q.Filters{if !allowed[f.Field]||len(f.Values)>50||len(f.Values)==0{return search.Query[string]{},bad}}
 filters:=make([]search.QueryFilter,0,len(q.Filters))
 for _,f:=range q.Filters{filters=append(filters,search.QueryFilter{Field:f.Field,Operator:f.Operator,Values:f.Values})}
 cursor,_,err:=search.NormalizeCursor(q.PageInfo.Cursor,sort,sortFields[sort.Field],search.ValidateUUIDv7Cursor,bad)
 if err!=nil{return search.Query[string]{},err}
 return search.Query[string]{SearchTerms:terms,Filters:filters,Sort:sort,Cursor:cursor,Limit:rows},nil
}
