package service

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/viveknathani/kkrh/cache"
	"github.com/viveknathani/kkrh/database"
	"github.com/viveknathani/kkrh/entity"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const dsn = "postgres://viveknathani:root@localhost:5432/kkrhdb?sslmode=disable"
const redisurl = "127.0.0.1:6379"

var service *Service

func TestMain(t *testing.M) {

	service = &Service{}
	db := &database.Database{}
	c := &cache.Cache{}
	err := db.Initialize(dsn)
	if err != nil {
		log.Fatal(err)
	}
	c.Initialize(redisurl, "", "")
	service.Repo = db
	service.Conn = c.Pool.Get()
	service.JwtSecret = []byte("secret")
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevel(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,
			TimeKey:     "ts",
			EncodeTime:  zapcore.EpochMillisTimeEncoder,
		},
	}
	logger, _ := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	service.Logger = logger
	code := t.Run()
	if err != nil {
		log.Fatal(err)
	}
	db.Close()
	_ = logger.Sync()
	os.Exit(code)
}

func TestUserAuth(t *testing.T) {

	var u entity.User = entity.User{
		Name:     "John Doe",
		Email:    "johnhello@gmail.com",
		Password: []byte("Hello2402@#"),
	}

	p := u.Password

	err := service.Signup(context.Background(), &u)
	if err != nil {
		log.Fatal(err)
	}

	u.Password = p
	token, err := service.Login(context.Background(), &u)
	if err != nil {
		log.Fatal(err)
	}

	payload, err := service.VerifyAndDecodeToken(context.Background(), token)
	if err != nil {
		log.Fatal(err)
	}

	if payload != u.Id {
		log.Fatal("id mismatch")
	}

	err = service.Logout(context.Background(), token)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Repo.DeleteUser(u.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateLog(t *testing.T) {

	var u entity.User = entity.User{
		Name:     "John Doe",
		Email:    "nobrainer@gmail.com",
		Password: []byte("Hello2402@#"),
	}

	err := service.Signup(context.Background(), &u)
	if err != nil {
		log.Fatal(err)
	}

	l := &entity.Log{
		UserId:    u.Id,
		Latitude:  78.1,
		Longitude: 78.2,
		Activity:  "sleep",
		StartTime: 45678,
		EndTime:   4569000000,
		Notes:     "",
	}

	err = service.CreateLog(context.Background(), l)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Repo.DeleteLog(l.Id)
	if err != nil {
		log.Fatal(err)
	}

	err = service.Repo.DeleteUser(u.Id)
	if err != nil {
		log.Fatal(err)
	}
}
