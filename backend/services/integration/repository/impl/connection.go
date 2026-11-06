package impl

import (
 "context"
 "errors"
 "time"
 "gorm.io/gorm"
 "gorm.io/gorm/clause"

 coretransaction "grovehub/libs/core/database/transaction"
 "grovehub/services/integration/entity"
 integrationrepository "grovehub/services/integration/repository"
)

type ConnectionRepository struct{db *gorm.DB}
func NewConnectionRepository(db *gorm.DB) integrationrepository.ConnectionRepository{return &ConnectionRepository{db:db}}
func(r *ConnectionRepository)Create(ctx context.Context,c entity.Connection)error{return coretransaction.Database(ctx,r.db).Create(&c).Error}
func(r *ConnectionRepository)Get(ctx context.Context,tenantID,id string)(entity.Connection,error){
 var c entity.Connection
 err:=coretransaction.Database(ctx,r.db).Where("tenant_id = ? AND id = ?",tenantID,id).Take(&c).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return c,integrationrepository.ErrConnectionNotFound}
 return c,err
}
func(r *ConnectionRepository)List(ctx context.Context,tenantID string)([]entity.Connection,error){
 rows:=make([]entity.Connection,0)
 err:=coretransaction.Database(ctx,r.db).Where("tenant_id = ?",tenantID).Order("created_at DESC, id DESC").Limit(200).Find(&rows).Error
 return rows,err
}
func(r *ConnectionRepository)Disconnect(ctx context.Context,tenantID,id string)(entity.Connection,error){
 err:=r.db.WithContext(ctx).Transaction(func(tx *gorm.DB)error{
  result:=tx.Model(&entity.Connection{}).Where("tenant_id = ? AND id = ?",tenantID,id).Updates(map[string]any{
   "status":entity.StatusDisconnected,"credential_ref":nil,"updated_at":time.Now().UTC(),
  })
  if result.Error!=nil{return result.Error}
  if result.RowsAffected==0{return integrationrepository.ErrConnectionNotFound}
  if err:=tx.Where("connection_id = ?",id).Delete(&entity.OAuthCredential{}).Error;err!=nil{return err}
  return tx.Where("connection_id = ?",id).Delete(&entity.OAuthState{}).Error
 })
 if err!=nil{return entity.Connection{},err}
 return r.Get(ctx,tenantID,id)
}
func(r *ConnectionRepository)SaveOAuthState(ctx context.Context,state entity.OAuthState)error{
 return r.db.WithContext(ctx).Create(&state).Error
}
func(r *ConnectionRepository)ConsumeOAuthState(ctx context.Context,hash string,now time.Time)(entity.OAuthState,error){
 var state entity.OAuthState
 err:=r.db.WithContext(ctx).Clauses(clause.Returning{}).
  Where("hash = ? AND expires_at > ?",hash,now).
  Delete(&state).Error
 if err!=nil{return entity.OAuthState{},err}
 if state.Hash==""{return entity.OAuthState{},integrationrepository.ErrOAuthStateInvalid}
 return state,nil
}
func(r *ConnectionRepository)Activate(ctx context.Context,connectionID,externalID string,ciphertext []byte)error{
 return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB)error{
  result:=tx.Model(&entity.Connection{}).Where("id = ? AND provider = ? AND status = ?",connectionID,"GMAIL",entity.StatusPendingSetup).
   Updates(map[string]any{"external_id":externalID,"status":entity.StatusActive,"updated_at":time.Now().UTC()})
  if result.Error!=nil{return result.Error}
  if result.RowsAffected==0{return integrationrepository.ErrConnectionNotFound}
  return tx.Create(&entity.OAuthCredential{ConnectionID:connectionID,Ciphertext:ciphertext,UpdatedAt:time.Now().UTC()}).Error
 })
}
func(r *ConnectionRepository)Credential(ctx context.Context,connectionID string)([]byte,error){
 var credential entity.OAuthCredential
 err:=r.db.WithContext(ctx).Where("connection_id = ?",connectionID).Take(&credential).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return nil,integrationrepository.ErrConnectionNotFound}
 return credential.Ciphertext,err
}
func(r *ConnectionRepository)UpdateCredential(ctx context.Context,connectionID string,ciphertext []byte)error{
 return r.db.WithContext(ctx).Model(&entity.OAuthCredential{}).Where("connection_id = ?",connectionID).
 Updates(map[string]any{"ciphertext":ciphertext,"updated_at":time.Now().UTC()}).Error
}
func(r *ConnectionRepository)SaveTest(ctx context.Context,test entity.ConnectionTest)error{
 return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB)error{
  if err:=tx.Create(&test).Error;err!=nil{return err}
  return tx.Model(&entity.Connection{}).Where("id = ? AND (last_tested_at IS NULL OR last_tested_at <= ?)",test.ConnectionID,test.CompletedAt).
   Updates(map[string]any{"last_tested_at":test.CompletedAt,"last_test_status":test.Status,"updated_at":time.Now().UTC()}).Error
 })
}
func(r *ConnectionRepository)ListTests(ctx context.Context,connectionID string)([]entity.ConnectionTest,error){
 tests:=make([]entity.ConnectionTest,0)
 err:=r.db.WithContext(ctx).Where("connection_id = ?",connectionID).
 Order("completed_at DESC, id DESC").Limit(100).Find(&tests).Error
 return tests,err
}
