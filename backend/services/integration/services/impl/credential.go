package impl

import (
 "crypto/aes"
 "crypto/cipher"
 "crypto/rand"
 "encoding/json"
 "errors"

 "grovehub/services/integration/providers"
)
type tokenCipher struct{key []byte}
func newTokenCipher(key []byte)(tokenCipher,error){
 if len(key)!=32{return tokenCipher{},errors.New("integration encryption key must have 32 bytes")}
 return tokenCipher{key:append([]byte(nil),key...)},nil
}
func(c tokenCipher)Seal(token providers.OAuthToken)([]byte,error){
 block,err:=aes.NewCipher(c.key);if err!=nil{return nil,err}
 gcm,err:=cipher.NewGCM(block);if err!=nil{return nil,err}
 nonce:=make([]byte,gcm.NonceSize())
 if _,err=rand.Read(nonce);err!=nil{return nil,err}
 payload,err:=json.Marshal(token);if err!=nil{return nil,err}
 return gcm.Seal(nonce,nonce,payload,nil),nil
}
func(c tokenCipher)Open(ciphertext []byte)(providers.OAuthToken,error){
 var token providers.OAuthToken
 block,err:=aes.NewCipher(c.key);if err!=nil{return token,err}
 gcm,err:=cipher.NewGCM(block);if err!=nil{return token,err}
 if len(ciphertext)<gcm.NonceSize(){return token,errors.New("invalid credential ciphertext")}
 payload,err:=gcm.Open(nil,ciphertext[:gcm.NonceSize()],ciphertext[gcm.NonceSize():],nil)
 if err!=nil{return token,err}
 err=json.Unmarshal(payload,&token)
 return token,err
}
