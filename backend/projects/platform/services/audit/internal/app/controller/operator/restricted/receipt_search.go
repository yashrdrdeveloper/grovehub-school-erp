package restricted

import (
 "context"
 "errors"
 "github.com/google/uuid"
 "encoding/json"
 "io"
 "net/http"
 "strconv"
 "time"

 "gorm.io/gorm"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const ReceiptSearchPattern="POST /audit/restricted/school-mutation-receipts/search"
const ReceiptReadScope="internal.audit.school.mutation-receipt.read"
const ReceiptReader="internal-tool-school-receipt-reader"

type ReceiptView struct {
 ExecutionID string `json:"executionId"`
 SchoolID string `json:"schoolId"`
 Operation string `json:"operation"`
 PreviousVersion string `json:"previousVersion"`
 CorrelationID *string `json:"correlationId,omitempty"`
 OccurredAt time.Time `json:"occurredAt"`
}
type ReceiptPage struct {
 Records []ReceiptView `json:"records"`
 PageNumber int `json:"pageNumber"`
 RowSize int `json:"rowSize"`
 HasMore bool `json:"hasMore"`
}
type ReceiptSearchInput struct {
 SearchInput
 SchoolID string `json:"schoolId,omitempty"`
 Operation string `json:"operation,omitempty"`
}
func(in ReceiptSearchInput) Valid()bool{
 if in.PageNumber<1||in.PageNumber>10000||in.RowSize<1||in.RowSize>100{return false}
 if in.SchoolID!=""{id,err:=uuid.Parse(in.SchoolID);if err!=nil||id==uuid.Nil||id.String()!=in.SchoolID{return false}}
 switch in.Operation{case "","RENAME","ACTIVATE","DISABLE":return true;default:return false}
}
func SearchReceipts(ctx context.Context,db *gorm.DB,in SearchInput)(ReceiptPage,error){
 return SearchFilteredReceipts(ctx,db,ReceiptSearchInput{SearchInput:in})
}
func SearchFilteredReceipts(ctx context.Context,db *gorm.DB,in ReceiptSearchInput)(ReceiptPage,error){
 if !in.Valid(){return ReceiptPage{},errors.New("invalid receipt search")}
 query:="SELECT execution_id,school_id,previous_version,correlation_id,occurred_at,operation FROM audit.school_mutation_receipts WHERE TRUE"
 args:=[]any{}
 if in.SchoolID!=""{query+=" AND school_id = ?";args=append(args,in.SchoolID)}
 if in.Operation!=""{query+=" AND operation = ?";args=append(args,in.Operation)}
 query+=" ORDER BY occurred_at DESC,execution_id DESC LIMIT ? OFFSET ?"
 args=append(args,in.RowSize+1,(in.PageNumber-1)*in.RowSize)
 rows,err:=db.WithContext(ctx).Raw(query,args...).Rows()
 if err!=nil{return ReceiptPage{},err};defer rows.Close()
 p:=ReceiptPage{Records:[]ReceiptView{},PageNumber:in.PageNumber,RowSize:in.RowSize}
 for rows.Next(){
  var v ReceiptView;var version int64
  if err=rows.Scan(&v.ExecutionID,&v.SchoolID,&version,&v.CorrelationID,&v.OccurredAt,&v.Operation);err!=nil{return ReceiptPage{},err}
  v.PreviousVersion=strconv.FormatInt(version,10);p.Records=append(p.Records,v)
 }
 if err=rows.Err();err!=nil{return ReceiptPage{},err}
 if len(p.Records)>in.RowSize{p.HasMore=true;p.Records=p.Records[:in.RowSize]}
 return p,nil
}
func ReceiptSearchHandler(db *gorm.DB)http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  claims,ok:=core.WorkloadFrom(r.Context())
  if !ok||claims.Subject!=ReceiptReader||len(claims.Scopes)!=1||claims.Scopes[0]!=ReceiptReadScope||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  var in ReceiptSearchInput;var extra any
  d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,4096));d.DisallowUnknownFields()
  if d.Decode(&in)!=nil||d.Decode(&extra)!=io.EOF||!in.Valid(){w.WriteHeader(400);return}
  ctx,cancel:=context.WithTimeout(r.Context(),5*time.Second);defer cancel()
  page,err:=SearchFilteredReceipts(ctx,db,in);if err!=nil{w.WriteHeader(503);return}
  w.Header().Set("Content-Type","application/json");_ =json.NewEncoder(w).Encode(page)
 }
}
