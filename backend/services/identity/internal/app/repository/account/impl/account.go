package impl
import (
 "context"
 "errors"
 "strings"
 "gorm.io/gorm"
 "github.com/jackc/pgx/v5/pgconn"
 coretransaction "grovehub/libs/core/database/transaction"
 coresearch "grovehub/libs/core/search"
 "grovehub/services/identity/internal/app/entity"
 accountrepository "grovehub/services/identity/internal/app/repository/account"
 searchconfig "grovehub/services/identity/internal/app/repository/account/helpers/search"
)
type Repository struct{db *gorm.DB}
func New(db *gorm.DB)accountrepository.Repository{return &Repository{db:db}}
func (r *Repository) Create(ctx context.Context,a entity.Account)error{return mapWriteError(r.dbFor(ctx).Create(&a).Error)}
func (r *Repository) dbFor(ctx context.Context)*gorm.DB{return coretransaction.Database(ctx,r.db)}
func (r *Repository) find(ctx context.Context,field string,value any)(entity.Account,error){
 var a entity.Account
 err:=r.dbFor(ctx).Where(field+" = ?",value).Take(&a).Error
 if errors.Is(err,gorm.ErrRecordNotFound){return a,accountrepository.ErrNotFound}
 return a,err
}
func (r *Repository) FindByID(ctx context.Context,id string)(entity.Account,error){return r.find(ctx,"account_id",id)}
func (r *Repository) FindByEmail(ctx context.Context,email string)(entity.Account,error){return r.find(ctx,"email",strings.ToLower(email))}
func (r *Repository) FindByUsername(ctx context.Context,name string)(entity.Account,error){return r.find(ctx,"lower(username)",strings.ToLower(name))}
func (r *Repository) FindByPersonGID(ctx context.Context,gid string)(entity.Account,error){return r.find(ctx,"person_gid",gid)}
func (r *Repository) FindByIdempotencyKey(ctx context.Context,key string)(entity.Account,error){return r.find(ctx,"registration_idempotency_key",key)}
func (r *Repository) Update(ctx context.Context,a entity.Account)error{
 result:=r.dbFor(ctx).Model(&entity.Account{}).Where("account_id = ? AND version = ?",a.AccountID,a.Version-1).
 Updates(map[string]any{"username":a.Username,"email":a.Email,"status":a.Status,"version":a.Version,"search_details":a.SearchDetails,"updated_at":a.UpdatedAt})
 if result.Error!=nil{return mapWriteError(result.Error)};if result.RowsAffected!=1{return accountrepository.ErrVersionConflict};return nil
}
// Search uses the trusted AccountConfig with a mandatory tenant predicate.
func (r *Repository) Search(ctx context.Context,tenantID string,q coresearch.Query[string])(coresearch.QueryResult[entity.Account],error){
 if tenantID==""{return coresearch.QueryResult[entity.Account]{},errors.New("account search requires tenant scope")}
 return coresearch.Execute[string,entity.Account](r.dbFor(ctx).Model(&entity.Account{}).Table("identity.accounts AS a").Where("a.tenant_id = ?",tenantID),q,searchconfig.AccountConfig)
}

func mapWriteError(err error)error{
 if err==nil{return nil}
 var pg *pgconn.PgError
 if errors.As(err,&pg)&&pg.Code=="23505"{return accountrepository.ErrConflict}
 return err
}
