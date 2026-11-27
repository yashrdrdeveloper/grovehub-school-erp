package migrations
import("context";"database/sql";"errors";"fmt";"github.com/pressly/goose/v3";resources "grovehub/services/identity/internal/resources")
func Migrate(ctx context.Context,db *sql.DB)error{
 if db==nil{return errors.New("identity requires a database")}
 if err:=goose.SetDialect("postgres");err!=nil{return fmt.Errorf("set Identity dialect: %w",err)}
 goose.SetBaseFS(resources.Migrations);defer goose.SetBaseFS(nil)
 goose.SetTableName("identity_schema_state")
 return goose.UpContext(ctx,db,"migrations")
}
