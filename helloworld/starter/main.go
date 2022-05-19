package main

import (
	"context"
	"log"
	"time"

	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/helloworld"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: "106.13.193.55:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	memo := make(map[string]interface{})
	memo["火影"] = helloworld.Param{Age: 18, Name: "Memo漩涡"}
	memo["死神"] = helloworld.Param{Age: 18, Name: "Memo一户"}
	workflowOptions := client.StartWorkflowOptions{
		ID:                       "hello_world_workflowID",
		TaskQueue:                "hello-world",
		WorkflowRunTimeout:       300 * time.Second,
		WorkflowExecutionTimeout: 300 * time.Second,
		Memo:                     memo,
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworld.Workflow, helloworld.Param{Age: 19, Name: "漩涡鸣人"}, helloworld.Param{Age: 19, Name: "刘亦菲"})
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
}
