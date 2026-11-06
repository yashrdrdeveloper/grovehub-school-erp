package config

import (
 "encoding/base64"
 "errors"
 "os"
)
func CredentialEncryptionKey()([]byte,error){
 value:=os.Getenv("GROVEHUB_INTEGRATION_ENCRYPTION_KEY")
 key,err:=base64.StdEncoding.DecodeString(value)
 if err!=nil||len(key)!=32{return nil,errors.New("integration encryption key must be base64-encoded 32 bytes")}
 return key,nil
}
