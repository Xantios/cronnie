package Cronnie

import (
	"context"
	"github.com/jackc/pgx/v5"
	"time"
)

func (ci *Instance) Create(functionName string, arguments map[string]string, when time.Time) error {

	ctx := context.Background()

	if when.IsZero() { // The initial state of a time.Time{} is 001-01-01 00:00:00
		when = time.Now()
	}

	//language=postgresql
	query := `
		INSERT INTO 
			public.jobs (id, function, arguments,run_at, created_at, completed_at)
		VALUES 
		    (DEFAULT, $1, $2, $3,now(), null);`

	_, e := ci.conn.Query(ctx, query, functionName, arguments, when)
	return e
}

func (ci *Instance) GetJobs() (pgx.Rows, error) {
	ctx := context.Background()
	rows, err := ci.conn.Query(ctx, "SELECT * FROM jobs WHERE completed_at is null")
	return rows, err
}

func (ci *Instance) MarkCompleted(id int) error {

	// language=postgresql
	q := `
		UPDATE jobs 
		SET completed_at = now()
		WHERE id = $1
	`

	_, e := ci.conn.Query(ci.ctx, q, id)
	return e
}

func setReminder(job JobModel) {
	timeout := time.Until(job.RunAt.Time)
	go func() {
		time.Sleep(timeout)
		jobChannel <- job
	}()
}
