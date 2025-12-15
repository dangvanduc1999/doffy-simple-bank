package models

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/doffy/simple-bank/internals/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

var AccountQueries *Queries

func TestMain(m *testing.M) {
	appContext := utils.GetAppTestContext()
	dbPool, err := pgxpool.New(context.Background(), appContext.MainDBConnection)

	if err != nil {
		log.Fatal("Cannot access to database", err)
		return
	}
	defer dbPool.Close()

	AccountQueries = New(dbPool)

	code := m.Run()
	appContext.Teardown()
	os.Exit(code)
}
