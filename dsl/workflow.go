package dsl

import (
	"fmt"
	"reflect"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	// Workflow is the type used to express the workflow definition. Variables are a map of valuables. Variables can be
	// used as input to Activity.
	Workflow struct {
		Variables map[string]interface{}
		Root      Statement
	}

	// Statement is the building block of dsl workflow. A Statement can be a simple ActivityInvocation or it
	// could be a Sequence or Parallel.
	Statement struct {
		Childworkflow *ChildWorkflowInvocation
		Activity      *ActivityInvocation
		Sequence      *Sequence
		Parallel      *Parallel
	}

	// Sequence consist of a collection of Statements that runs in sequential.
	Sequence struct {
		Elements []*Statement
	}

	// Parallel can be a collection of Statements that runs in parallel.
	Parallel struct {
		Branches []*Statement
	}

	// ActivityInvocation is used to express invoking an Activity. The Arguments defined expected arguments as input to
	// the Activity, the result specify the name of variable that it will store the result as which can then be used as
	// arguments to subsequent ActivityInvocation.
	ActivityInvocation struct {
		Name      string
		Arguments []string
		Result    string
		Option    *workflow.ActivityOptions
	}

	ChildWorkflowInvocation struct {
		Name      string
		Arguments []string
		Result    string
		//考虑到Childworkflow执行的时候所有param来源于父workflow，不再从配置文件中获取信息，所以这里是Statement
		Root   *Statement
		Option *workflow.ChildWorkflowOptions
	}

	executable interface {
		execute(ctx workflow.Context, bindings map[string]interface{}) error
	}
)

// SimpleDSLWorkflow workflow definition
func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	bindings := make(map[string]interface{})
	for k, v := range dslWorkflow.Variables {
		//存储到参数表bingdings里
		bindings[k] = v
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	err := dslWorkflow.Root.execute(ctx, bindings)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return nil, err
	}

	logger.Info("DSL Workflow completed.")
	return nil, err
}

//ChildWorkflow第二个参数必须为root
//第三个参数固定为bindings，也就是参数表
//后序把这两个参数放到ctx里
func SimpleDSLWorkflowChild(ctx workflow.Context, root Statement, bindings map[string]interface{}, stu Student) (interface{}, error) {
	logger := workflow.GetLogger(ctx)

	err := root.execute(ctx, bindings)
	if err != nil {
		logger.Error("DSL Workflow failed.", "Error", err)
		return -1, err
	}

	//logger.Info("DSL Workflow completed.")
	fmt.Printf("\nDSL-Child Workflow completed.\n")
	return 100, err
}

func (b *Statement) execute(ctx workflow.Context, bindings map[string]interface{}) error {
	if b.Parallel != nil {
		err := b.Parallel.execute(ctx, bindings)
		if err != nil {
			return err
		}
		//一个StateMent只允许是一种[ChildWorkflow,parallel,Sequence,Activity其中的一种]
		return nil
	}
	if b.Sequence != nil {
		err := b.Sequence.execute(ctx, bindings)
		if err != nil {
			return err
		}
		//一个StateMent只允许是一种[ChildWorkflow,parallel,Sequence,Activity其中的一种]
		return nil
	}
	if b.Activity != nil {
		err := b.Activity.execute(ctx, bindings)
		if err != nil {
			return err
		}
		//一个StateMent只允许是一种[ChildWorkflow,parallel,Sequence,Activity其中的一种]
		return nil
	}

	if b.Childworkflow != nil {
		err := b.Childworkflow.execute(ctx, bindings)
		if err != nil {
			return err
		}
		//一个StateMent只允许是一种[ChildWorkflow,parallel,Sequence,Activity其中的一种]
		return nil
	}
	return nil
}

func (a ActivityInvocation) execute(ctx workflow.Context, bindings map[string]interface{}) error {
	if a.Name == "SampleActivity6" {
		fmt.Println(a.Name)
	}
	//没写配置就使用默认配置
	defaultOption := workflow.ActivityOptions{
		StartToCloseTimeout:    10 * time.Second,
		ScheduleToCloseTimeout: 10 * time.Second,
	}
	var ao *workflow.ActivityOptions
	if a.Option == nil {
		ao = &defaultOption
	} else {
		ao = a.Option
	}
	ctx = workflow.WithActivityOptions(ctx, *ao)

	parasList := makeInput(a.Arguments, bindings, false)
	parasList[0] = reflect.ValueOf(ctx)
	parasList[1] = reflect.ValueOf(a.Name)

	future := ExecuteActivityDSL(parasList)
	var result interface{}
	err := future.Get(ctx, &result)
	if err != nil {
		return err
	}
	//配置了Result的Activity，其返回值放入参数表
	if a.Result != "" {
		bindings[a.Result] = result
	}

	//考虑future作为value传入binding参数表
	//如果Activity1和Activity2并行执行，但Activity1需要使用Activity2的结果，但是是并行的，不知道是否执行完，所以需要使用future
	//但这样又感觉与Parallel和Sequence机制重复了
	//if a.Result!="" {
	//	bindings[a.Result] = future
	//}
	return nil
}

func (a ChildWorkflowInvocation) execute(ctx workflow.Context, bindings map[string]interface{}) error {
	if a.Root == nil {
		//抛出错误，此ChilWorkflow缺少Root
		fmt.Printf("childWorkflow %v 缺少root", a.Root)
	}

	paramsList := makeInput(a.Arguments, bindings, true)
	paramsList[0] = reflect.ValueOf(ctx)
	paramsList[1] = reflect.ValueOf(a.Name)
	paramsList[2] = reflect.ValueOf((*a.Root))

	var cwo *workflow.ChildWorkflowOptions
	defaultCwo := workflow.ChildWorkflowOptions{
		TaskQueue:  "dsl",
		WorkflowID: "DSL-Workflow",
	}
	//没有配置Option就加载默认的Option
	if a.Option == nil {
		cwo = &defaultCwo
	} else {
		cwo = a.Option
	}

	ctx = workflow.WithChildOptions(ctx, *cwo)
	childFuture := ExecuteChildWorkflowDSL(paramsList)
	var result interface{}
	err := childFuture.Get(ctx, &result)
	if err != nil {
		return err
	}

	//配置了Result的ChildWorkflow，其返回值放入参数表
	if a.Result != "" {
		bindings[a.Result] = result
	}
	return nil
}

func (s Sequence) execute(ctx workflow.Context, bindings map[string]interface{}) error {
	for _, a := range s.Elements {
		err := a.execute(ctx, bindings)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Parallel) execute(ctx workflow.Context, bindings map[string]interface{}) error {
	//
	// You can use the context passed in to activity as a way to cancel the activity like standard GO way.
	// Cancelling a parent context will cancel all the derived contexts as well.
	//

	// In the parallel block, we want to execute all of them in parallel and wait for all of them.
	// if one activity fails then we want to cancel all the rest of them as well.
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)
	var activityErr error
	for _, s := range p.Branches {
		f := executeAsync(s, childCtx, bindings)
		selector.AddFuture(f, func(f workflow.Future) {
			err := f.Get(ctx, nil)
			if err != nil {
				// cancel all pending activities
				cancelHandler()
				activityErr = err
			}
		})
	}

	for i := 0; i < len(p.Branches); i++ {
		selector.Select(ctx) // this will wait for one branch
		if activityErr != nil {
			return activityErr
		}
	}

	return nil
}

func executeAsync(exe executable, ctx workflow.Context, bindings map[string]interface{}) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, bindings)
		settable.Set(nil, err)
	})
	return future
}

func makeInput(argNames []string, argsMap map[string]interface{}, isChildWorkflow bool) []reflect.Value {
	//给ctx,name留出空间
	offset := 2
	if isChildWorkflow {
		offset = 4
	}
	args := make([]reflect.Value, len(argNames)+offset)

	if isChildWorkflow {
		bingdings := make(map[string]interface{})
		for _, argName := range argNames {
			bingdings[argName] = argsMap[argName]
		}
		args[3] = reflect.ValueOf(bingdings)
	}
	for i, arg := range argNames {
		args[i+offset] = reflect.ValueOf(argsMap[arg])
	}
	return args
}

func ExecuteActivityDSL(paramList []reflect.Value) workflow.Future {
	funcValue := reflect.ValueOf(workflow.ExecuteActivity)
	futureValue := funcValue.Call(paramList)
	future, _ := futureValue[0].Interface().(workflow.Future)
	return future
}

func ExecuteChildWorkflowDSL(paramList []reflect.Value) workflow.ChildWorkflowFuture {
	funcValue := reflect.ValueOf(workflow.ExecuteChildWorkflow)
	futureValue := funcValue.Call(paramList)
	future, _ := futureValue[0].Interface().(workflow.ChildWorkflowFuture)
	return future
}
