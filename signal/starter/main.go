package main

import (
	"context"
	signalDemo "github.com/temporalio/samples-go/signal"
	"log"
	"time"

	"go.temporal.io/sdk/client"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:                       "Signal-workflow",
		TaskQueue:                "DIY-Signal",
		WorkflowExecutionTimeout: 1000 * time.Second,
		WorkflowTaskTimeout:      1000 * time.Second,
		WorkflowRunTimeout:       1000 * time.Second,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, signalDemo.SignalWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
