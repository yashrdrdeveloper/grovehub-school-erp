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
 "grovehub/services/integration/providers"
 integrationrepository "grovehub/services/integration/repository"
)

var (
 ErrOAuthUnavailable=errors.New("provider OAuth is not configured")
 ErrOAuthInvalid=errors.New("invalid OAuth callback")
 ErrConnectionInactive=errors.New("connection not active")
)

func(s *ConnectionService)BeginOAuth(ctx context.Context,tenantID,id string)(integrationresponse.OAuthStart,error){
 c,err:=s.repository.Get(ctx,tenantID,id)
 if err!=nil{return integrationresponse.OAuthStart{},err}
 if c.Status!=entity.StatusPendingSetup{return integrationresponse.OAuthStart{},ErrInvalidConnection}
 switch c.Provider{
 case "GMAIL":
  if s.gmail==nil{return integrationresponse.OAuthStart{},ErrOAuthUnavailable}
 case "OUTLOOK":
  if s.outlook==nil{return integrationresponse.OAuthStart{},ErrOAuthUnavailable}
 default:return integrationresponse.OAuthStart{},ErrOAuthUnavailable
 }
 nonce:=make([]byte,32)
 if _,err=rand.Read(nonce);err!=nil{return integrationresponse.OAuthStart{},err}
 state:=base64.RawURLEncoding.EncodeToString(nonce)
 hash:=hashState(state)
 err=s.repository.SaveOAuthState(ctx,entity.OAuthState{
  Hash:hash,ConnectionID:id,ExpiresAt:time.Now().UTC().Add(10*time.Minute),
 })
 if err!=nil{return integrationresponse.OAuthStart{},err}
 if c.Provider=="GMAIL"{
  return integrationresponse.OAuthStart{AuthorizationURL:s.gmail.AuthorizationURL(state)},nil
 }
 return integrationresponse.OAuthStart{AuthorizationURL:s.outlook.AuthorizationURL(state)},nil
}

func(s *ConnectionService)FinishOAuth(ctx context.Context,provider,code,state string)(integrationresponse.OAuthCompletion,error){
 if code==""||len(state)<40{return integrationresponse.OAuthCompletion{},ErrOAuthInvalid}
 oauthState,err:=s.repository.ConsumeOAuthState(ctx,hashState(state),time.Now().UTC())
 if err!=nil{
  if errors.Is(err,integrationrepository.ErrOAuthStateInvalid){return integrationresponse.OAuthCompletion{},ErrOAuthInvalid}
  return integrationresponse.OAuthCompletion{},err
 }
 connection,err:=s.repository.GetByID(ctx,oauthState.ConnectionID)
 if err!=nil{return integrationresponse.OAuthCompletion{},err}
 if connection.Provider!=provider || connection.Status!=entity.StatusPendingSetup{
  return integrationresponse.OAuthCompletion{},ErrOAuthInvalid
 }
 var token providers.OAuthToken
 var mailbox string
 switch provider{
 case "GMAIL":
  if s.gmail==nil{return integrationresponse.OAuthCompletion{},ErrOAuthUnavailable}
  token,err=s.gmail.Exchange(ctx,code)
  if err==nil && token.RefreshToken!=""{mailbox,err=s.gmail.Mailbox(ctx,token)}
 case "OUTLOOK":
  if s.outlook==nil{return integrationresponse.OAuthCompletion{},ErrOAuthUnavailable}
  token,err=s.outlook.Exchange(ctx,code)
  if err==nil && token.RefreshToken!=""{mailbox,err=s.outlook.Mailbox(ctx,token)}
 default:return integrationresponse.OAuthCompletion{},ErrOAuthUnavailable
 }
 if err!=nil{return integrationresponse.OAuthCompletion{},err}
 if mailbox==""||token.RefreshToken==""{return integrationresponse.OAuthCompletion{},ErrOAuthInvalid}
 encrypted,err:=s.cipher.Seal(token)
 if err!=nil{return integrationresponse.OAuthCompletion{},err}
 if err=s.repository.Activate(ctx,connection.ID,provider,mailbox,encrypted);err!=nil{
  return integrationresponse.OAuthCompletion{},err
 }
 return integrationresponse.OAuthCompletion{Status:entity.StatusActive},nil
}

func hashState(raw string)string{
 sum:=sha256.Sum256([]byte(raw))
 return hex.EncodeToString(sum[:])
}

func(s *ConnectionService)validProviderToken(ctx context.Context,connectionID,provider string)(providers.OAuthToken,error){
 var empty providers.OAuthToken
 encrypted,err:=s.repository.Credential(ctx,connectionID)
 if err!=nil{return empty,err}
 token,err:=s.cipher.Open(encrypted)
 if err!=nil{return empty,err}
 if time.Until(token.Expiry)>time.Minute{return token,nil}
 switch provider{
 case "GMAIL":
  if s.gmail==nil{return empty,ErrOAuthUnavailable}
  token,err=s.gmail.Refresh(ctx,token)
 case "OUTLOOK":
  if s.outlook==nil{return empty,ErrOAuthUnavailable}
  token,err=s.outlook.Refresh(ctx,token)
 default:return empty,ErrOAuthUnavailable
 }
 if err!=nil{return empty,err}
 encrypted,err=s.cipher.Seal(token)
 if err!=nil{return empty,err}
 if err=s.repository.UpdateCredential(ctx,connectionID,encrypted);err!=nil{return empty,err}
 return token,nil
}

// Existing Gmail watch/connection probes use this compatibility method.
// Provider dispatch for normal actions goes through validProviderToken.
func(s *ConnectionService)validToken(ctx context.Context,connectionID string)(providers.OAuthToken,error){
 return s.validProviderToken(ctx,connectionID,"GMAIL")
}
