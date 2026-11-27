package workflow

import (
 "context"
 "strings"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
 repo "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/repository/workflow"
)

type BulkApprovalInput struct{
 Items []repo.BulkItem `json:"items"`
 Comment string `json:"comment"`
}
func(s *Service) BulkApprove(ctx context.Context,actor entity.Operator,in BulkApprovalInput)([]entity.Approval,error){
 if actor.Persona!=entity.L2{return nil,ErrForbidden}
 in.Comment=strings.TrimSpace(in.Comment)
 if len(in.Items)==0||len(in.Items)>50||len(in.Comment)>1000{return nil,ErrInvalid}
 seen:=map[string]bool{}
 for _,item:=range in.Items{
  if !validID(item.ID)||item.Version<1||seen[item.ID]{return nil,ErrInvalid}
  seen[item.ID]=true
  if err:=s.checkSchoolReview(ctx,actor,item.ID);err!=nil{return nil,err}
 }
 return s.repo.BulkApprove(ctx,actor,in.Items,in.Comment)
}
