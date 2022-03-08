package ctxpropagation

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
	"time"
)

func ChildWorkflowSample(ctx workflow.Context) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Second, // such a short timeout to make sample fail over very fast
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	//if val := ctx.Value(PropagateKey); val != nil {
	//	vals := val.(Values)
	//	workflow.GetLogger(ctx).Info("custom context propagated to childworkflow", vals.Key, vals.Value)
	//}
	logger := workflow.GetLogger(ctx)
	if val := ctx.Value(LogIdPropagateKey); val != nil {
		fmt.Println(val)
		logger.Info("ChildWorkflow PropagateKey", "Logid", val)

	}
	//if val := ctx.Value(LogIdPropagateKey); val != nil {
	//	vals := val.(Values)
	//	fmt.Printf("ChildWorkflow Key is : %v",vals.Key)
	//	logger.Info("ChildWorkflow PropagateKey", "Key", vals.Key)
	//}
	var values Values
	var helloActivity *HelloActivity
	if err := workflow.ExecuteActivity(ctx, helloActivity.MkDir1,"strings").Get(ctx, &values); err != nil {
		workflow.GetLogger(ctx).Error("child Workflow failed.", "Error", err)
		return err
	}
	//workflow.GetLogger(ctx).Info("context propagated to activity", values.Key, values.Value)
	//workflow.GetLogger(ctx).Info("Workflow completed.")
	return nil
}


type HelloActivity struct {
}

//func (a *RollbackActivity) MkDir1(ctx workflow.Context, fileDir string) error {
func (a *HelloActivity) MkDir1(ctx context.Context,fileDir string) error {
	// return sangfor.NewError("Error message of Mine", "MyErrorTypeInactivity", nil, &sangfor.ErrResp{123})
	//attr1 := map[string]interface{}{
	//	"CustomIntField":  77,
	//	"CustomBoolField": true,
	//}
	//context := ctx.(workflow.Context)
	//
	//workflow.UpsertSearchAttributes(context, attr1)
	logger:=activity.GetLogger(ctx)
	fmt.Println("Acitvity from childworkflow\n")

	if val := ctx.Value(LogIdPropagateKey); val != nil {
		fmt.Println(val)
		logger.Info("ChildWorkflow PropagateKey", "Logid", val)

	}
	return nil
	// f(3)
	// cmd := exec.Command("mkdir", fileDir)
	// err := cmd.Run()
	// if err != nil {
	// 	fmt.Println("RollbackActivity1 fail err: ", err)
	// } else {
	// 	fmt.Println("RollbackActivity1 suc")
	// }
	// return nil
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
