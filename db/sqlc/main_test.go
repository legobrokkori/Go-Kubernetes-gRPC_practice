package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/legobrokkori/go-kubernetes-grpc_practice/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	testDB, err := sql.Open(config.DbDriver, config.DbServer)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
