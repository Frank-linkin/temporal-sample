package main

import (
	"context"
	"fmt"
	"log"

	"github.com/temporalio/samples-go/interceptor"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/metrics"
)

func main() {
	metricsScope, metricsScopeCloser, metricsReporter := metrics.NewTaggedMetricsScope()

	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort:     "106.13.193.55:7233",
		MetricsScope: metricsScope,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "interceptor_workflowID",
		TaskQueue: "interceptor",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, interceptor.Workflow, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result string
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)
	//start report
	fmt.Printf("count=%v", metricsReporter.Counts())
	metricsScopeCloser.Close()
}
