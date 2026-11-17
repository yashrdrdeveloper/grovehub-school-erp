package credential_test
import(
 "testing"
 "grovehub/services/identity/internal/app/services/credential/helpers"
)
func TestArgon2HashAndVerify(t *testing.T){
 hash,err:=helpers.HashPassword("correct horse battery staple")
 if err!=nil{t.Fatal(err)}
 if hash=="correct horse battery staple"{t.Fatal("plaintext credential was retained")}
 ok,err:=helpers.VerifyPassword("correct horse battery staple",hash)
 if err!=nil||!ok{t.Fatalf("valid credential rejected: %v",err)}
 ok,err=helpers.VerifyPassword("wrong-password",hash)
 if err!=nil||ok{t.Fatal("invalid password was accepted")}
}
