package impl

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"gorm.io/gorm"
	"grovehub/libs/schoolgateway"
	inframessaging "grovehub/libs/infrastructure/messaging"
	"grovehub/services/audit/internal/app/bootstrap"
	"grovehub/services/audit/internal/app/constants/events"
	protected "grovehub/services/audit/internal/app/controller/activity/protected"
 operator "grovehub/services/audit/internal/app/controller/operator/restricted"
	controllerimpl "grovehub/services/audit/internal/app/controller/activity/protected/impl"
	consumer "grovehub/services/audit/internal/app/events/consumers/activity"
	repository "grovehub/services/audit/internal/app/repository/activity/impl"
	rejection "grovehub/services/audit/internal/app/repository/rejection/impl"
	"grovehub/services/audit/internal/app/routes"
	"grovehub/services/audit/internal/app/services/monitoring"
	service "grovehub/services/audit/internal/app/services/activity/impl"
	"grovehub/services/infrastructure/client/messaging"
)

type builder struct {
	logger *slog.Logger
	db *gorm.DB
	messaging messaging.MessagingClient
	access schoolgateway.ProtectedRouteAuthorizer
}

func New(logger *slog.Logger, db *gorm.DB, client messaging.MessagingClient, access schoolgateway.ProtectedRouteAuthorizer) bootstrap.Bootstrap {
	return &builder{logger:logger,db:db,messaging:client,access:access}
}

type runtime struct {
	controller protected.ActivityController
 operator http.HandlerFunc
 operatorSearch http.HandlerFunc
	stop context.CancelFunc
	wait sync.WaitGroup
	once sync.Once
	monitor *monitoring.Monitor
}

func (b *builder) Build() (bootstrap.Runtime,error) {
	if b.db==nil || b.messaging==nil || b.access==nil { return nil,errors.New("audit requires PostgreSQL, Infrastructure and access policy") }
	if b.logger==nil { b.logger=slog.Default() }
	if err:=b.db.Exec("SELECT subject_account_id, occurrence_id FROM audit.activities LIMIT 0").Error; err!=nil { return nil,errors.New("audit migrations required") }
	if err:=b.db.Exec("SELECT fingerprint FROM audit.rejected_events LIMIT 0").Error; err!=nil { return nil,errors.New("audit rejection migrations required") }
	if err:=b.db.Exec("SELECT source_event_id FROM audit.event_receipts LIMIT 0").Error;err!=nil { return nil,errors.New("audit receipt migrations required") }
	if err:=b.db.Exec("SELECT event_id,correlation_id FROM audit.operator_activities LIMIT 0").Error;err!=nil{return nil,errors.New("operator audit migration required")}
 repo:=repository.New(b.db)
	handler:=consumer.New(repo,rejection.New(b.db))
	monitor,err:=monitoring.New(events.Topics())
	if err!=nil { return nil,errors.New("audit monitoring initialization failed") }
	ctx,stop:=context.WithCancel(context.Background())
	module:=&runtime{controller:controllerimpl.New(service.New(repo,b.access)),operator:operator.Handler(b.db),operatorSearch:operator.SearchHandler(b.db),stop:stop,monitor:monitor}
	module.wait.Add(1)
	go module.probe(ctx,b.db,b.messaging)
	for topic,eventType:=range events.Topics() {
		module.wait.Add(1)
		go module.consume(ctx,b.logger,b.messaging,handler,topic,eventType)
	}
	return module,nil
}

func (module *runtime) consume(ctx context.Context, logger *slog.Logger, client messaging.MessagingClient, handler *consumer.Consumer, topic string, eventType string) {
	defer module.wait.Done()
	for ctx.Err()==nil {
		subscription:=inframessaging.WithSubscriptionObserver(inframessaging.WithConsumerGroup(ctx,"audit"),func(ready bool){module.monitor.Subscription(topic,ready)})
		_ = client.Consume(subscription,topic,func(deliveryCtx context.Context,message messaging.Message) error {
			started:=time.Now()
			outcome,err:=handler.Process(deliveryCtx,eventType,message)
			module.monitor.Delivery(deliveryCtx,topic,outcome,time.Since(started),err)
			if err!=nil {
				logger.Warn("audit event not persisted","topic",topic)
				return err
			}
			return nil
		})
 		module.monitor.Subscription(topic,false)
		if ctx.Err()!=nil { return }
		module.monitor.Interrupted(ctx,topic)
		logger.Warn("audit subscription interrupted","topic",topic)
		select { case <-ctx.Done(): return; case <-time.After(time.Second): }
	}
}

func (module *runtime) RegisterHTTP(mux *http.ServeMux) { routes.RegisterHTTP(mux,module.controller); mux.HandleFunc(operator.Pattern,module.operator); mux.HandleFunc(operator.SearchPattern,module.operatorSearch) }
func (module *runtime) Close() { module.once.Do(func(){module.monitor.Close();module.stop();module.wait.Wait()}) }
