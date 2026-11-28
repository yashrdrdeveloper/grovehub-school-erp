package workflow

import (
 "context"
 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)
func(s *Service) SearchEnvironments(ctx context.Context,p entity.PageRequest)(entity.Page[entity.Environment],error){
 p,ok:=p.Normalize();if !ok{return entity.Page[entity.Environment]{},ErrInvalid};return s.repo.SearchEnvironments(ctx,p)
}
func(s *Service) SearchTemplates(ctx context.Context,p entity.PageRequest)(entity.Page[entity.Template],error){
 p,ok:=p.Normalize();if !ok{return entity.Page[entity.Template]{},ErrInvalid};return s.repo.SearchTemplates(ctx,p)
}
func(s *Service) SearchApprovals(ctx context.Context,a entity.Operator,p entity.PageRequest,inbox bool)(entity.Page[entity.Approval],error){
 p,ok:=p.Normalize();if !ok{return entity.Page[entity.Approval]{},ErrInvalid};return s.repo.SearchApprovals(ctx,a,p,inbox)
}
func(s *Service) SearchEvents(ctx context.Context,a entity.Operator,p entity.PageRequest)(entity.Page[entity.WorkflowEvent],error){
 p,ok:=p.Normalize();if !ok{return entity.Page[entity.WorkflowEvent]{},ErrInvalid};return s.repo.SearchEvents(ctx,a,p)
}
