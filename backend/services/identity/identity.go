// Package identity exposes only the module entrypoint; implementation is private.
package identity
import(
 "context"
 "errors"
 "fmt"
 "log/slog"
 "gorm.io/gorm"
 "grovehub/services/identity/internal/app/bootstrap"
 bootstrapimpl "grovehub/services/identity/internal/app/bootstrap/impl"
 migrations "grovehub/services/identity/internal/resources/migrations"
 foundationclient "grovehub/services/foundation/client/person"
 integrationclient "grovehub/services/integration/client/connection"
)
type Module=bootstrap.Runtime
func New(logger *slog.Logger,db *gorm.DB,foundation foundationclient.Client,integration integrationclient.Client)(Module,error){
 return bootstrapimpl.New(logger,db,foundation,integration).Build()
}
func Migrate(ctx context.Context,db *gorm.DB)error{
 if db==nil{return errors.New("identity migration requires database")}
 sqlDB,err:=db.DB();if err!=nil{return fmt.Errorf("identity migration DB: %w",err)}
 return migrations.Migrate(ctx,sqlDB)
}
