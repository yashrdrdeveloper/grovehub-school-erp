package impl

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	"gorm.io/gorm"
	schoolgateway "grovehub/libs/platform/session"

	"grovehub/services/audit/internal/app/bootstrap"
	"grovehub/services/audit/internal/app/bootstrap/constants"
	"grovehub/services/audit/internal/app/bootstrap/helpers"
	"grovehub/services/audit/internal/app/constants/events"
	protected "grovehub/services/audit/internal/app/controller/activity/protected"
	controllerimpl "grovehub/services/audit/internal/app/controller/activity/protected/impl"
	operator "grovehub/services/audit/internal/app/controller/operator/restricted"
	consumer "grovehub/services/audit/internal/app/events/consumers/activity"
	repository "grovehub/services/audit/internal/app/repository/activity/impl"
	rejection "grovehub/services/audit/internal/app/repository/rejection/impl"
	schemarepo "grovehub/services/audit/internal/app/repository/schema/impl"
	"grovehub/services/audit/internal/app/routes"
	service "grovehub/services/audit/internal/app/services/activity/impl"
	"grovehub/services/audit/internal/app/services/monitoring"
	"grovehub/services/infrastructure/client/messaging"
)

type builder struct {
	logger    *slog.Logger
	db        *gorm.DB
	messaging messaging.MessagingClient
	access    schoolgateway.ProtectedRouteAuthorizer
}

func New(logger *slog.Logger, db *gorm.DB, client messaging.MessagingClient, access schoolgateway.ProtectedRouteAuthorizer) bootstrap.Bootstrap {
	return &builder{logger: logger, db: db, messaging: client, access: access}
}

type runtime struct {
	controller          protected.ActivityController
	operator            http.HandlerFunc
	operatorSearch      http.HandlerFunc
	schoolCheck         http.HandlerFunc
	schoolReceipt       http.HandlerFunc
	schoolStatusReceipt http.HandlerFunc
	schoolReceiptSearch http.HandlerFunc
	stop                context.CancelFunc
	wait                sync.WaitGroup
	once                sync.Once
	monitor             *monitoring.Monitor
}

func (b *builder) Build() (bootstrap.Runtime, error) {
	if b.db == nil || b.messaging == nil || b.access == nil {
		return nil, coreerrors.NewError(constants.DependenciesRequired)
	}

	if b.logger == nil {
		b.logger = slog.Default()
	}

	if err := helpers.CheckSchema(schemarepo.New(b.db)); err != nil {
		return nil, err
	}

	repo := repository.New(b.db)
	handler := consumer.New(repo, rejection.New(b.db))
	topics := events.Topics()
	monitor, err := monitoring.New(topics)
	if err != nil {
		return nil, coreerrors.Wrap(constants.MonitoringFailed, err)
	}

	ctx, stop := context.WithCancel(context.Background())
	module := &runtime{
		controller:          controllerimpl.New(service.New(repo, b.access)),
		operator:            operator.Handler(b.db),
		operatorSearch:      operator.SearchHandler(b.db),
		schoolCheck:         operator.SchoolCheckHandler(b.db),
		schoolReceipt:       operator.ReceiptHandler(b.db),
		schoolStatusReceipt: operator.StatusReceiptHandler(b.db),
		schoolReceiptSearch: operator.ReceiptSearchHandler(b.db),
		stop:                stop,
		monitor:             monitor,
	}

	module.wait.Add(1)
	go module.probe(ctx, b.db, b.messaging)

	for topic, eventType := range topics {
		module.wait.Add(1)
		go module.consume(ctx, b.logger, b.messaging, handler, topic, eventType)
	}

	return module, nil
}

func (module *runtime) consume(ctx context.Context, logger *slog.Logger, client messaging.MessagingClient, handler *consumer.Consumer, topic string, eventType string) {
	defer module.wait.Done()

	for ctx.Err() == nil {
		subscription := messaging.WithSubscriptionObserver(
			messaging.WithConsumerGroup(ctx, constants.ServiceName),
			func(ready bool) { module.monitor.Subscription(topic, ready) },
		)

		_ = client.Consume(subscription, topic, func(deliveryCtx context.Context, message messaging.Message) error {
			started := time.Now()
			outcome, err := handler.Process(deliveryCtx, eventType, message)
			module.monitor.Delivery(deliveryCtx, topic, outcome, time.Since(started), err)

			if err != nil {
				logger.Warn("audit event not persisted", "topic", topic)
				return err
			}

			return nil
		})

		module.monitor.Subscription(topic, false)
		if ctx.Err() != nil {
			return
		}

		module.monitor.Interrupted(ctx, topic)
		logger.Warn("audit subscription interrupted", "topic", topic)

		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Second):
		}
	}
}

func (module *runtime) RegisterHTTP(mux *http.ServeMux) {
	routes.RegisterHTTP(mux, module.controller)
	mux.HandleFunc(operator.Pattern, module.operator)
	mux.HandleFunc(operator.SearchPattern, module.operatorSearch)
	mux.HandleFunc(operator.SchoolCheckPattern, module.schoolCheck)
	mux.HandleFunc(operator.ReceiptPattern, module.schoolReceipt)
	mux.HandleFunc(operator.StatusReceiptPattern, module.schoolStatusReceipt)
	mux.HandleFunc(operator.ReceiptSearchPattern, module.schoolReceiptSearch)
}

func (module *runtime) Close() {
	module.once.Do(func() {
		module.monitor.Close()
		module.stop()
		module.wait.Wait()
	})
}
