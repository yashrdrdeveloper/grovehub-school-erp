package restricted

import (
 "context"
 "net/http"
 "time"

 "gorm.io/gorm"
 core "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/gateway"
)

const SchoolCheckPattern="GET /audit/restricted/schools/{schoolId}/activity-check"

// SchoolCheck exposes only whether this exact tenant has recorded activity.
// No actor data, event content, tenant enumeration or write access is provided.
func SchoolCheck(ctx context.Context,db *gorm.DB,school string)(bool,error){
 var exists bool
 err:=db.WithContext(ctx).Raw("SELECT EXISTS(SELECT 1 FROM audit.activities WHERE tenant_id=?)",school).Scan(&exists).Error
 return exists,err
}
func SchoolCheckHandler(db *gorm.DB)http.HandlerFunc{
 return func(w http.ResponseWriter,r *http.Request){
  w.Header().Set("Cache-Control","no-store")
  school:=r.PathValue("schoolId")
  claims,ok:=core.WorkloadFrom(r.Context())
  if !ok||!validID(school)||claims.Subject!="internal-tool-school:"+school||len(claims.Scopes)!=1||claims.Scopes[0]!="internal.audit.school.activity-check"||claims.ActorGID!=""||claims.TenantID!=""||claims.SessionID!=""{w.WriteHeader(403);return}
  ctx,cancel:=context.WithTimeout(r.Context(),3*time.Second);defer cancel()
  exists,err:=SchoolCheck(ctx,db,school)
  if err!=nil{w.WriteHeader(503);return}
  if !exists{w.WriteHeader(404);return}
  w.WriteHeader(204)
 }
}
