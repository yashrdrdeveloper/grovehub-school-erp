package gmail

import (
 "context"
 "crypto"
 "crypto/rsa"
 "crypto/sha256"
 "encoding/base64"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "math/big"
 "net/http"
 "strings"
 "sync"
 "time"
)

const jwksEndpoint="https://www.googleapis.com/oauth2/v3/certs"

type PushVerifier struct{
 audience,email string
 client *http.Client
 mu sync.Mutex
 keys map[string]*rsa.PublicKey
 expiry time.Time
}
func NewPushVerifier(audience,email string)*PushVerifier{
 return &PushVerifier{audience:audience,email:strings.ToLower(email),client:&http.Client{Timeout:8*time.Second}}
}
type PushNotification struct{EmailAddress string `json:"emailAddress"`;HistoryID string `json:"historyId"`}
func DecodePush(body io.Reader)(PushNotification,error){
 var envelope struct{Message struct{Data string `json:"data"`} `json:"message"`}
 if err:=json.NewDecoder(io.LimitReader(body,32768)).Decode(&envelope);err!=nil{return PushNotification{},err}
 if len(envelope.Message.Data)==0||len(envelope.Message.Data)>8192{return PushNotification{},errors.New("invalid Pub/Sub payload")}
 raw,err:=base64.StdEncoding.DecodeString(envelope.Message.Data)
 if err!=nil{raw,err=base64.RawURLEncoding.DecodeString(envelope.Message.Data)}
 if err!=nil{return PushNotification{},err}
 var data PushNotification
 if err:=json.Unmarshal(raw,&data);err!=nil{return data,err}
 data.EmailAddress=strings.ToLower(strings.TrimSpace(data.EmailAddress))
 if data.EmailAddress==""||data.HistoryID==""{return data,errors.New("missing Gmail push fields")}
 return data,nil
}
func(v *PushVerifier)Verify(ctx context.Context,authorization string)error{
 if v==nil||v.audience==""||v.email==""||!strings.HasPrefix(authorization,"Bearer "){return errors.New("push token unavailable")}
 parts:=strings.Split(strings.TrimPrefix(authorization,"Bearer "),".")
 if len(parts)!=3{return errors.New("malformed JWT")}
 headerBytes,err:=base64.RawURLEncoding.DecodeString(parts[0]);if err!=nil{return err}
 var header struct{Alg string `json:"alg"`;Kid string `json:"kid"`}
 if err=json.Unmarshal(headerBytes,&header);err!=nil{return err}
 if header.Alg!="RS256"||header.Kid==""{return errors.New("invalid signing algorithm")}
 keys,err:=v.publicKeys(ctx);if err!=nil{return err}
 key,ok:=keys[header.Kid];if !ok{return errors.New("unknown signing key")}
 signature,err:=base64.RawURLEncoding.DecodeString(parts[2]);if err!=nil{return err}
 digest:=sha256.Sum256([]byte(parts[0]+"."+parts[1]))
 if rsa.VerifyPKCS1v15(key,crypto.SHA256,digest[:],signature)!=nil{return errors.New("invalid token signature")}
 body,err:=base64.RawURLEncoding.DecodeString(parts[1]);if err!=nil{return err}
 var claims struct{
  Iss string `json:"iss"`
  Aud string `json:"aud"`
  Email string `json:"email"`
  Verified bool `json:"email_verified"`
  Exp int64 `json:"exp"`
  Iat int64 `json:"iat"`
  Nbf int64 `json:"nbf"`
 }
 if err=json.Unmarshal(body,&claims);err!=nil{return err}
 now:=time.Now().Unix()
 if (claims.Iss!="https://accounts.google.com"&&claims.Iss!="accounts.google.com")||
  claims.Aud!=v.audience||!strings.EqualFold(claims.Email,v.email)||!claims.Verified||
  claims.Exp<=now||claims.Iat<=0||claims.Iat>now+60||claims.Nbf>now{
  return errors.New("invalid signed Pub/Sub identity claims")
 }
 return nil
}
func(v *PushVerifier)publicKeys(ctx context.Context)(map[string]*rsa.PublicKey,error){
 v.mu.Lock();defer v.mu.Unlock()
 if time.Now().Before(v.expiry)&&len(v.keys)>0{return v.keys,nil}
 request,err:=http.NewRequestWithContext(ctx,http.MethodGet,jwksEndpoint,nil);if err!=nil{return nil,err}
 response,err:=v.client.Do(request);if err!=nil{return nil,err}
 defer response.Body.Close()
 if response.StatusCode!=http.StatusOK{return nil,fmt.Errorf("Google JWKS unavailable (%d)",response.StatusCode)}
 var doc struct{Keys []struct{
  Kid string `json:"kid"`;Kty string `json:"kty"`;Alg string `json:"alg"`;N string `json:"n"`;E string `json:"e"`
 } `json:"keys"`}
 if err=json.NewDecoder(io.LimitReader(response.Body,65536)).Decode(&doc);err!=nil{return nil,err}
 keys:=map[string]*rsa.PublicKey{}
 for _,entry:=range doc.Keys{
  if entry.Kty!="RSA"||entry.Alg!="RS256"||entry.Kid==""{continue}
  n,err:=base64.RawURLEncoding.DecodeString(entry.N);if err!=nil||len(n)<256{continue}
  e,err:=base64.RawURLEncoding.DecodeString(entry.E);if err!=nil||len(e)==0||len(e)>4{continue}
  exp:=0;for _,byteValue:=range e{exp=exp*256+int(byteValue)}
  if exp<3{continue}
  keys[entry.Kid]=&rsa.PublicKey{N:new(big.Int).SetBytes(n),E:exp}
 }
 if len(keys)==0{return nil,errors.New("Google JWKS contains no usable keys")}
 v.keys=keys;v.expiry=time.Now().Add(time.Hour)
 return keys,nil
}
