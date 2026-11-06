package impl

import (
 "context"
 "crypto/rand"
 "crypto/sha256"
 "encoding/base64"
 "encoding/hex"
 "errors"
 "time"

 integrationresponse "grovehub/services/integration/dto/response"
 "grovehub/services/integration/entity"
 "grovehub/services/integration/providers/gmail"
 integrationrepository "grovehub/services/integration/repository"
)
var (
 ErrOAuthUnavailable=errors.New("gmail OAuth is not configured")
 ErrOAuthInvalid=errors.New("invalid OAuth callback")
 ErrConnectionInactive=errors.New("connection not active")
)

func(s *ConnectionService)BeginOAuth(ctx context.Context,tenantID,id string)(integrationresponse.OAuthStart,error){
 if s.gmail==nil{return integrationresponse.OAuthStart{},ErrOAuthUnavailable}
 c,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.OAuthStart{},err}
 if c.Provider!="GMAIL"||c.Status!=entity.StatusPendingSetup{return integrationresponse.OAuthStart{},ErrInvalidConnection}
 nonce:=make([]byte,32)
 if _,err=rand.Read(nonce);err!=nil{return integrationresponse.OAuthStart{},err}
 state:=base64.RawURLEncoding.EncodeToString(nonce)
 hash:=hashState(state)
 err=s.repository.SaveOAuthState(ctx,entity.OAuthState{Hash:hash,ConnectionID:id,ExpiresAt:time.Now().UTC().Add(10*time.Minute)})
 if err!=nil{return integrationresponse.OAuthStart{},err}
 return integrationresponse.OAuthStart{AuthorizationURL:s.gmail.AuthorizationURL(state)},nil
}
func(s *ConnectionService)FinishOAuth(ctx context.Context,code,state string)(integrationresponse.OAuthCompletion,error){
 if s.gmail==nil{return integrationresponse.OAuthCompletion{},ErrOAuthUnavailable}
 if code==""||len(state)<40{return integrationresponse.OAuthCompletion{},ErrOAuthInvalid}
 oauthState,err:=s.repository.ConsumeOAuthState(ctx,hashState(state),time.Now().UTC())
 if err!=nil {
  if errors.Is(err,integrationrepository.ErrOAuthStateInvalid){return integrationresponse.OAuthCompletion{},ErrOAuthInvalid}
  return integrationresponse.OAuthCompletion{},err
 }
 token,err:=s.gmail.Exchange(ctx,code)
 if err!=nil{return integrationresponse.OAuthCompletion{},err}
 if token.RefreshToken==""{return integrationresponse.OAuthCompletion{},ErrOAuthInvalid}
 mailbox,err:=s.gmail.Mailbox(ctx,token)
 if err!=nil{return integrationresponse.OAuthCompletion{},err}
 encrypted,err:=s.cipher.Seal(token)
 if err!=nil{return integrationresponse.OAuthCompletion{},err}
 if err=s.repository.Activate(ctx,oauthState.ConnectionID,mailbox,encrypted);err!=nil{return integrationresponse.OAuthCompletion{},err}
 return integrationresponse.OAuthCompletion{Status:entity.StatusActive},nil
}
func hashState(raw string)string{sum:=sha256.Sum256([]byte(raw));return hex.EncodeToString(sum[:])}
func(s *ConnectionService)validToken(ctx context.Context,connectionID string)(gmail.Token,error){
 data,err:=s.repository.Credential(ctx,connectionID)
 if err!=nil{return gmail.Token{},err}
 token,err:=s.cipher.Open(data)
 if err!=nil{return gmail.Token{},err}
 if time.Until(token.Expiry)>time.Minute{return token,nil}
 token,err=s.gmail.Refresh(ctx,token)
 if err!=nil{return gmail.Token{},err}
 encrypted,err:=s.cipher.Seal(token)
 if err!=nil{return gmail.Token{},err}
 if err=s.repository.UpdateCredential(ctx,connectionID,encrypted);err!=nil{return gmail.Token{},err}
 return token,nil
}
