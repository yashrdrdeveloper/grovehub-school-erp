package invitation
import(
 "context"
 "crypto/aes"
 "crypto/cipher"
 "encoding/json"
 "errors"
 "time"
 "grovehub/services/identity/internal/app/config"
 "grovehub/services/identity/internal/app/repository/invitation"
 mailhelpers "grovehub/services/identity/internal/app/services/registration/operations/register/helpers"
 integrationclient "grovehub/services/integration/client/connection"
)
type Worker struct{repository invitation.Repository;client integrationclient.Client;config config.Config}
func NewWorker(r invitation.Repository,cl integrationclient.Client,c config.Config)*Worker{return &Worker{repository:r,client:cl,config:c}}
func(w *Worker)RunOnce(ctx context.Context)(bool,error){
 m,ok,err:=w.repository.Claim(ctx,time.Now().UTC());if err!=nil||!ok{return ok,err}
 if len(w.config.EncryptionKey)!=32{return true,errors.New("mail encryption key unavailable")}
 block,err:=aes.NewCipher(w.config.EncryptionKey);if err!=nil{return true,err}
 gcm,err:=cipher.NewGCM(block);if err!=nil{return true,err}
 plain,err:=gcm.Open(nil,m.Nonce,m.Ciphertext,[]byte(m.AccountID));if err!=nil{return true,err}
 var data Encrypted;if err=json.Unmarshal(plain,&data);err!=nil{return true,err}
 mail:=mailhelpers.Mail{Title:data.Title,Name:data.Name,Description:data.Description,Button:data.Button,URL:data.URL,Minutes:data.Minutes}
 err=mailhelpers.Send(ctx,w.client,w.config.MailTenantID,w.config.MailConnectionID,m.Email,mail)
 if err!=nil{
  delay:=time.Duration(1<<min(m.Attempts+1,9))*time.Second
  _=w.repository.Retry(ctx,m.ID,time.Now().UTC().Add(delay))
  return true,err
 }
 return true,w.repository.Delivered(ctx,m.ID)
}
