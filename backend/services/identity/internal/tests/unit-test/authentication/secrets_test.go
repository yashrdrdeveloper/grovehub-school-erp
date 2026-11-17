package authentication_test
import(
 "testing"
 "grovehub/services/identity/internal/app/services/authentication/helpers"
)
func TestTokenSecretsAreRandomAndHashed(t *testing.T){
 a,err:=helpers.NewSecret(32);if err!=nil{t.Fatal(err)}
 b,err:=helpers.NewSecret(32);if err!=nil{t.Fatal(err)}
 if a==b{t.Fatal("random tokens collide")}
 hash,err:=helpers.Digest(a,32);if err!=nil||len(hash)!=64||hash==a{t.Fatal("invalid digest")}
 if _,err:=helpers.Digest(a,48);err==nil{t.Fatal("incorrect token length accepted")}
}
