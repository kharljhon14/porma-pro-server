package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var testStore Store

func TestMain(m *testing.M) {
	var err error

	testDBPool, err := pgxpool.New(context.Background(), os.Getenv("DSN"))
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testStore = NewStore(testDBPool)

	os.Exit(m.Run())
}
