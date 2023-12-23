package Cronnie

import (
	"context"
	"github.com/jackc/pgx/v5"
)

func (ci *Instance) Create(functionName string, arguments map[string]string) error {
	ctx := context.Background()
	//language=postgresql
	query := `
		INSERT INTO 
			public.jobs (id, function, arguments, created_at, completed_at)
		VALUES 
		    (DEFAULT, $1, $2, now(), null);`

	_, e := ci.conn.Query(ctx, query, functionName, arguments)
	return e
}

func (ci *Instance) GetJobs() (pgx.Rows, error) {
	ctx := context.Background()
	rows, err := ci.conn.Query(ctx, "SELECT * FROM jobs WHERE completed_at is null")
	return rows, err
}
