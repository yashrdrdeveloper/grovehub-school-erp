package delivery

import (
 "crypto/ed25519"
 "encoding/base64"
 "encoding/json"
 "errors"
 "io"
 "strings"

 core "grovehub/libs/core/gateway"
)

type ReaderGrant struct{Key core.WorkloadKeyConfig `json:"key"`}

// Configure validates all credentials before any worker starts. The relay
// must never reuse an ingress credential or another school's publisher key.
func Configure(raw string,legacy Config,workloadID,workloadPrivate,readRaw,renameRaw string,statusRaw ...string)([]Config,error){
 configs:=[]Config{}
 if strings.TrimSpace(raw)!=""{
  if legacy!=(Config{})||len(raw)>256*1024{return nil,errors.New("conflicting receipt configuration")}
  d:=json.NewDecoder(strings.NewReader(raw));d.DisallowUnknownFields()
  var extra any
  if d.Decode(&configs)!=nil||configs==nil||len(configs)>100||d.Decode(&extra)!=io.EOF{return nil,errors.New("invalid receipt publisher configuration")}
 }else if legacy!=(Config{}){configs=append(configs,legacy)}
 if len(configs)==0{return configs,nil}
 ids:=map[string]bool{}
 keys:=map[string]bool{}
 schools:=map[string]bool{}
 reserve:=func(id,pub string){if id!=""{ids[id]=true};if pub!=""{keys[pub]=true}}
 if workloadID!=""{ids[workloadID]=true}
 if b,err:=base64.RawURLEncoding.DecodeString(workloadPrivate);err==nil&&len(b)==ed25519.PrivateKeySize{reserve(workloadID,base64.RawURLEncoding.EncodeToString(b[32:]))}
 var readers []ReaderGrant
 if strings.TrimSpace(readRaw)!=""&&json.Unmarshal([]byte(readRaw),&readers)!=nil{return nil,errors.New("invalid registry reader configuration")}
 for _,r:=range readers{reserve(r.Key.KeyID,r.Key.PublicKey)}
 var renamers []core.WorkloadKeyConfig
 if strings.TrimSpace(renameRaw)!=""&&json.Unmarshal([]byte(renameRaw),&renamers)!=nil{return nil,errors.New("invalid registry rename configuration")}
 for _,k:=range renamers{reserve(k.KeyID,k.PublicKey)}
 for _,raw:=range statusRaw{
  var statusKeys []core.WorkloadKeyConfig
  if strings.TrimSpace(raw)!=""&&json.Unmarshal([]byte(raw),&statusKeys)!=nil{return nil,errors.New("invalid registry status configuration")}
  for _,k:=range statusKeys{reserve(k.KeyID,k.PublicKey)}
 }
 for _,c:=range configs{
  operation:=c.Operation
  if operation==""{operation="rename"}
  if operation!="rename"&&operation!="status"{return nil,errors.New("invalid receipt operation")}
  private,err:=base64.RawURLEncoding.DecodeString(c.PrivateKey)
  if err!=nil||len(private)!=ed25519.PrivateKeySize||schools[c.SchoolID+":"+operation]{return nil,errors.New("invalid or duplicate receipt publisher")}
  pub:=base64.RawURLEncoding.EncodeToString(private[32:])
  if ids[c.KeyID]||keys[pub]{return nil,errors.New("receipt publisher credential must be independent")}
  if operation=="status"{
   if publisher,err:=StatusClient(c);err!=nil||publisher==nil{return nil,errors.New("invalid status receipt publisher")}
  }else if publisher,err:=Client(c);err!=nil||publisher==nil{return nil,errors.New("invalid rename receipt publisher")}
  reserve(c.KeyID,pub);schools[c.SchoolID+":"+operation]=true
 }
 return configs,nil
}
