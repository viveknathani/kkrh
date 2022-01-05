package database

import (
	"log"
	"testing"
)

func TestInitializeAndClose(t *testing.T) {

	db := &Database{}
	dsn := "postgres://viveknathani:root@localhost:5432/kkrhdb"
	err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.pool.Query("select version();")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Close()
	if err != nil {
		log.Fatal(err)
	}
}
