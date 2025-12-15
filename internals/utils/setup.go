package utils

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type AppTestContext struct {
	MainDBConnection string
	Teardown         func()
}

var lock = &sync.Mutex{}

var AppTestContextInstance *AppTestContext

func GetAppTestContext() *AppTestContext {
	if AppTestContextInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		AppTestContextInstance = Setup()
	}
	return AppTestContextInstance
}
func Setup() *AppTestContext {
	ctx := context.Background()
	dbClient, err := testcontainers.Run(
		ctx,
		"postgres:18.1-alpine3.23",
		testcontainers.WithExposedPorts("5432/tcp"),
		testcontainers.WithEnv(map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "simple_bank",
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	host, err := dbClient.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}

	port, err := dbClient.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatal(err)
	}

	endpoint := host + ":" + port.Port()

	teardown := func() {
		if err := dbClient.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err)
		}
	}

	connStr := "postgres://test:test@" + endpoint + "/simple_bank?sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("failed to open db connection for migration: %s", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %s", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %s", err)
	}

	migrationsDir, err := findMigrationsDir()
	if err != nil {
		log.Fatalf("failed to find migrations directory: %s", err)
	}

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}

	return &AppTestContext{
		MainDBConnection: connStr,
		Teardown:         teardown,
	}
}

func findMigrationsDir() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		foundPath := filepath.Join(path, "migrations")
		if _, err := os.Stat(foundPath); err == nil {
			return foundPath, nil
		}
		parent := filepath.Dir(path)
		if parent == path {
			return "", os.ErrNotExist
		}
		path = parent
	}
}
