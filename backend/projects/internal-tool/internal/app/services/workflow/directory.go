package workflow

import (
 "context"
 "sort"
 "time"

 "github.com/yashrdrdeveloper/grovehub-school-erp/backend/projects/internal-tool/internal/app/entity"
)

type SchoolRecord struct{ID,Code,Name,Status string}
type Directory func(context.Context,[]string)(map[string]SchoolRecord,error)

// WithDirectory is called once at startup, before serving requests.
func(s *Service) WithDirectory(directory Directory)*Service{s.directory=directory;return s}
func(s *Service) ResolvedSchoolGrants(ctx context.Context,actor entity.Operator)([]SchoolGrantView,error){
 grants:=s.SchoolGrants(actor)
 if s.directory==nil||len(grants)==0{return grants,nil}
 ids:=map[string]bool{}
 for _,g:=range grants{ids[g.SchoolID]=true}
 ordered:=make([]string,0,len(ids));for id:=range ids{ordered=append(ordered,id)};sort.Strings(ordered)
 probe,cancel:=context.WithTimeout(ctx,4*time.Second);defer cancel()
 records,err:=s.directory(probe,ordered);if err!=nil{return nil,err}
 result:=[]SchoolGrantView{}
 for _,grant:=range grants{
  row,ok:=records[grant.SchoolID]
  if !ok||row.Status!="ACTIVE"{continue}
  grant.SchoolName=row.Name;grant.SchoolCode=row.Code;result=append(result,grant)
 }
 return result,nil
}
func(s *Service) checkDirectory(ctx context.Context,id string)error{
 if s.directory==nil{return nil}
 probe,cancel:=context.WithTimeout(ctx,4*time.Second);defer cancel()
 rows,err:=s.directory(probe,[]string{id});if err!=nil{return err}
 if rows[id].ID!=id||rows[id].Status!="ACTIVE"{return ErrForbidden}
 return nil
}
