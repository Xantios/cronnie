package Cronnie

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Instance struct {
	ctx    context.Context
	conn   *pgxpool.Pool
	jobMap map[string]Job
}

func New(config *Config) (Instance, error) {

	var instance Instance

	instance.ctx = context.Background()

	// check if we have a direct connection set or if we have a connection string to connect with
	if config.Connection == nil && config.Uri != "" {
		conn, err := pgxpool.New(instance.ctx, config.Uri)
		if err != nil {
			return Instance{}, err
		}

		instance.conn = conn
	} else {
		instance.conn = config.Connection
	}

	// Run Seeder
	e := instance.Seed()
	if e != nil {
		return Instance{}, e
	}

	// Set jobMap
	instance.jobMap = config.JobMap

	return instance, nil
}
