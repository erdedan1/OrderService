package order

import (
	"context"
	"testing"
	"time"

	"OrderService/internal/dto"
	"OrderService/internal/errors"
	"OrderService/internal/model"
	"OrderService/mocks"

	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/otel/trace/noop"
)

func preparingTests(t *testing.T) (
	*Service,
	*mocks.OrderRepo,
	*mocks.UserRepo,
	*mocks.MarketCacheRepo,
	*mocks.MarketService,
	*mocks.OrderStatusSubscriber,
	*mocks.OrderStatusPublisher,

) {
	orderRepo := mocks.NewOrderRepo(t)
	userRepo := mocks.NewUserRepo(t)
	cache := mocks.NewMarketCacheRepo(t)
	marketSrv := mocks.NewMarketService(t)
	subscriber := mocks.NewOrderStatusSubscriber(t)
	publisher := mocks.NewOrderStatusPublisher(t)

	logger, _ := log.NewLogger("debug")
	defer logger.Sync()

	service := New(
		orderRepo,
		userRepo,
		cache,
		marketSrv,
		subscriber,
		publisher,
		logger,
		noop.NewTracerProvider(),
	)
	return service, orderRepo, userRepo, cache, marketSrv, subscriber, publisher
}
func TestCreateOrder_Success(t *testing.T) {
	service, orderRepo, userRepo, cache, marketSrv, _, publisher := preparingTests(t)
	ctx := context.Background()

	userID := uuid.MustParse("1179803e-06f0-4369-b94f-14e26ec190a3")
	user := &model.User{ID: userID, Roles: []string{"USER_ROLE_TRADER"}}
	order := &model.Order{ID: uuid.New(), Status: model.StatusCreated, UserID: userID}

	userRepo.On("GetUserById", mock.Anything, userID).
		Return(user, nil)
	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil, nil)
	marketSrv.On("ViewMarketsByRoles", mock.Anything, mock.Anything).
		Return([]dto.ViewMarketsResponse{{ID: uuid.New()}}, nil)
	cache.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	publisher.On("PublishOrderStatus", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	orderRepo.On("CreateOrder", mock.Anything, mock.Anything).
		Return(order, nil)

	res, err := service.CreateOrder(ctx, &dto.CreateOrderRequest{
		MarketID:  uuid.New(),
		UserID:    userID,
		OrderType: "Test_type",
		Price:     120,
		UserRoles: user.Roles,
		Quantity:  1,
	})

	assert.Nil(t, err)
	assert.Equal(t, order.ID, res.ID)
	assert.Equal(t, order.Status.ToString(), res.Status)

	userRepo.AssertExpectations(t)
	cache.AssertExpectations(t)
	marketSrv.AssertExpectations(t)
	publisher.AssertExpectations(t)
	orderRepo.AssertExpectations(t)
}

func TestCreateOrder_User_No_Acess(t *testing.T) {
	service, _, userRepo, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	userID := uuid.MustParse("1179803e-06f0-4369-b94f-14e26ec190a3")
	user := &model.User{ID: userID, Roles: []string{"USER_ROLE_TRADER"}}

	userRepo.On("GetUserById", mock.Anything, userID).
		Return(user, nil)

	res, err := service.CreateOrder(ctx, &dto.CreateOrderRequest{
		MarketID:  uuid.New(),
		UserID:    userID,
		OrderType: "Test_type",
		Price:     120,
		UserRoles: []string{"TEST_ROLE"},
		Quantity:  1,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, errors.ErrUserHasNoAccessToMarket.Message, err.Message)

	userRepo.AssertExpectations(t)
}

func TestCreateOrder_UserRepo_Error(t *testing.T) {
	service, _, userRepo, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	userID := uuid.New()

	userRepo.On("GetUserById", mock.Anything, userID).
		Return(nil, errors.ErrUserNotFound)

	res, err := service.CreateOrder(ctx, &dto.CreateOrderRequest{
		MarketID:  uuid.New(),
		UserID:    userID,
		OrderType: "Test_type",
		Price:     120,
		UserRoles: []string{"USER_ROLE_TRADER"},
		Quantity:  1,
	})

	assert.Nil(t, res)
	assert.Error(t, err)

	userRepo.AssertExpectations(t)
}

func TestCreateOrder_Market_Not_Found(t *testing.T) {
	service, _, userRepo, cache, marketSrv, _, _ := preparingTests(t)
	ctx := context.Background()

	userID := uuid.MustParse("1179803e-06f0-4369-b94f-14e26ec190a3")
	userRoles := []string{"USER_ROLE_TRADER"}
	user := &model.User{ID: userID, Roles: []string{"USER_ROLE_TRADER"}}

	userRepo.On("GetUserById", mock.Anything, userID).
		Return(user, nil)
	cache.On("Get", mock.Anything, mock.Anything).
		Return(nil, nil)
	marketSrv.On("ViewMarketsByRoles", mock.Anything, mock.Anything).
		Return(nil, errors.ErrMarketNotFound)

	res, err := service.CreateOrder(ctx, &dto.CreateOrderRequest{
		MarketID:  uuid.New(),
		UserID:    userID,
		OrderType: "Test_type",
		Price:     120,
		UserRoles: userRoles,
		Quantity:  1,
	})

	assert.Nil(t, res)
	assert.Error(t, err)

	userRepo.AssertExpectations(t)
	cache.AssertExpectations(t)
	marketSrv.AssertExpectations(t)
}
func TestGetOrderStatus_Success(t *testing.T) {
	service, orderRepo, _, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	userID := uuid.New()
	orderID := uuid.New()

	order := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.StatusCreated,
	}

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(order, nil)

	res, err := service.GetOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, order.Status.ToString(), res.Status)
	assert.NotNil(t, res.UpdatedAt)

	orderRepo.AssertExpectations(t)
}

func TestGetOrderStatus_OrderRepo_Error(t *testing.T) {
	service, orderRepo, _, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	userID := uuid.New()
	orderID := uuid.New()

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(nil, errors.ErrOrderNotFound)

	res, err := service.GetOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.Nil(t, res)
	assert.Error(t, err)

	orderRepo.AssertExpectations(t)
}

func TestGetOrderStatus_InvalidUserID(t *testing.T) {
	service, orderRepo, _, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	userID := uuid.New()
	anotherUserID := uuid.New()
	orderID := uuid.New()

	order := &model.Order{
		ID:     orderID,
		UserID: anotherUserID,
		Status: model.StatusCreated,
	}

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(order, nil)

	res, err := service.GetOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.Nil(t, res)
	assert.Error(t, err)
	assert.Equal(t, errors.ErrInvalidUserID.Message, err.Message)

	orderRepo.AssertExpectations(t)
}

func TestSubscribeOrderStatus_GetOrder_Error(t *testing.T) {
	service, orderRepo, _, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(nil, errors.ErrOrderNotFound)

	ch, err := service.SubscribeOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.Nil(t, ch)
	assert.NotNil(t, err)

	orderRepo.AssertExpectations(t)
}

func TestSubscribeOrderStatus_InvalidUserID(t *testing.T) {
	service, orderRepo, _, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	order := &model.Order{
		ID:     orderID,
		UserID: uuid.New(),
		Status: model.StatusCreated,
	}

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(order, nil)

	ch, err := service.SubscribeOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.NotNil(t, ch)
	assert.Error(t, err)
	assert.Equal(t, errors.ErrInvalidUserID.Message, err.Message)

	orderRepo.AssertExpectations(t)
}

func TestSubscribeOrderStatus_OrderClosed(t *testing.T) {
	service, orderRepo, _, _, _, _, _ := preparingTests(t)
	ctx := context.Background()

	orderID := uuid.New()
	userID := uuid.New()

	now := time.Now()

	order := &model.Order{
		ID:        orderID,
		UserID:    userID,
		Status:    model.StatusClosed,
		UpdatedAt: &now,
	}

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(order, nil)

	ch, err := service.SubscribeOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.Nil(t, err)

	res := <-ch

	assert.Equal(t, model.StatusClosed.ToString(), res.Status)

	orderRepo.AssertExpectations(t)
}

func TestSubscribeOrderStatus_Subscribe_Error(t *testing.T) {
	service, orderRepo, _, _, _, subscriber, _ := preparingTests(t)

	ctx := context.Background()
	orderID := uuid.New()
	userID := uuid.New()

	order := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.StatusCreated,
	}

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(order, nil)

	subscriber.On("SubscribeOrderStatus", mock.Anything, orderID).
		Return(nil, errors.ErrUnavailableRedis)

	ch, err := service.SubscribeOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.Nil(t, ch)
	assert.NotNil(t, err)

	orderRepo.AssertExpectations(t)
	subscriber.AssertExpectations(t)
}

func TestSubscribeOrderStatus_Success(t *testing.T) {
	service, orderRepo, _, _, _, subscriber, _ := preparingTests(t)

	ctx := context.Background()
	orderID := uuid.New()
	userID := uuid.New()

	statusCh := make(chan model.OrderStatus, 2)

	order := &model.Order{
		ID:     orderID,
		UserID: userID,
		Status: model.StatusCreated,
	}

	orderRepo.On("GetOrder", mock.Anything, orderID).
		Return(order, nil)

	subscriber.On("SubscribeOrderStatus", mock.Anything, orderID).
		Return((<-chan model.OrderStatus)(statusCh), nil)

	ch, err := service.SubscribeOrderStatus(ctx, &dto.GetOrderStatusRequest{
		UserID:  userID,
		OrderID: orderID,
	})

	assert.Nil(t, err)

	first := <-ch
	assert.Equal(t, model.StatusCreated.ToString(), first.Status)

	statusCh <- model.StatusClosed

	second := <-ch

	assert.Equal(t, model.StatusClosed.ToString(), second.Status)
	orderRepo.AssertNotCalled(t, "UpdateOrderStatus", mock.Anything, mock.Anything, mock.Anything)

	orderRepo.AssertExpectations(t)
	subscriber.AssertExpectations(t)
}
