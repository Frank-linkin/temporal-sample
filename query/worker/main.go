package main

import (
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/query"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: "106.13.193.55:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "DIYQuery", worker.Options{
		EnableLoggingInReplay:        true,
		StickyScheduleToStartTimeout: 1000 * time.Second,
		WorkerStopTimeout:            500 * time.Second,
	})

	w.RegisterWorkflow(query.QueryWorkflow)
	var hello query.HelloActivity
	w.RegisterActivity(&hello)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
