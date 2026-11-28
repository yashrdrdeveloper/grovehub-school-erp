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
        // Run immediately after restart, then make bounded progress every
        // five minutes. The cursor advances past disabled/failing schools.
        for ctx.Err()==nil {
            renewalCtx,stopRenewal:=context.WithTimeout(ctx,4*time.Minute)
            if err:=service.RenewExpiring(renewalCtx);err!=nil&&ctx.Err()==nil {b.logger.Warn("integration watch renewal incomplete")}
            stopRenewal()
            timer:=time.NewTimer(5*time.Minute)
            select {case <-ctx.Done():timer.Stop();return;case <-timer.C:}
        }
    }()
    return cancel,done
}
