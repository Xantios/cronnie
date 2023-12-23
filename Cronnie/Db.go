package Cronnie

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupDb(host, username, password, dbname string, port int, ssl bool, timezone string) (*pgxpool.Pool, error) {

	sslMode := "disable"
	if ssl {
		sslMode = "enable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, username, password, dbname, port, sslMode, timezone)

	ctx := context.Background()
	conn, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return conn, Seed(conn)
}

func Seed(conn *pgxpool.Pool) error {
	var e error
	e = createJobsTable(conn)
	if e != nil {
		return e
	}

	e = createNotificationProcedure(conn)
	if e != nil {
		return e
	}

	e = createNotificationTrigger(conn)
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
