package impl
import(
 "errors"
 "strings"
 "github.com/google/uuid"
 corenormalize "grovehub/libs/core/normalize"
 "grovehub/libs/core/search"
 dto "grovehub/services/identity/internal/app/dto/account"
 contract "grovehub/services/identity/internal/app/validators/account"
)
type Validator struct{}
func New()contract.Validator{return Validator{}}
func (Validator)NormalizeCreate(r dto.Create)(dto.Create,error){
 email,err:=corenormalize.Email(r.Email);if err!=nil{return r,errors.New("invalid email")};r.Email=email
 if strings.TrimSpace(r.PhoneNumber)!=""{phone,e:=corenormalize.Phone(r.PhoneNumber);if e!=nil{return r,errors.New("invalid phone")};r.PhoneNumber=phone}
 r.Person.FirstName=strings.TrimSpace(r.Person.FirstName);r.Person.LastName=strings.TrimSpace(r.Person.LastName)
 if r.Person.MiddleName!=nil{v:=strings.TrimSpace(*r.Person.MiddleName);r.Person.MiddleName=&v}
 if r.Person.FirstName==""||r.Person.LastName==""||len(r.Person.FirstName)>100||len(r.Person.LastName)>100{return r,errors.New("invalid registration")}
 return r,nil
}
func (Validator)NormalizeUsername(s string)(string,error){
 s=strings.ToLower(strings.TrimSpace(s));if len(s)<3||len(s)>40{return "",errors.New("invalid username")}
 for _,c:=range s{if !((c>='a'&&c<='z')||(c>='0'&&c<='9')||c=='.'||c=='_'||c=='-'){return "",errors.New("invalid username")}}
 return s,nil
}
func (Validator)Search(q search.Request)(search.Query[string],error){
 bad:=errors.New("invalid account search")
 terms,err:=search.NormalizeSearchTerms(q.SearchText,256,10,bad);if err!=nil{return search.Query[string]{},err}
 rows,err:=search.ResolveRows(q.PageInfo.Rows,bad);if err!=nil{return search.Query[string]{},err}
 sortFields:=map[string]search.CursorValueValidator{"createdAt":func(v string)error{return search.ValidateTimeCursor(v,bad)},"updatedAt":func(v string)error{return search.ValidateTimeCursor(v,bad)},"accountId":func(v string)error{return search.ValidateUUIDv7Cursor(v,bad)},"status":func(v string)error{return search.ValidateNonEmptyCursor(v,bad)},"username":func(v string)error{return search.ValidateNonEmptyCursor(v,bad)}}
 sort,err:=search.NormalizeSort(q.Sort,search.QuerySort{Field:"createdAt",Direction:search.DirectionDescending},sortFields,bad);if err!=nil{return search.Query[string]{},err}
 fields:=map[string]search.FieldSpec{
  "accountId":{Operators:search.EqualityOperators(),Normalizer:func(_ search.Operator,v string)(string,error){if _,err:=uuid.Parse(v);err!=nil{return "",bad};return v,nil}},
  "personGid":{Operators:search.EqualityOperators(),Normalizer:func(_ search.Operator,v string)(string,error){if _,err:=uuid.Parse(v);err!=nil{return "",bad};return v,nil}},
  "username":{Operators:search.TextOperators(),Normalizer:func(_ search.Operator,v string)(string,error){v=strings.ToLower(strings.TrimSpace(v));if v==""{return "",bad};return v,nil}},
  "email":{Operators:search.TextOperators(),Normalizer:func(_ search.Operator,v string)(string,error){v=strings.ToLower(strings.TrimSpace(v));if v==""{return "",bad};return v,nil}},
  "status":{Operators:search.EqualityOperators(),Normalizer:func(_ search.Operator,v string)(string,error){v=strings.ToUpper(strings.TrimSpace(v));switch v{case "ACTIVE","PENDING_SETUP","SUSPENDED","DEACTIVATED":return v,nil};return "",bad}},
  "createdAt":{Operators:search.TimeOperators(),Normalizer:func(_ search.Operator,v string)(string,error){v,e:=corenormalize.RFC3339UTC(v);if e!=nil{return "",bad};return v,nil}},
  "updatedAt":{Operators:search.TimeOperators(),Normalizer:func(_ search.Operator,v string)(string,error){v,e:=corenormalize.RFC3339UTC(v);if e!=nil{return "",bad};return v,nil}},
 }
 filters,err:=search.NormalizeFilters(q.Filters,fields,20,50,bad)
 if err!=nil{return search.Query[string]{},err}
 cursor,_,err:=search.NormalizeCursor(q.PageInfo.Cursor,sort,sortFields[sort.Field],search.ValidateUUIDv7Cursor,bad)
 if err!=nil{return search.Query[string]{},err}
 return search.Query[string]{SearchTerms:terms,Filters:filters,Sort:sort,Cursor:cursor,Limit:rows},nil
}
