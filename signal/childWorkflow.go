package signalDemo

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"
	"log"
	"time"
)

// Workflow is to demo how to setup query handler
func SignalWorkflow_child(ctx workflow.Context) error {
	fmt.Println("SignalWorkflow_child started")
	//向父workflow发送Signal信号
	temporalClient, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
		return err
	}
	wfExe := workflow.GetInfo(ctx).ParentWorkflowExecution
	//signalVal := SignalData{Age: 17, Name: "Joehanm"}
	signalVal := "my Signal Message"
	err = temporalClient.SignalWorkflow(context.Background(), wfExe.ID, wfExe.RunID, "MySignal", signalVal)
	if err != nil {
		log.Fatalln("Error signaling client", err)
		return err
	}
	//err := workflow.SignalExternalWorkflow(ctx, wfExe.ID, wfExe.RunID, "MySignal", signalVal).Get(ctx, nil)
	//if err != nil {
	//	log.Fatalln("Error signaling client", err)
	//	return err
	//}
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Second, // such a short timeout to make sample fail over very fast
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var values string
	var helloActivity *HelloActivity_child
	if err := workflow.ExecuteActivity(ctx, helloActivity.MkDir1, "strings").Get(ctx, &values); err != nil {
		workflow.GetLogger(ctx).Error(" Activity failed.", "Error", err)
		return err
	}
	fmt.Println("SignalWorkflow_child complete")

	return nil
}

type HelloActivity_child struct {
}

//func (a *RollbackActivity) MkDir1(ctx workflow.Context, fileDir string) error {
func (a *HelloActivity_child) MkDir1(ctx context.Context, fileDir string) error {
	fmt.Println("Mkdir1_child form child Workflow   ")
	return nil
}

func (a *HelloActivity_child) MkDir1Rollback(fileDir string, ctx interface{}) error {
	// cmd := exec.Command("rm", "-rf", fileDir)
	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Println("Rollback fail err: ", err)
	// } else {
	// 	fmt.Println("Rollback suc")
	// }
	fmt.Println("Mkdir1 Rollback is done")
	// return &sangfor.SfError{
	// 	Code:    111,
	// 	Message: "ErrMessafeOfMine\n",
	// 	ErrType: "myType\n",
	// }
	return nil
}
