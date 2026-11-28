package workflow

import (
 "context"
 "strings"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
type ReconcileInput struct{
 Outcome string `json:"outcome"`
 Evidence string `json:"evidence"`
}
func(s *Service) Reconcile(ctx context.Context,actor entity.Operator,id string,in ReconcileInput)(entity.Reconciliation,error){
 if actor.Persona!=entity.L2{return entity.Reconciliation{},ErrForbidden}
 in.Evidence=strings.TrimSpace(in.Evidence)
 if !validID(id)||len(in.Evidence)==0||len(in.Evidence)>1000||(in.Outcome!="CONFIRMED_SUCCEEDED"&&in.Outcome!="CONFIRMED_FAILED"&&in.Outcome!="UNRESOLVED"){return entity.Reconciliation{},ErrInvalid}
 return s.repo.Reconcile(ctx,actor,id,in.Outcome,in.Evidence)
}
func(s *Service) Reconciliation(ctx context.Context,actor entity.Operator,id string)(entity.Reconciliation,error){
 if !validID(id){return entity.Reconciliation{},ErrInvalid}
 return s.repo.Reconciliation(ctx,actor,id)
}
