package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/viveknathani/kkrh/cache"
	"github.com/viveknathani/kkrh/database"
	"github.com/viveknathani/kkrh/server"
	"github.com/viveknathani/kkrh/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	port           string = ""
	databaseServer string = ""
	redisServer    string = ""
	jwtSecret      string = ""
	allowedOrigin  string = "*"
)

func init() {

	mode := os.Getenv("MODE")

	if mode == "dev" {
		port = os.Getenv("DEV_PORT")
		databaseServer = os.Getenv("DEV_DATABASE_URL")
		redisServer = os.Getenv("DEV_REDIS_URL")
		jwtSecret = os.Getenv("DEV_JWT_SECRET")
	} else if mode == "prod" {
		port = os.Getenv("PORT")
		databaseServer = os.Getenv("DATABASE_URL")
		redisServer = os.Getenv("REDIS_URL")
		jwtSecret = os.Getenv("JWT_SECRET")
		allowedOrigin = os.Getenv("ALLOWED_ORIGIN")
	} else {
		port = "8080"
		databaseServer = "postgres://viveknathani:root@localhost:5432/kkrhdb"
		redisServer = "127.0.0.1:6379"
		jwtSecret = "hey"
	}
}

func main() {

	// Setup logger
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
	logger, err := cfg.Build()
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// Setup database
	db := &database.Database{}
	err = db.Initialize(databaseServer)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// Setup cache
	memory := &cache.Cache{}
	memory.Initialize(redisServer)
	memoryConn := memory.Pool.Get()

	// Setup the web server
	srv := &server.Server{
		Service: &service.Service{
			Repo:      db,
			Conn:      memoryConn,
			JwtSecret: []byte(jwtSecret),
			Logger:    logger,
		},
		Router: mux.NewRouter(),
	}

	srv.SetupRoutes()

	crossOrigin := cors.New(cors.Options{
		AllowedOrigins:   []string{allowedOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut},
	})

	// Listen
	err = http.ListenAndServe(":"+port, crossOrigin.Handler(handlers.ProxyHeaders(srv)))
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	err = srv.Service.Conn.Close()
	if err != nil {
		fmt.Print(err)
	}
	err = memory.Close()
	if err != nil {
		fmt.Print(err)
	}
	err = db.Close()
	if err != nil {
		fmt.Print(err)
	}
	err = srv.Service.Logger.Sync()
	if err != nil {
		fmt.Print(err)
	}
}
