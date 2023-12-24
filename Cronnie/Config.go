package Cronnie

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

type Config struct {
	Connection    *pgxpool.Pool
	Uri           string
	JobMap        map[string]Job
	Logger        *log.Logger
	KeepCompleted time.Duration
}

func (c *Config) SetDB(pool *pgxpool.Pool) *Config {
	c.Connection = pool
	return c
}

// addJob
// removeJob
