package main

import (
	"Cronnie/Cronnie"
	"fmt"
	"log"
	"os"
	"time"
)

func ExampleJobHandler(args map[string]string) bool {

	fmt.Printf("this is ran totally async as a job managed by Cronnie")

	// Function always succeeds
	return true
}

func main() {

	fmt.Println("This is an example on how to use Cronnie")

	// Craft a logger which we (optionally) can pass around
	logger := log.Default()

	config := &Cronnie.Config{
		// TimeZone=Europe/Amsterdam
		Uri: "host=db user=postgres password=EcArEmp7YTDqCk8Yl61RQEZnLF9oXnjE dbname=app port=5432 sslmode=disable",
		JobMap: map[string]Cronnie.Job{
			"example": ExampleJobHandler,
		},
		Logger: logger,
	}

	cron, err := Cronnie.New(config)
	if err != nil {
		fmt.Printf("Error during setup :: %s\n", err)
		os.Exit(1)
	}

	// Create a job
	// This is a small convenience function, feel free to do an insert

	// Run NOW
	cron.Create("example", map[string]string{
		"arg_1": "This is the first argument",
		"arg_2": "This is the second argument",
	}, time.Time{})

	// run in minute
	cron.Create("example", map[string]string{
		"example_2": "this functions runs in a minute",
	}, time.Now().Add(time.Minute))

	// Example CLI runner
	fmt.Println("Running job worker ")
	go func() {
		e := cron.Run()
		if e != nil {
			fmt.Printf("Error while starting runner :: %s\n", e)
		}
	}()

	// Infinite loop so the example runner keeps running.
	for {
		//
	}
}
