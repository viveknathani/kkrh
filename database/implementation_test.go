package database

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/viveknathani/kkrh/entity"
)

const dsn = "postgres://viveknathani:root@localhost:5432/kkrhdb?sslmode=disable"

var db *Database

func TestMain(t *testing.M) {

	db = &Database{}
	err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}

	// create tables
	_, err = db.pool.Exec("create table if not exists users(id uuid primary key,name varchar not null,email varchar(319) not null,password bytea not null);")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.pool.Exec("create table if not exists logs(id uuid primary key,userId uuid references users(id),latitude double precision,longitude double precision,activity varchar not null,startTime bigint not null,endTime bigint not null,notes varchar);")
	code := t.Run()
	if err != nil {
		log.Fatal(err)
	}

	db.Close()
	os.Exit(code)
}

func TestCreateAndGetUser(t *testing.T) {

	u := &entity.User{
		Name:     "john",
		Email:    "john@gmail.com",
		Password: []byte("someHashedPwd44555"),
	}

	err := db.CreateUser(u)

	if err != nil {
		log.Fatal(err)
	}

	user, err := db.GetUser(u.Email)
	if err != nil {
		log.Fatal(err)
	}

	if user == nil {
		log.Fatal("Got nothing")
	}

	// since entity.User contains a byte array,
	// we cannot use the equality operator to test
	if !reflect.DeepEqual(*u, *user) {
		log.Println("Inequality.")
		log.Println("Created: ", u)
		log.Println("Got: ", user)
		log.Fatal()
	}

	// clean up
	err = db.execWithTransaction("delete from users where name = 'john'")
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateAndGetLogs(t *testing.T) {

	u := &entity.User{
		Name:     "forlogs",
		Email:    "forlogs@gmail.com",
		Password: []byte("someHashedPwd44555"),
	}

	err := db.CreateUser(u)

	if err != nil {
		log.Fatal(err)
	}

	l := &[]entity.Log{
		{
			UserId:    u.Id,
			Latitude:  78.1,
			Longitude: 78.2,
			Activity:  "sleep",
			StartTime: 45678,
			EndTime:   45690,
			Notes:     "",
		},
		{
			UserId:    u.Id,
			Latitude:  78.1,
			Longitude: 78.2,
			Activity:  "work",
			StartTime: 45678,
			EndTime:   0,
		},
	}

	for _, item := range *l {
		err := db.CreateLog(&item)
		if err != nil {
			log.Fatal(err)
		}

		got, err := db.GetLogs(u.Id, item.EndTime)
		if err != nil {
			log.Fatal(err)
		}

		if (*got)[0] != item {
			log.Println(item)
			log.Println((*got)[0])
			log.Fatal("Inequality")
		}

		item.EndTime = 9999
		err = db.UpdateLog(item.Id, 9999)
		if err != nil {
			log.Fatal(err)
		}

		got, err = db.GetLogs(u.Id, item.EndTime)
		if err != nil {
			log.Fatal(err)
		}

		if (*got)[0] != item {
			log.Println(item)
			log.Println((*got)[0])
			log.Fatal("Inequality")
		}

		err = db.execWithTransaction(fmt.Sprintf("delete from logs where id='%s'", item.Id))
		if err != nil {
			log.Fatal(err)
		}
	}

	err = db.execWithTransaction("delete from users where name = 'forlogs'")
	if err != nil {
		log.Fatal(err)
	}
}
