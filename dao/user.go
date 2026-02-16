package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type User struct {
	rclient *redis.Client
	mclient *sqlx.DB
}
