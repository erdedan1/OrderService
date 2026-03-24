package order

import (
	"OrderService/internal/usecase"

	log "github.com/erdedan1/shared/logger"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	orderRepo             usecase.OrderRepo
	userRepo              usecase.UserRepo
	marketCache           usecase.MarketCacheRepo
	marketSrv             usecase.MarketService
	orderStatusSubscriber usecase.OrderStatusSubscriber
	orderStatusPublisher  usecase.OrderStatusPublisher
	log                   log.Logger
	tracer                trace.Tracer
}

func New(
	repo usecase.OrderRepo,
	userRepo usecase.UserRepo,
	marketCache usecase.MarketCacheRepo,
	marketSrv usecase.MarketService,
	orderStatusSubscriber usecase.OrderStatusSubscriber,
	orderStatusPublisher usecase.OrderStatusPublisher,
	log log.Logger,
	tp trace.TracerProvider,
) *Service {
	return &Service{
		orderRepo:             repo,
		userRepo:              userRepo,
		marketCache:           marketCache,
		marketSrv:             marketSrv,
		orderStatusSubscriber: orderStatusSubscriber,
		orderStatusPublisher:  orderStatusPublisher,
		log:                   log,
		tracer:                tp.Tracer("order-service/Service"),
	}
}

const layer = "OrderService"
