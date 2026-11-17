package impl
import ("encoding/base64";"os";"strings";"errors";"net/url";"github.com/google/uuid";"grovehub/services/identity/internal/app/config")
func Load() (config.Config,error) {
 c:=config.Config{AdminIDs:map[string]bool{},SetupURL:os.Getenv("IDENTITY_SETUP_URL"),ResetURL:os.Getenv("IDENTITY_RESET_URL"),MailTenantID:os.Getenv("IDENTITY_SYSTEM_MAIL_TENANT_ID"),MailConnectionID:os.Getenv("IDENTITY_SYSTEM_MAIL_CONNECTION_ID"),AllowedOrigin:os.Getenv("IDENTITY_ALLOWED_ORIGIN")}
 for _,s:=range strings.Split(os.Getenv("IDENTITY_ADMIN_ACCOUNT_IDS"),","){if id,err:=uuid.Parse(strings.TrimSpace(s));err==nil{c.AdminIDs[id.String()]=true}}
 key,err:=base64.StdEncoding.DecodeString(os.Getenv("IDENTITY_INVITATION_KEY_B64"));if err==nil&&len(key)==32{c.EncryptionKey=key}
 if c.AllowedOrigin!=""{u,e:=url.Parse(c.AllowedOrigin);if e!=nil||u.Scheme!="https"||u.Host==""||u.Path!=""{return c,errors.New("identity allowed origin must be an HTTPS origin")}}
 return c,nil
}
