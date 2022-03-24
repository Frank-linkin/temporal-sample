package dsl

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/workflow"
)

type (
	// Workflow is the type used to express the workflow definition. Variables are a map of valuables. Variables can be
	// used as input to Activity.
	Workflow struct {
		Variables map[string]int
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
		//SfChildWorkflowOpt{}配置文件的路径
		//OptionPath string
	}

	ChildWorkflowInvocation struct {
		Name      string
		Arguments []string
		Result    string
		//考虑到Childworkflow执行的时候所有param来源于父workflow，不再从配置文件中获取信息，所以这里是Statement
		Root Statement
		//SfChildWorkflowOpt{}配置文件的路径
		//OptionPath string
	}

	executable interface {
		execute(ctx workflow.Context, bindings map[string]int) error
	}
)

// SimpleDSLWorkflow workflow definition
func SimpleDSLWorkflow(ctx workflow.Context, dslWorkflow Workflow) ([]byte, error) {
	bindings := make(map[string]int)
	for k, v := range dslWorkflow.Variables {
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

func SimpleDSLWorkflowChild(ctx workflow.Context, root Statement, bindings map[string]int) (int, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
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

func (b *Statement) execute(ctx workflow.Context, bindings map[string]int) error {
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

func (a ActivityInvocation) execute(ctx workflow.Context, bindings map[string]int) error {
	inputParam := makeInput(a.Arguments, bindings)
	var result int
	err := workflow.ExecuteActivity(ctx, a.Name, inputParam).Get(ctx, &result)
	if err != nil {
		return err
	}
	if a.Result != "" {
		bindings[a.Result] = result
	}
	return nil
}

func (a ChildWorkflowInvocation) execute(ctx workflow.Context, bindings map[string]int) error {
	//inputParam := makeInput(a.Arguments, bindings)
	var result int
	//配置的Option单独写yaml文件
	//第三个参数必须是a.workflow
	cwo := workflow.ChildWorkflowOptions{
		TaskQueue:  "dsl",
		WorkflowID: "DSL-Workflow",
	}
	ctx = workflow.WithChildOptions(ctx, cwo)
	err := workflow.ExecuteChildWorkflow(ctx, a.Name, a.Root, bindings).Get(ctx, &result)
	if err != nil {
		return err
	}
	if a.Result != "" {
		bindings[a.Result] = result
	}
	return nil
}

func (s Sequence) execute(ctx workflow.Context, bindings map[string]int) error {
	for _, a := range s.Elements {
		err := a.execute(ctx, bindings)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p Parallel) execute(ctx workflow.Context, bindings map[string]int) error {
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

func executeAsync(exe executable, ctx workflow.Context, bindings map[string]int) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := exe.execute(ctx, bindings)
		settable.Set(nil, err)
	})
	return future
}

func makeInput(argNames []string, argsMap map[string]int) []int {
	var args []int
	for _, arg := range argNames {
		args = append(args, argsMap[arg])
	}
	return args
}
