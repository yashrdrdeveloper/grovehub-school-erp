package invitation
import(
 "context"
 "crypto/aes"
 "crypto/cipher"
 "crypto/rand"
 "encoding/json"
 "errors"
 "net/url"
 "time"
 "github.com/google/uuid"
 "grovehub/services/identity/internal/app/entity"
 "grovehub/services/identity/internal/app/config"
 "grovehub/services/identity/internal/app/repository/passwordtoken"
 "grovehub/services/identity/internal/app/repository/invitation"
 "grovehub/services/identity/internal/app/services/authentication/helpers"
)
type Service struct{Tokens passwordtoken.Repository;Mail invitation.Repository;Config config.Config}
type Encrypted struct{Title,Name,Description,Button,URL string;Minutes int}
func New(t passwordtoken.Repository,m invitation.Repository,c config.Config)*Service{return &Service{Tokens:t,Mail:m,Config:c}}
func(s *Service)Queue(ctx context.Context,account entity.Account,purpose entity.TokenPurpose,name string)error{
 if len(s.Config.EncryptionKey)!=32||s.Config.MailTenantID==""||s.Config.MailConnectionID==""{return errors.New("system mail unavailable")}
 baseURL:=s.Config.SetupURL
 if purpose==entity.PasswordReset{baseURL=s.Config.ResetURL}
 parsed,err:=url.Parse(baseURL);if err!=nil||parsed.Scheme!="https"||parsed.Host==""{return errors.New("invalid trusted reset/setup URL")}
 token,err:=helpers.NewSecret(32);if err!=nil{return err}
 hash,err:=helpers.Digest(token,32);if err!=nil{return err}
 minutes:=30;if purpose==entity.PasswordReset{minutes=15}
 if err:=s.Tokens.Invalidate(ctx,account.AccountID,purpose);err!=nil{return err}
 now:=time.Now().UTC()
 if err:=s.Tokens.Issue(ctx,entity.PasswordToken{TokenHash:hash,AccountID:account.AccountID,Purpose:purpose,ExpiresAt:now.Add(time.Duration(minutes)*time.Minute),CreatedAt:now});err!=nil{return err}
 // URL fragment keeps bearer proof out of the initial HTTP GET request path.
 parsed.Fragment="proof="+url.QueryEscape(token)
 mail:=Encrypted{Title:"Welcome to GrooveHub",Name:name,Description:"Your GrooveHub school account is ready. Choose your password to finish setup.",Button:"Set up your account",URL:parsed.String(),Minutes:minutes}
 if purpose==entity.PasswordReset{mail.Title="Reset your GrooveHub password";mail.Description="You requested to reset your GrooveHub password.";mail.Button="Reset password"}
 plain,err:=json.Marshal(mail);if err!=nil{return err}
 block,err:=aes.NewCipher(s.Config.EncryptionKey);if err!=nil{return err}
 gcm,err:=cipher.NewGCM(block);if err!=nil{return err}
 nonce:=make([]byte,gcm.NonceSize());if _,err=rand.Read(nonce);err!=nil{return err}
 id:=uuid.NewString()
 data:=gcm.Seal(nil,nonce,plain,[]byte(account.AccountID))
 return s.Mail.Queue(ctx,entity.Invitation{ID:id,AccountID:account.AccountID,Email:account.Email,Purpose:string(purpose),Nonce:nonce,Ciphertext:data,NextAttemptAt:now,CreatedAt:now})
}
