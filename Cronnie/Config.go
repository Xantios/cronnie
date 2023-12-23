package Cronnie

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Config struct {
	Connection *pgxpool.Pool
	Uri        string
	JobMap     map[string]Job
	Logger     *log.Logger
}

func (c *Config) SetDB(pool *pgxpool.Pool) *Config {
	c.Connection = pool
	return c
}

// addJob
// removeJob
