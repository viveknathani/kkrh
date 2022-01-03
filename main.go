package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"database/sql"

	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq"
)

func main() {

	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		io.WriteString(rw, "hello from kkrh-backend")
	})

	PORT := ""
	REDIS_URL := ""
	DATABASE_URL := ""

	if os.Getenv("MODE") == "dev" {
		PORT = os.Getenv("DEV_PORT")
		REDIS_URL = os.Getenv("DEV_REDIS_URL")
		DATABASE_URL = os.Getenv("DEV_DATABASE_URL")
	} else {
		PORT = os.Getenv("PORT")
		REDIS_URL = os.Getenv("REDIS_URL")
		DATABASE_URL = os.Getenv("DEV_DATABASE_URL")
	}

	_, err := sql.Open("postgres", DATABASE_URL)
	if err != nil {
		log.Fatal(err)
	}

	_, err = redis.Dial("tcp", REDIS_URL)
	if err != nil {
		log.Fatal(err)
	}

	err = http.ListenAndServe(":"+PORT, nil)
	if err != nil {
		log.Fatal(err)
	}
}
