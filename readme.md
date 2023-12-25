<p align="center">
    <img src="assets/logo.png" alt="Cronnie Logo" width="320px">
</p>

# Cronnie

A easy to use but yet powerfull job schedueler written in Go

### How does it work?
The basic concept here is the Cronnie connects to a postgres table which it manages itself. 
by the magic of postgres we can subscribe to the changes on this table, so if a new job is created Cronnie figures out what to do and when to do it

most of the database specific handling is in the db file if you want to deep dive into it. 
the idea is that Cronnie is fully event based so there is no polling in this code 

## How to install ?
As with any golang package you can just use the `go get` command. 

```bash
go get github.com/xantios/cronnie
```

## How to use ?
For a full code example see [example/main.go](example/main.go)

### JobMaps 
Cronie heavily relies on the concept of a JobMap. this is simple a `map[string]cronnie.Job` a cronnie.Job is a function which excepts a `map[string]string` 
A quick example

```golang
package main 

func ExampleJobHandler(arg map[string]string) bool {
	// Example 
	return true
}

func main() {
    exampleMap := map[string]Cronnie.Job{
    "example": ExampleJobHandler,
    },
}
```

#### Create a config
A simple config can be created using the following snippet

```golang
package main

import Cronnie "github.com/xantios/cronnie"

func main() {
	config := &Cronnie.Config{
		Uri: "host=db user=MY_USER password=PASSWORD dbname=app port=5432 sslmode=disable",
		JobMap: map[string]Cronnie.Job{
			"example": ExampleJobHandler,
		},
		Logger: logger,
	}
}
```

### Create a instance
Create a instance to work with, if you have a config this is easy.

```golang
cron, err := Cronnie.New(config)
```

### Add a task
To add a task, take the reference created above (or any instance for that matter) and take a look at the following example code.
a more involved example can be found in the example directory

```golang
	cron.Create("example", map[string]string{
		"arg_1": "This is the first argument",
		"arg_2": "This is the second argument",
	}, time.Time{})
```

### Running job worker
This is where the actual magic happens!
Just simply run the `Run()` function should suffice. 
```golang
e := cron.Run()
```
