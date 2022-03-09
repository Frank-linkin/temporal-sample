package query

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/workflow"
	"time"
)

// Workflow is to demo how to setup query handler
func QueryWorkflow(ctx workflow.Context) error {
	queryResult := "started"
	//attr2 := map[string]interface{}{
	//	//"CustomIntField": -1,
	//	"CustomKeywordField": "seattle",
	//}
	//workflow.UpsertSearchAttributes(ctx, attr2)
	logger := workflow.GetLogger(ctx)
	logger.Info("QueryWorkflow started")
	// setup query handler for query type "state"
	err := workflow.SetQueryHandler(ctx, "state", func(str string) (string ,error) {
		var res string

		attr1 := map[string]interface{}{
			"CustomIntField": 15,
			"CustomKeywordField": "seattle",
		}
		workflow.UpsertSearchAttributes(ctx, attr1)
		fmt.Printf("++++++joehanms is : %v",str)
		res = queryResult
		return res, nil
	})
	if err != nil {
		logger.Info("SetQueryHandler failed: " + err.Error())
		return err
	}

	attr1 := map[string]interface{}{
		"CustomIntField": -1,
	}
	workflow.UpsertSearchAttributes(ctx, attr1)
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 2 * time.Second, // such a short timeout to make sample fail over very fast
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var values string
	var helloActivity *HelloActivity
	if err := workflow.ExecuteActivity(ctx, helloActivity.MkDir1,"strings").Get(ctx, &values); err != nil {
		workflow.GetLogger(ctx).Error(" Activity failed.", "Error", err)
		return err
	}
	queryResult = "waiting on timer"
	// to simulate workflow been blocked on something, in reality, workflow could wait on anything like activity, signal or timer
	_ = workflow.NewTimer(ctx, time.Second*300).Get(ctx, nil)
	logger.Info("Timer fired")

	queryResult = "done"
	fmt.Println(queryResult)
	logger.Info("QueryWorkflow completed")

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
	for i:=0;i<=1000000000;i++{
		//fmt.Println("this is Acitvity")
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
