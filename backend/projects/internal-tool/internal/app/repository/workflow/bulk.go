package workflow

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/http/correlation"
 "sort"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

type BulkItem struct{
 ID string `json:"id"`
 Version int64 `json:"version"`
}

func(r *Repository) BulkApprove(ctx context.Context,actor entity.Operator,items []BulkItem,comment string)([]entity.Approval,error){
 ctx=correlation.Ensure(ctx)
 if actor.Persona!=entity.L2{return nil,ErrForbidden}
 tx,err:=r.db.BeginTx(ctx,nil);if err!=nil{return nil,err};defer tx.Rollback()
 var active bool
 if err=tx.QueryRowContext(ctx,"SELECT EXISTS(SELECT 1 FROM internal_tool.operators WHERE id=$1 AND persona='L2' AND status='ACTIVE')",actor.ID).Scan(&active);err!=nil{return nil,err}
 if !active{return nil,ErrForbidden}
 ordered:=append([]BulkItem(nil),items...)
 sort.Slice(ordered,func(i,j int)bool{return ordered[i].ID<ordered[j].ID})
 result:=make([]entity.Approval,0,len(ordered))
 for _,item:=range ordered{
  a,err:=decide(ctx,tx,actor,item.ID,item.Version,"APPROVED",comment);if err!=nil{return nil,err}
  result=append(result,a)
 }
 if err=tx.Commit();err!=nil{return nil,err}
 return result,nil
}
