package helloworld

import (
	"context"
	"fmt"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
	"time"
)

// Workflow is a Hello World workflow definition.
type Param struct {
	Name string
	Age  int
}

func (a *Param) say() {
	fmt.Printf("\n[name=%v,age=%v]", a.Name, a.Age)
}
func Workflow(ctx workflow.Context, actor1 Param, actor2 Param) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 300 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger := workflow.GetLogger(ctx)
	logger.Info("HelloWorld workflow started", "name", actor1.Name)

	var result string
	err := workflow.ExecuteActivity(ctx, "prefix_Activity", "小张").Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}
	err = workflow.ExecuteActivity(ctx, "prefix_Activity2", "小明").Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("HelloWorld workflow completed.", "result", result)
	actor1.say()
	actor2.say()
	return actor1.Name, nil
}

type Act struct {
}

//1,2,2003,5,9,10,11,15,16
//1,5,9
//2,10,15
//3,11,16
func (a Act) Activity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	fmt.Printf("Activity1 is Running,name is %v ", name)
	fmt.Println("this is the 抛瓦")

	return "Hello " + name + "!", nil
}

func (a Act) Activity2(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Activity", "name", name)
	fmt.Printf("Activity2 is Running,name is %v ", name)
	fmt.Println("this is the 抛瓦2")
	return "Hello2 " + name + "!", nil
}
