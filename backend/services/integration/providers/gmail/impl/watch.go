package impl

import (
 "context"
 "errors"
 "net/http"
 "strconv"
 "strings"
 "time"
)
func(c *Client)Watch(ctx context.Context,t Token,topic string)(string,time.Time,error){
 if !strings.HasPrefix(topic,"projects/")||!strings.Contains(topic,"/topics/"){return "",time.Time{},errors.New("invalid Pub/Sub topic")}
 var p struct{HistoryID string `json:"historyId"`; Expiration string `json:"expiration"`}
 err:=c.do(ctx,t,http.MethodPost,"/watch",map[string]string{"topicName":topic},&p)
 if err!=nil{return "",time.Time{},err}
 millis,err:=strconv.ParseInt(p.Expiration,10,64)
 if err!=nil||p.HistoryID==""||millis<=0{return "",time.Time{},errors.New("invalid Gmail watch result")}
 return p.HistoryID,time.UnixMilli(millis).UTC(),nil
}