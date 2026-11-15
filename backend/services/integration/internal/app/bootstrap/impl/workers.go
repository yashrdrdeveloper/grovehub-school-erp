package impl

import (
    "context"
    "time"

    connections "grovehub/services/integration/internal/app/services/connection"
)

// Renewal follows module lifecycle cancellation; in-flight work can be canceled
// instead of delaying shutdown until its five-minute deadline.
func (b *builder) newWatchRenewal(service connections.ConnectionService, enabled bool) (context.CancelFunc,chan struct{}) {
    ctx,cancel:=context.WithCancel(context.Background())
    done:=make(chan struct{})
    if !enabled {close(done);return cancel,done}
    go func(){
        defer close(done)
        ticker:=time.NewTicker(6*time.Hour)
        defer ticker.Stop()
        for {
            select {
            case <-ctx.Done():return
            case <-ticker.C:
                renewalCtx,stopRenewal:=context.WithTimeout(ctx,5*time.Minute)
                if err:=service.RenewExpiring(renewalCtx);err!=nil&&ctx.Err()==nil {b.logger.Warn("integration watch renewal failed","error",err)}
                stopRenewal()
            }
        }
    }()
    return cancel,done
}
