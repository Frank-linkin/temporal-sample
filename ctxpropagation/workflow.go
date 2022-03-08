package ctxpropagation

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

// CtxPropWorkflow workflow definition
func CtxPropWorkflow(ctx workflow.Context) (err error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Second, // such a short timeout to make sample fail over very fast
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	info:=workflow.GetInfo(ctx)
	if(info.ParentWorkflowExecution == nil) {
		fmt.Println(info.WorkflowExecution.ID)
	}
	if val := ctx.Value(LogIdPropagateKey); val != nil {
		fmt.Println("Father:"+val.(string))
	}



	var values Values
	if err = workflow.ExecuteActivity(ctx, SampleActivity).Get(ctx, &values); err != nil {
		workflow.GetLogger(ctx).Error("Workflow failed.", "Error", err)
		return err
	}
	//workflow.GetLogger(ctx).Info("context propagated to activity", values.Key, values.Value)
	childWorkflowOptions := workflow.ChildWorkflowOptions{}
	ctx = workflow.WithChildOptions(ctx, childWorkflowOptions)

	var result interface{}
	err = workflow.ExecuteChildWorkflow(ctx, ChildWorkflowSample).Get(ctx, &result)
	if err != nil {
		// ...
	}
	workflow.GetLogger(ctx).Info("Workflow completed.")
	return nil
}
