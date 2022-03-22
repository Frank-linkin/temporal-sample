package main

import (
	signalDemo "github.com/temporalio/samples-go/signal"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "DIY-Signal", worker.Options{
		EnableLoggingInReplay:        true,
		StickyScheduleToStartTimeout: 1000 * time.Second,
		WorkerStopTimeout:            500 * time.Second,
	})

	w.RegisterWorkflow(signalDemo.SignalWorkflow)
	var hello signalDemo.HelloActivity
	w.RegisterActivity(&hello)
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
