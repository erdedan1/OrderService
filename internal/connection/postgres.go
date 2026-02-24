package connection

import (
	"context"
	"fmt"
	"log"
	"time"

	"OrderService/config"

	error "github.com/erdedan1/shared/errs"
	"github.com/jmoiron/sqlx"
)

const (
	defaultConnMaxLifetime = time.Minute * 3
	defaultMaxIdleConns    = 10
	defaultMaxOpenConns    = 10
	defaultRetryDelay      = time.Second
	defaultRetryTimeout    = time.Second * 10
)

func New(ctx context.Context, config config.PostgresDB) (*sqlx.DB, *error.CustomError) {
	dbURI := getDBURI(config)

	delayTimer := time.NewTimer(time.Duration(0))
	timeoutExceeded := time.After(defaultRetryTimeout)

	for {
		select {
		case <-timeoutExceeded:
			return nil, error.New(2, "timeout connect")
		case <-delayTimer.C:
			client, err := sqlx.ConnectContext(ctx, "postgres", dbURI)
			if err != nil {
				log.Printf("db connection failed: %s", err)
				delayTimer.Reset(defaultRetryDelay)
				continue
			}

			client.SetConnMaxLifetime(defaultConnMaxLifetime)
			client.SetMaxIdleConns(defaultMaxIdleConns)
			client.SetMaxOpenConns(defaultMaxOpenConns)

			return client, nil
		}
	}
}

func NewHandle(config config.PostgresDB) *sqlx.DB {
	client, err := sqlx.Open("postgres", getServerURI(config))
	if err != nil {
		panic(err)
	}
	return client
}

func getServerURI(config config.PostgresDB) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
	)
}

func getDBURI(config config.PostgresDB) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
}
