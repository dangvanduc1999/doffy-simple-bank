package coredb

import (
	"context"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestNewCoreDB(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18.1-alpine3.23",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "simple_bank",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer container.Terminate(ctx)

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatal(err)
	}

	config := CoreDBConfig{
		UserName: "test",
		Password: "test",
		Host:     host,
		Port:     port.Port(),
		DBName:   "simple_bank",
		Options:  "sslmode=disable",
	}

	queries, pool, err := NewCoreDB(config)
	if err != nil {
		t.Fatalf("NewCoreDB failed: %v", err)
	}
	defer pool.Close()

	if queries == nil {
		t.Error("expected non-nil queries")
	}
	if pool == nil {
		t.Error("expected non-nil pool")
	}

	if err := pool.Ping(ctx); err != nil {
		t.Errorf("failed to ping database: %v", err)
	}
}
