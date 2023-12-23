package Cronnie

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (ci *Instance) Seed() error {
	var e error
	e = createJobsTable(ci.conn)
	if e != nil {
		return e
	}

	e = createNotificationProcedure(ci.conn)
	if e != nil {
		return e
	}

	e = createNotificationTrigger(ci.conn)
	if e != nil {
		return e
	}

	return nil
}

func createJobsTable(conn *pgxpool.Pool) error {
	q := `create table if not exists jobs
	(
		id           serial,
		function     varchar(128) not null,
		arguments    json,
		created_at   timestamp,
		completed_at date
	);`

	ctx := context.Background()
	_, e := conn.Query(ctx, q)
	return e
}

func createNotificationProcedure(conn *pgxpool.Pool) error {

	ctx := context.Background()

	//language=postgresql
	q := `CREATE OR REPLACE FUNCTION job_notification()
			RETURNS trigger AS
	   	$$
	
				DECLARE
					payload TEXT;
	
				BEGIN
				    payload := json_build_object('data',row_to_json(new));
				    PERFORM pg_notify('job_channel',payload);
				    RETURN NEW;
				END;
	   	$$
	   	LANGUAGE plpgsql VOLATILE`

	_, e := conn.Query(ctx, q)

	return e
}

func createNotificationTrigger(conn *pgxpool.Pool) error {

	ctx := context.Background()

	//language=postgresql
	_, e := conn.Query(ctx, `CREATE OR REPLACE TRIGGER new_job_trigger
    	AFTER INSERT ON "jobs"
    	FOR EACH ROW
    	EXECUTE PROCEDURE job_notification();
    `)

	return e
}
