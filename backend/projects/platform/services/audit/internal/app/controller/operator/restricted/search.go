package restricted

import (
 "context"
 "encoding/json"
 "io"
 "net/http"
 "time"

 "gorm.io/gorm"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
 "grovehub/services/audit/client/producer"
)

const SearchPattern="POST /audit/restricted/operator-activities/search"
const ReadScope="internal.audit.operator.read"

type SearchInput struct{
 PageNumber int `json:"pageNumber"`
 RowSize int `json:"rowSize"`
}
type SearchPage struct{
 Records []producer.OperatorActivity `json:"records"`
 PageNumber int `json:"pageNumber"`
 RowSize int `json:"rowSize"`
 HasMore bool `json:"hasMore"`
}
func Search(ctx context.Context,db *gorm.DB,in SearchInput)(SearchPage,error){
 rows,err:=db.WithContext(ctx).Raw("SELECT event_id,actor_id,resource_id,operation,occurred_at FROM audit.operator_activities ORDER BY occurred_at DESC,event_id DESC LIMIT ? OFFSET ?",in.RowSize+1,(in.PageNumber-1)*in.RowSize).Rows()
 if err!=nil{return SearchPage{},err};defer rows.Close()
 p:=SearchPage{Records:[]producer.OperatorActivity{},PageNumber:in.PageNumber,RowSize:in.RowSize}
 for rows.Next(){var a producer.OperatorActivity;if err=rows.Scan(&a.EventID,&a.ActorID,&a.ResourceID,&a.Operation,&a.OccurredAt);err!=nil{return SearchPage{},err};p.Records=append(p.Records,a)}
 if err=rows.Err();err!=nil{return SearchPage{},err}
 if len(p.Records)>in.RowSize{p.HasMore=true;p.Records=p.Records[:in.RowSize]}
 return p,nil
}
func SearchHandler(db *gorm.DB)http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  claims,ok:=core.WorkloadFrom(r.Context())
  if !ok||claims.Subject!="internal-tool-audit-reader"||len(claims.Scopes)!=1||claims.Scopes[0]!=ReadScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  var in SearchInput;var extra any
  d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));d.DisallowUnknownFields()
  if d.Decode(&in)!=nil||d.Decode(&extra)!=io.EOF||in.PageNumber<1||in.PageNumber>10000||in.RowSize<1||in.RowSize>100{w.WriteHeader(400);return}
  ctx,cancel:=context.WithTimeout(r.Context(),5*time.Second);defer cancel()
  page,err:=Search(ctx,db,in);if err!=nil{w.WriteHeader(503);return}
  w.Header().Set("Content-Type","application/json");_ = json.NewEncoder(w).Encode(page)
 }
}
