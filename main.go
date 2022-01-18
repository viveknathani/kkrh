package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/viveknathani/kkrh/cache"
	"github.com/viveknathani/kkrh/database"
	"github.com/viveknathani/kkrh/server"
	"github.com/viveknathani/kkrh/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Hold environment variables
var (
	port           string = ""
	databaseServer string = ""
	redisServer    string = ""
	redisUsername  string = ""
	redisPassword  string = ""
	jwtSecret      string = ""
)

// Setup environment variables
func init() {

	mode := os.Getenv("MODE")

	if mode == "dev" {
		port = os.Getenv("DEV_PORT")
		databaseServer = os.Getenv("DEV_DATABASE_URL")
		redisServer = os.Getenv("DEV_REDIS_URL")
		jwtSecret = os.Getenv("DEV_JWT_SECRET")
		fmt.Println("here")
	}

	if mode == "prod" {
		port = os.Getenv("PORT")
		databaseServer = os.Getenv("DATABASE_URL")
		redisServer = os.Getenv("REDIS_URL")
		redisUsername = os.Getenv("REDIS_USERNAME")
		redisPassword = os.Getenv("REDIS_PASSWORD")
		jwtSecret = os.Getenv("JWT_SECRET")
	}
}

// getLogger will configure and return a uber/zap logger
func getLogger() *zap.Logger {

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

	return logger
}

// getDatabase will init and return a db
func getDatabase() *database.Database {

	db := &database.Database{}
	err := db.Initialize(databaseServer)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	return db
}

// getCache will return a connection to Redis from the pool
func getCache() (*cache.Cache, redis.Conn) {

	memory := &cache.Cache{}
	memory.Initialize(redisServer, redisUsername, redisPassword)
	memoryConn := memory.Pool.Get()
	return memory, memoryConn
}

func main() {

	logger := getLogger()
	db := getDatabase()
	memory, memoryConn := getCache()

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

	var secureHandler *server.SecurityHandler = nil

	// middleware for better HTTP headers
	if os.Getenv("MODE") == "prod" {
		secureHandler = server.NewSecurityHandler(srv)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Listen
	go func() {

		if secureHandler != nil {

			err := http.ListenAndServe(":"+port, secureHandler)
			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
		} else {
			err := http.ListenAndServe(":"+port, srv)
			if err != nil {
				fmt.Print(err)
				os.Exit(1)
			}
		}
	}()
	fmt.Println("Server started!")
	<-done
	shutdown(srv, db, memory)
}

func shutdown(srv *server.Server, db *database.Database, memory *cache.Cache) {

	err := srv.Service.Conn.Close()
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
	fmt.Println("goodbye!")
}
