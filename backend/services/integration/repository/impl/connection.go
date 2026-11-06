package impl

import (
 "context"
 "errors"
 "time"

 "gorm.io/gorm"

 coretransaction "grovehub/libs/core/database/transaction"
 "grovehub/services/integration/entity"
 integrationrepository "grovehub/services/integration/repository"
)

type ConnectionRepository struct { db *gorm.DB }
func NewConnectionRepository(db *gorm.DB) integrationrepository.ConnectionRepository {
 return &ConnectionRepository{db:db}
}
func (r *ConnectionRepository) Create(ctx context.Context, c entity.Connection) error {
 return coretransaction.Database(ctx,r.db).Create(&c).Error
}
func (r *ConnectionRepository) Get(ctx context.Context, tenantID,id string)(entity.Connection,error){
 var c entity.Connection
 err:=coretransaction.Database(ctx,r.db).Where("tenant_id = ? AND id = ?",tenantID,id).Take(&c).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return c,integrationrepository.ErrConnectionNotFound}
 return c,err
}
func (r *ConnectionRepository) List(ctx context.Context,tenantID string)([]entity.Connection,error){
 c:=make([]entity.Connection,0)
 err:=coretransaction.Database(ctx,r.db).Where("tenant_id = ?",tenantID).
 Order("created_at DESC, id DESC").Limit(200).Find(&c).Error
 return c,err
}
func (r *ConnectionRepository) Disconnect(ctx context.Context,tenantID,id string)(entity.Connection,error){
 db:=coretransaction.Database(ctx,r.db)
 result:=db.Model(&entity.Connection{}).Where("tenant_id = ? AND id = ?",tenantID,id).
 Updates(map[string]any{"status":entity.StatusDisconnected,"credential_ref":nil,"updated_at":time.Now().UTC()})
 if result.Error!=nil{return entity.Connection{},result.Error}
 if result.RowsAffected==0{return entity.Connection{},integrationrepository.ErrConnectionNotFound}
 return r.Get(ctx,tenantID,id)
}
