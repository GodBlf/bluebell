package dao

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type User struct {
	rclient *redis.Client
	mclient *sqlx.DB
}

func NewUser(r *redis.Client, m *sqlx.DB) *User {
	return &User{
		rclient: r,
		mclient: m,
	}
}
