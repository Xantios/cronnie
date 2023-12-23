package main

import (
	"Cronnie/Cronnie"
	"fmt"
	"log"
)

func ExampleJobHandler(args map[string]string) bool {

	fmt.Printf("this is ran totally async as a job managed by Cronnie")

	// Function always succeeds
	return true
}

func main() {

	fmt.Println("This is an example on how to use Cronnie")

	// Set up a connection and use it
	conn, e := Cronnie.SetupDb("db", "postgres", "EcArEmp7YTDqCk8Yl61RQEZnLF9oXnjE", "app", 5432, false, "Europe/Amsterdam")
	if e != nil {
		log.Fatalf("Cant setup DB for use with Cronnie. Error :: %s\n", e)
	}

	// Or bring your own connection
	// conn,err := pgxpool.New(context.Background(),"string")
	// e := Cronnie.Seed(conn)

	fmt.Println("Ready to store jobs!")

	// Bind names to functions, so we can run them when the time is right
	Cronnie.Runners(map[string]Cronnie.Job{
		"example": ExampleJobHandler,
	})

	// Create a job
	// This is a small convenience function, feel free to do an insert
	arguments := map[string]string{}
	arguments["arg_1"] = "This is the first argument"
	arguments["arg_2"] = "This is the second argument"

	e = Cronnie.Create(conn, "example", arguments)
	if e != nil {
		fmt.Printf("Cant create job. something wrong with your database? Error :: %s\n", e)
	}

	// Example CLI runner
	fmt.Println("Running job worker ")
	go func() {
		e := Cronnie.Run(conn)
		if e != nil {
			fmt.Printf("Error while starting runner :: %s\n", e)
		}
	}()

	// Infinite loop so the example runner keeps running.
	for {
		//
	}
}
