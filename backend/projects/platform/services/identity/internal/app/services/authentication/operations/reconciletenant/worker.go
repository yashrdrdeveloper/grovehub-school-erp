package reconciletenant

import (
 "context"
 "errors"
 "log/slog"
 "time"

 schoolgateway "grovehub/libs/platform/session"
)

const (
 batchSize = 20
 interval = 15 * time.Second
 operationTimeout = 5 * time.Second
)

type Store interface {
 ListTenants(context.Context, string, int) ([]string, error)
 RevokeBefore(context.Context, string, int64) (int64, error)
}

type Worker struct {
 store Store
 lookup schoolgateway.TenantStateLookup
 logger *slog.Logger
 cursor string
}

func New(store Store, lookup schoolgateway.TenantStateLookup, logger *slog.Logger) *Worker {
 if logger == nil { logger = slog.Default() }
 return &Worker{store: store, lookup: lookup, logger: logger}
}

// Step is called serially. Failed tenants are retried on the next sweep,
// while unrelated tenants can continue. The database remains the durable
// source of pending work; restarting safely begins another sweep.
func (w *Worker) Step(ctx context.Context) error {
 probe, cancel := context.WithTimeout(ctx, operationTimeout)
 tenants, err := w.store.ListTenants(probe, w.cursor, batchSize)
 cancel()
 if err != nil { return err }
 failed := false
 for _, tenant := range tenants {
  if ctx.Err() != nil { return ctx.Err() }
  probe, cancel := context.WithTimeout(ctx, operationTimeout)
  state, err := w.lookup(probe, tenant)
  cancel()
  if err == nil && state.Generation > 0 {
   write, cancel := context.WithTimeout(ctx, operationTimeout)
   _, err = w.store.RevokeBefore(write, tenant, state.Generation)
   cancel()
  } else {
   err = errors.New("tenant generation unavailable")
  }
  if err != nil { failed = true }
  w.cursor = tenant
 }
 if len(tenants) < batchSize { w.cursor = "" }
 if failed { return errors.New("tenant reconciliation incomplete") }
 return nil
}

func (w *Worker) Run(ctx context.Context) {
 for ctx.Err() == nil {
  if err := w.Step(ctx); err != nil && ctx.Err() == nil {
   w.logger.Warn("Identity tenant session reconciliation incomplete; retry on next sweep")
  }
  timer := time.NewTimer(interval)
  select {
  case <-ctx.Done(): timer.Stop(); return
  case <-timer.C:
  }
 }
}
