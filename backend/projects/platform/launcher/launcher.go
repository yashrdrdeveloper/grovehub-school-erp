package launcher

import (
 "context"
 "io"
 "log/slog"
 "gorm.io/gorm"
 bootstrap "grovehub/launcher/internal/src/bootstrap/impl"
)

func Run(logger *slog.Logger) error {
 runtime,err:=bootstrap.New(logger).Build()
 if err!=nil{return err}
 defer runtime.Close()
 return runtime.Run()
}

// Operator commands share the same Platform-owned database and logger setup.
func NewLogger(output io.Writer)(*slog.Logger,error){return bootstrap.NewLogger(output)}
func NewPostgresDB(ctx context.Context)(*gorm.DB,error){return bootstrap.NewPostgresDB(ctx)}
func NewAuditPostgresDB(ctx context.Context,shared *gorm.DB)(*gorm.DB,func(),error){return bootstrap.NewAuditPostgresDB(ctx,shared)}
var ErrDatabaseURLRequired=bootstrap.ErrDatabaseURLRequired
