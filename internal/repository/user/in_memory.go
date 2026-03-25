package user

import (
	"context"
	"sync"

	"OrderService/internal/model"

	errs "OrderService/internal/errors"

	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// тут трейсы как будто бы тоже закинуть нужно но это нигде не используется так то
type Repo struct {
	Users  map[uuid.UUID]model.User
	mu     *sync.RWMutex
	log    log.Logger
	tracer trace.Tracer
}

func NewRepo(logger log.Logger, tp trace.TracerProvider) *Repo {
	repo := &Repo{
		Users:  make(map[uuid.UUID]model.User),
		mu:     &sync.RWMutex{},
		log:    logger,
		tracer: tp.Tracer("order-service/UserRepo"),
	}
	users := []model.User{
		{
			ID:    uuid.MustParse("1179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Gleb",
			Roles: []string{"TRADER"},
		},
		{
			ID:    uuid.MustParse("2179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Oleg",
			Roles: []string{"ADMIN"},
		},
		{
			ID:    uuid.MustParse("3179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Vova",
			Roles: []string{"TRADER"},
		},
		{
			ID:    uuid.MustParse("4179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Arsen",
			Roles: []string{"ADMIN"},
		},
	}
	for _, user := range users {
		repo.Users[user.ID] = user
	}
	return repo
}

const layer = "UserInMemoryRepo"

func (r *Repo) CreateUser(ctx context.Context, user model.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.log.Debug(layer, "CreateUser", "user success created")

	r.Users[user.ID] = user
}

func (r *Repo) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *errors.CustomError) {
	const method = "GetUserById"

	ctx, span := r.tracer.Start(ctx, "UserRepo.GetUserById")
	defer span.End()

	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, found := r.Users[id]; found {
		r.log.Debug(
			layer,
			method,
			"found user",
			"user_id", id,
		)
		return &u, nil
	}

	span.RecordError(errs.ErrUserNotFound)
	span.SetStatus(codes.Error, errs.ErrUserNotFound.Error())

	r.log.Error(layer, method, "user not found", errs.ErrUserNotFound)

	return nil, errs.ErrUserNotFound
}
