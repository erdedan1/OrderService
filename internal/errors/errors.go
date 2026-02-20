package errors

import (
	"github.com/erdedan1/shared/errs"
)

var (
	ErrOrderNotFound           = errs.New(errs.NOT_FOUND, "order not found")
	ErrIvalidArgument          = errs.New(errs.INVALID_ARGUMENT, "invalid argument")
	ErrUserHasNoAccessToMarket = errs.New(errs.PERMISSION_DENIED, "user has no acces to market")
	ErrNoMarketsAvailable      = errs.New(errs.NOT_FOUND, "no markets available for user")
	ErrInvalidUserID           = errs.New(errs.PERMISSION_DENIED, "invalid user id")

	ErrUserNotFound = errs.New(errs.NOT_FOUND, "user not found")

	ErrMarketNotFound = errs.New(errs.NOT_FOUND, "market not found")

	ErrFailedSerializeRedis   = errs.New(errs.INTERNAL, "failed to serialize markets redis")
	ErrUnavailableRedis       = errs.New(errs.UNAVAILABLE, "redis is unavailable")
	ErrUnavailableDataRedis   = errs.New(errs.UNAVAILABLE, "failed to get data from redis")
	ErrFailedDeserializeRedis = errs.New(errs.INTERNAL, "failed to deserialize markets redis")
	ErrDeleteRedis            = errs.New(errs.UNAVAILABLE, "failed to delete redis key")
)
