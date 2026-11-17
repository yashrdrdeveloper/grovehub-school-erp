package validator_test
import(
 "testing"
 "grovehub/libs/core/search"
 validator "grovehub/services/identity/internal/app/validators/account/impl"
)
func TestAccountUsernamePolicy(t *testing.T){
 v:=validator.New()
 normalized,err:=v.NormalizeUsername(" John.Doe ");if err!=nil||normalized!="john.doe"{t.Fatal("username normalization failed",err)}
 if _,err:=v.NormalizeUsername("bad username");err==nil{t.Fatal("space should be rejected")}
}
func TestAccountSearchRejectsUntrustedFields(t *testing.T){
 v:=validator.New()
 _,err:=v.Search(search.Request{Filters:[]search.Filter{{Field:"passwordHash",Operator:search.OperatorEquals,Values:[]string{"bad"}}}})
 if err==nil{t.Fatal("passwordHash filter accepted")}
 _,err=v.Search(search.Request{Filters:[]search.Filter{{Field:"status",Operator:search.OperatorEquals,Values:[]string{"ACTIVE OR 1=1"}}}})
 if err==nil{t.Fatal("untrusted status accepted")}
}
