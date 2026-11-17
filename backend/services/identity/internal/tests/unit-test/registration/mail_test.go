package registration_test
import(
 "strings"
 "testing"
 mail "grovehub/services/identity/internal/app/services/registration/operations/register/helpers"
)
func TestRegistrationMailHTMLIsSafe(t *testing.T){
 html,err:=mail.RenderMail(mail.Mail{Title:"Welcome",Name:"<script>alert(1)</script>",Description:"Set up your password",Button:"Set password",URL:"https://app.example.test/setup#proof=abc",Minutes:30})
 if err!=nil{t.Fatal(err)}
 if strings.Contains(html,"<script>alert"){t.Fatal("HTML injection was not escaped")}
 if !strings.Contains(html,"Set password"){t.Fatal("missing button")}
}
