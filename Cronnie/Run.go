package Cronnie

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobModel struct {
	ID          int
	Function    string
	Arguments   map[string]string
	CreatedAt   pgtype.Timestamp
	CompletedAt pgtype.Timestamp
}

var jobChannel chan JobModel

func init() {
	jobChannel = make(chan JobModel, 256)
}

func (ci *Instance) Run() error {

	ctx := context.Background()

	// Convert pool connection to singular connection
	conf := ci.conn.Config()
	conn, e := pgx.Connect(ctx, conf.ConnString())

	// Populate cache if recovering from crash
	jobs, err := ci.GetJobs()
	if err != nil {
		return err
	}

	for jobs.Next() {
		item := JobModel{}
		err = jobs.Scan(&item.ID, &item.Function, &item.Arguments, &item.CreatedAt, &item.CompletedAt)
		if err != nil {
			return err
		}

		ci.logger.Printf("Got item from crash-recovery. %#v\n", item)
		jobChannel <- item
	}

	if e != nil {
		return e
	}

	// Start listening for new events
	_, e = conn.Exec(ctx, `LISTEN job_channel`)
	if e != nil {
		return e
	}

	// Run queueHandler
	go ci.queueHandler()

	for {
		notification, err := conn.WaitForNotification(ctx)
		if err != nil {
			return err
		}

		// Create item
		item := JobModel{}
		e := json.Unmarshal([]byte(notification.Payload), &item)
		if e != nil {
			ci.logger.Printf("Malformed event :: [[%s]] error :: %s\n", notification.Payload, e)
			continue
		}

		// Write to channel
		// ...
	}
}

func (ci *Instance) queueHandler() {
	for {
		select {
		case job := <-jobChannel:
			ci.logger.Printf("Got job :: %#v\n", job)

			ran := ci.executeRunnerFunction(job.Function, job.Arguments)
			if !ran {
				ci.logger.Printf("Error while running job. \n")
			} else {
				ci.logger.Printf("Successfully ran job. \n")
			}
		}
	}
}
