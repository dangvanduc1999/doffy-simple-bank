package coredb

import (
	"context"
	"fmt"

	"github.com/doffy/simple-bank/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CoreDBConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	DBName   string
	Options  string
}

func NewCoreDB(config CoreDBConfig) (*models.Queries, *pgxpool.Pool, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		config.UserName,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
		config.Options,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to create connection pool: %w", err)
	}

	return models.New(pool), pool, nil
}
