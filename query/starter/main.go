package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/query"
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
		ID:        "query_workflow",
		TaskQueue: "DIYQuery",
		WorkflowExecutionTimeout: 1000*time.Second,
		WorkflowTaskTimeout: 1000*time.Second,
		WorkflowRunTimeout: 1000*time.Second,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, query.QueryWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
