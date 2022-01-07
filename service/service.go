package service

import (
	"github.com/gomodule/redigo/redis"
	"github.com/viveknathani/kkrh/repository"
)

type Service struct {
	repo      repository.Repository
	conn      redis.Conn
	jwtSecret []byte
}
