package coredb

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/dangvanduc1999/doffy-simple-bank/internals/models"
	"github.com/dangvanduc1999/doffy-simple-bank/internals/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type CoreDBConfig struct {
	UserName string
	Password string
	Host     string
	Port     string
	DBName   string
	Options  string
}

func NewCoreDB(config CoreDBConfig, autoMigrate bool) (*models.Queries, *pgxpool.Pool, error) {
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

	model := models.New(pool)

	if autoMigrate {
		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatalf("failed to set goose dialect: %s", err)
		}

		migrationsDir, err := utils.FindMigrationsDir()
		if err != nil {
			log.Fatalf("failed to find migrations directory: %s", err)
		}

		dbMigration, err := sql.Open("pgx", connStr)
		if err != nil {
			log.Fatalf("failed to open db connection for migration: %s", err)
		}
		defer dbMigration.Close()

		if err := goose.Up(dbMigration, migrationsDir); err != nil {
			log.Fatalf("failed to run migrations: %s", err)
		}
	}

	return model, pool, nil
}
