package creation

import (
 "context"
 "crypto/sha256"
 "encoding/hex"
 "encoding/json"
 "errors"
 "strings"
 "time"

 "gorm.io/gorm"
 "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/transaction"
 "grovehub/services/audit/client/producer"
 "grovehub/services/accounting/internal/app/dto"
 "grovehub/services/accounting/internal/app/entity"
 "grovehub/services/accounting/internal/app/services/accountingentity"
 "grovehub/services/accounting/internal/app/services/accountingbook"
 "grovehub/services/accounting/internal/app/validators"
)

var ErrConflict = errors.New("idempotency key reused with different request")

type Service struct {
 db *gorm.DB
 unit transaction.UnitOfWork
 audit producer.Client
 entities accountingentity.Service
 books accountingbook.Service
}

func New(db *gorm.DB, unit transaction.UnitOfWork, audit producer.Client, entities accountingentity.Service, books accountingbook.Service) *Service {
 return &Service{db:db, unit:unit, audit:audit, entities:entities, books:books}
}

func (s *Service) CreateEntity(ctx context.Context, tenant, actor, key string, input dto.CreateAccountingEntity) (entity.AccountingEntity,error) {
 var result entity.AccountingEntity
 err:=s.run(ctx,tenant,actor,key,"ENTITY_CREATE","ACCOUNTING_ENTITY",input,func(ctx context.Context, existing string)(string,time.Time,error){
  var err error
  if existing!="" { result,err=s.entities.FindByID(ctx,tenant,existing) } else { result,err=s.entities.Create(ctx,tenant,input) }
  return result.EntityID,result.CreatedAt,err
 })
 if err!=nil { return entity.AccountingEntity{},err }
 return result,nil
}

func (s *Service) CreateBook(ctx context.Context, tenant, actor, key string, input dto.CreateAccountingBook) (entity.AccountingBook,error) {
 var result entity.AccountingBook
 err:=s.run(ctx,tenant,actor,key,"BOOK_CREATE","ACCOUNTING_BOOK",input,func(ctx context.Context, existing string)(string,time.Time,error){
  var err error
  if existing!="" { result,err=s.books.FindByID(ctx,tenant,existing) } else { result,err=s.books.Create(ctx,tenant,input) }
  return result.BookID,result.CreatedAt,err
 })
 if err!=nil { return entity.AccountingBook{},err }
 return result,nil
}

func (s *Service) run(ctx context.Context, tenant, actor, key, operation, assetType string, input any, create func(context.Context,string)(string,time.Time,error)) error {
 if tenant=="" || strings.TrimSpace(tenant)!=tenant || len(tenant)>128 || key=="" || strings.TrimSpace(key)!=key || len(key)>128 { return errors.New("invalid creation scope or idempotency key") }
 if err:=validators.ValidateID(actor);err!=nil{return err}
 raw,err:=json.Marshal([]any{actor,input});if err!=nil{return err}
 digest:=sha256.Sum256(raw)
 hash:=hex.EncodeToString(digest[:])
 return s.unit.Run(ctx,func(ctx context.Context)error{
  db:=transaction.Database(ctx,s.db)
  if err:=db.Exec("INSERT INTO accounting.creation_receipts (tenant_id,operation,idempotency_key,request_hash) VALUES (?,?,?,?) ON CONFLICT DO NOTHING",tenant,operation,key,hash).Error;err!=nil{return err}
  var receipt struct { RequestHash string; AssetID *string }
  if err:=db.Raw("SELECT request_hash,asset_id FROM accounting.creation_receipts WHERE tenant_id=? AND operation=? AND idempotency_key=? FOR UPDATE",tenant,operation,key).Scan(&receipt).Error;err!=nil{return err}
  if receipt.RequestHash!=hash{return ErrConflict}
  if receipt.AssetID!=nil { _,_,err:=create(ctx,*receipt.AssetID);return err }
  id,occurred,err:=create(ctx,"");if err!=nil{return err}
  if err:=s.audit.Record(ctx,producer.Activity{TenantID:tenant,Source:"accounting",AssetType:assetType,AssetID:id,ActorID:actor,AssetVersion:1,Operation:"CREATED",OccurredAt:occurred});err!=nil{return err}
  return db.Exec("UPDATE accounting.creation_receipts SET asset_id=? WHERE tenant_id=? AND operation=? AND idempotency_key=?",id,tenant,operation,key).Error
 })
}
