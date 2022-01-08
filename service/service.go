package service

import (
	"github.com/gomodule/redigo/redis"
	"github.com/viveknathani/kkrh/repository"
	"go.uber.org/zap"
)

type Service struct {
	repo      repository.Repository
	conn      redis.Conn
	jwtSecret []byte
	logger    *zap.Logger
}
