package service

import (
	"github.com/gomodule/redigo/redis"
	"github.com/viveknathani/kkrh/repository"
	"go.uber.org/zap"
)

type Service struct {
	Repo      repository.Repository
	Conn      redis.Conn
	JwtSecret []byte
	Logger    *zap.Logger
}
