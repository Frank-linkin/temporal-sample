package main

import (
	"github.com/temporalio/samples-go/interceptor"
	"go.temporal.io/sdk/client"
	sdkinterceptor "go.temporal.io/sdk/interceptors"
	"go.temporal.io/sdk/worker"
	"log"
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

	tracingInceptor := interceptor.NewTracingInterceptor()
	signalInterceptor := interceptor.NewSignalInterceptor()
	workerInterceptors := []sdkinterceptor.WorkflowInterceptor{tracingInceptor, signalInterceptor}
	w := worker.New(c, "interceptor", worker.Options{
		// Create interceptor that will put started time on the logger
		WorkflowInterceptorChainFactories: workerInterceptors,
	})

	w.RegisterWorkflow(interceptor.Workflow)
	w.RegisterActivity(interceptor.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
