package user

import (
	"context"
	"sync"

	"OrderService/internal/model"

	errs "OrderService/internal/errors"

	errors "github.com/erdedan1/shared/errs"
	log "github.com/erdedan1/shared/logger"
	"github.com/google/uuid"
)

type Repo struct {
	Users map[uuid.UUID]model.User
	mu    *sync.RWMutex
	l     log.Logger
}

func NewRepo(logger log.Logger) *Repo {
	repo := &Repo{
		Users: make(map[uuid.UUID]model.User),
		mu:    &sync.RWMutex{},
		l:     logger.Layer("User.Repository"),
	}
	users := []model.User{
		{
			ID:    uuid.MustParse("1179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Gleb",
			Roles: []string{"USER_ROLE_TRADER"},
		},
		{
			ID:    uuid.MustParse("2179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Oleg",
			Roles: []string{"USER_ROLE_ADMIN"},
		},
		{
			ID:    uuid.MustParse("3179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Vova",
			Roles: []string{"USER_ROLE_TRADER"},
		},
		{
			ID:    uuid.MustParse("4179803e-06f0-4369-b94f-14e26ec190a3"),
			Name:  "Arsen",
			Roles: []string{"USER_ROLE_ADMIN"},
		},
	}
	for _, user := range users {
		repo.Users[user.ID] = user
	}
	return repo
}

func (r *Repo) CreateUser(ctx context.Context, user model.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.l.Debug("CreateUser", "user success created")

	r.Users[user.ID] = user
}

func (r *Repo) GetUserById(ctx context.Context, id uuid.UUID) (*model.User, *errors.CustomError) {
	const method = "GetUserById"
	r.mu.RLock()
	defer r.mu.RUnlock()

	if u, found := r.Users[id]; found {
		r.l.Debug(method, "found user", id)
		return &u, nil
	}

	r.l.Error(method, "user not found", errs.ErrUserNotFound)

	return nil, errs.ErrUserNotFound
}
