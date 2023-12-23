package main

import (
	"Cronnie/Cronnie"
	"fmt"
	"os"
)

func ExampleJobHandler(args map[string]string) bool {

	fmt.Printf("this is ran totally async as a job managed by Cronnie")

	// Function always succeeds
	return true
}

func main() {

	fmt.Println("This is an example on how to use Cronnie")

	config := &Cronnie.Config{
		Uri: "host=db user=postgres password=EcArEmp7YTDqCk8Yl61RQEZnLF9oXnjE dbname=app port=5432 sslmode=disable TimeZone=Europe/Amsterdam",
		JobMap: map[string]Cronnie.Job{
			"example": ExampleJobHandler,
		},
	}

	cron, err := Cronnie.New(config)
	if err != nil {
		fmt.Printf("Error during setup :: %s\n", err)
		os.Exit(1)
	}

	// Create a job
	// This is a small convenience function, feel free to do an insert
	e := cron.Create("example", map[string]string{
		"arg_1": "This is the first argument",
		"arg_2": "This is the second argument",
	})

	if e != nil {
		fmt.Printf("Cant create job. something wrong with your database? Error :: %s\n", e)
	}

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
