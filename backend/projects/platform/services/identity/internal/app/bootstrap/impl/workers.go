package impl

import (
	"context"
	"time"

	invitation "grovehub/services/identity/internal/app/services/registration/operations/invitation"
)

const (
	invitationPollInterval = 10 * time.Second
	invitationRunTimeout   = 20 * time.Second
)

func (module *runtime) startWorkers(worker *invitation.Worker) {
	if worker == nil {
		return
	}
	ctx, stop := context.WithCancel(context.Background())
	module.stop = stop
	module.done = make(chan struct{})
	go module.runWorker(ctx, worker)
}

func (module *runtime) runWorker(ctx context.Context, worker *invitation.Worker) {
	defer close(module.done)

	ticker := time.NewTicker(invitationPollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			runCtx, cancel := context.WithTimeout(ctx, invitationRunTimeout)
			_, err := worker.RunOnce(runCtx)
			cancel()
			if err != nil {
				module.logger.Warn("Identity invitation delivery failed; retry scheduled")
			}
		}
	}
}
