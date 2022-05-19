package signalDemo

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

type SignalData struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

// Workflow is to demo how to setup query handler
func SignalWorkflow(ctx workflow.Context) error {
	signalChan := workflow.GetSignalChannel(ctx, "MySignal")
	selector := workflow.NewSelector(ctx)

	selector.AddReceive(signalChan, func(channel workflow.ReceiveChannel, more bool) {
		var signalData string
		channel.Receive(ctx, &signalData)
		//fmt.Printf("\nsignalData.name = %v", signalData.Name)
		//fmt.Printf("\nsignalData.age = %v", signalData.Age)
		fmt.Printf("my signal Data is %v", signalData)
	})

	fmt.Println("注册Signal Handler成功 ")
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Second, // such a short timeout to make sample fail over very fast
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var values string
	var helloActivity *HelloActivity
	if err := workflow.ExecuteActivity(ctx, helloActivity.MkDir1, "strings").Get(ctx, &values); err != nil {
		workflow.GetLogger(ctx).Error(" Activity failed.", "Error", err)
		return err
	}

	cwo := workflow.ChildWorkflowOptions{
		WorkflowID: "ABC-SIMPLE-CHILD-WORKFLOW-ID",
		TaskQueue:  "DIY-Signal-child",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	var result string
	err := workflow.ExecuteChildWorkflow(ctx, SignalWorkflow_child, "World").Get(ctx, &result)
	if err != nil {
		return err
	}
	_ = workflow.NewTimer(ctx, time.Second*5).Get(ctx, nil)
	selector.Select(ctx)

	fmt.Printf("\nroot workflow finished")

	return nil
}

type HelloActivity struct {
}

//func (a *RollbackActivity) MkDir1(ctx workflow.Context, fileDir string) error {
func (a *HelloActivity) MkDir1(ctx context.Context, fileDir string) error {
	fmt.Println("Mkdir1 form rootWorkflow   ")
	return nil
}

func (a *HelloActivity) MkDir1Rollback(fileDir string, ctx interface{}) error {
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
