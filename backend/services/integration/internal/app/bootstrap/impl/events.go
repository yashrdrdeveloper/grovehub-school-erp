package impl

import (
    "context"
    "time"

    services "grovehub/services/integration/internal/app/services"
)

// Renewal belongs to the module lifecycle, not the public module constructor.
func (b *builder) newEvents(service services.ConnectionService, enabled bool) (chan struct{},chan struct{}) {
    stop:=make(chan struct{})
    done:=make(chan struct{})
    if !enabled {close(done);return stop,done}
    go func(){
        defer close(done)
        ticker:=time.NewTicker(6*time.Hour)
        defer ticker.Stop()
        for {
            select {
            case <-stop:return
            case <-ticker.C:
                ctx,cancel:=context.WithTimeout(context.Background(),5*time.Minute)
                if err:=service.RenewExpiring(ctx);err!=nil {b.logger.Warn("integration watch renewal failed","error",err)}
                cancel()
            }
        }
    }()
    return stop,done
}
