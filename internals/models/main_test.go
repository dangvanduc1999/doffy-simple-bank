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
var dbPoolTest *pgxpool.Pool
var dbPoolError error

func TestMain(m *testing.M) {
	appContext := utils.GetAppTestContext()
	dbPoolTest, dbPoolError = pgxpool.New(context.Background(), appContext.MainDBConnection)

	if dbPoolError != nil {
		log.Fatal("Cannot access to database", dbPoolError)
		return
	}
	defer dbPoolTest.Close()

	AccountQueries = New(dbPoolTest)

	code := m.Run()
	appContext.Teardown()
	os.Exit(code)
}
