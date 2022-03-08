package query

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
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
		//info := workflow.GetInfo(ctx)
		//payload := info.SearchAttributes.IndexedFields["CustomKeywordField"]
		var res string
		//converter.GetDefaultDataConverter().FromPayload(payload,&res)
		//
		//
		////upsers Search Attribute
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

	queryResult = "waiting on timer"
	// to simulate workflow been blocked on something, in reality, workflow could wait on anything like activity, signal or timer
	_ = workflow.NewTimer(ctx, time.Second*300).Get(ctx, nil)
	logger.Info("Timer fired")

	queryResult = "done"
	fmt.Println(queryResult)
	logger.Info("QueryWorkflow completed")
	return nil
}
