package interceptor

import (
	"fmt"
	commonpb "go.temporal.io/api/common/v1"
	"go.temporal.io/sdk/interceptors"
	"go.temporal.io/sdk/workflow"
	"sync"
)

var _ interceptors.WorkflowInterceptor = (*tracingInterceptor)(nil)
var _ interceptors.WorkflowInboundCallsInterceptor = (*tracingInboundCallsInterceptor)(nil)
var _ interceptors.WorkflowOutboundCallsInterceptor = (*tracingOutboundCallsInterceptor)(nil)

type tracingInterceptor struct {
	sync.Mutex
	// key is workflow id
	instances map[string]*tracingInboundCallsInterceptor
}

type tracingInboundCallsInterceptor struct {
	Next  interceptors.WorkflowInboundCallsInterceptor
	trace []string
}

type tracingOutboundCallsInterceptor struct {
	interceptors.WorkflowOutboundCallsInterceptorBase
	inbound *tracingInboundCallsInterceptor
}

func (t *tracingOutboundCallsInterceptor) Go(ctx workflow.Context, name string, f func(ctx workflow.Context)) workflow.Context {
	t.inbound.trace = append(t.inbound.trace, "Go")
	return t.Next.Go(ctx, name, f)
}

func NewTracingInterceptor() *tracingInterceptor {
	return &tracingInterceptor{instances: make(map[string]*tracingInboundCallsInterceptor)}
}

func (t *tracingInterceptor) InterceptWorkflow(info *workflow.Info, next interceptors.WorkflowInboundCallsInterceptor) interceptors.WorkflowInboundCallsInterceptor {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	result := &tracingInboundCallsInterceptor{next, nil}
	t.instances[info.WorkflowType.Name] = result
	return result
}

func (t *tracingInterceptor) GetTrace(workflowType string) []string {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	if i, ok := t.instances[workflowType]; ok {
		return i.trace
	}
	panic(fmt.Sprintf("Unknown workflowType %v, known types: %v", workflowType, t.instances))
}

func (t *tracingInboundCallsInterceptor) Init(outbound interceptors.WorkflowOutboundCallsInterceptor) error {
	return t.Next.Init(&tracingOutboundCallsInterceptor{
		interceptors.WorkflowOutboundCallsInterceptorBase{Next: outbound}, t})
}

func (t *tracingOutboundCallsInterceptor) ExecuteActivity(ctx workflow.Context, activityType string, args ...interface{}) workflow.Future {
	t.inbound.trace = append(t.inbound.trace, "ExecuteActivity")
	return t.Next.ExecuteActivity(ctx, activityType, args...)
}

func (t *tracingOutboundCallsInterceptor) ExecuteChildWorkflow(ctx workflow.Context, childWorkflowType string, args ...interface{}) workflow.ChildWorkflowFuture {
	t.inbound.trace = append(t.inbound.trace, "ExecuteChildWorkflow")
	return t.Next.ExecuteChildWorkflow(ctx, childWorkflowType, args...)
}

func (t *tracingInboundCallsInterceptor) ExecuteWorkflow(ctx workflow.Context, workflowType string, args ...interface{}) []interface{} {
	t.trace = append(t.trace, "ExecuteWorkflow begin")
	result := t.Next.ExecuteWorkflow(ctx, workflowType, args...)
	t.trace = append(t.trace, "ExecuteWorkflow end")
	return result
}

func (t *tracingInboundCallsInterceptor) ProcessSignal(ctx workflow.Context, signalName string, arg interface{}) error {
	t.trace = append(t.trace, "ProcessSignal")
	return t.Next.ProcessSignal(ctx, signalName, arg)
}

func (t *tracingInboundCallsInterceptor) HandleQuery(ctx workflow.Context, queryType string, args *commonpb.Payloads,
	handler func(*commonpb.Payloads) (*commonpb.Payloads, error)) (*commonpb.Payloads, error) {
	t.trace = append(t.trace, "HandleQuery begin")
	result, err := t.Next.HandleQuery(ctx, queryType, args, handler)
	t.trace = append(t.trace, "HandleQuery end")
	return result, err
}

var _ interceptors.WorkflowInterceptor = (*artificialInterceptor)(nil)
var _ interceptors.WorkflowInboundCallsInterceptor = (*artificialInboundCallsInterceptor)(nil)
var _ interceptors.WorkflowOutboundCallsInterceptor = (*artificialOutboundCallsInterceptor)(nil)

type artificialInterceptor struct {
	sync.Mutex
	// key is workflow id
	instances map[string]*artificialInboundCallsInterceptor
}

type artificialInboundCallsInterceptor struct {
	Next  interceptors.WorkflowInboundCallsInterceptor
	trace []string
}

type artificialOutboundCallsInterceptor struct {
	interceptors.WorkflowOutboundCallsInterceptorBase
	inbound *artificialInboundCallsInterceptor
}

func (t *artificialOutboundCallsInterceptor) Go(ctx workflow.Context, name string, f func(ctx workflow.Context)) workflow.Context {
	t.inbound.trace = append(t.inbound.trace, "Go")
	return t.Next.Go(ctx, name, f)
}

func NewArtificialInterceptor() *artificialInterceptor {
	return &artificialInterceptor{instances: make(map[string]*artificialInboundCallsInterceptor)}
}

func (t *artificialInterceptor) InterceptWorkflow(info *workflow.Info, next interceptors.WorkflowInboundCallsInterceptor) interceptors.WorkflowInboundCallsInterceptor {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	result := &artificialInboundCallsInterceptor{next, nil}
	t.instances[info.WorkflowType.Name] = result
	return result
}

func (t *artificialInterceptor) GetTrace(workflowType string) []string {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	if i, ok := t.instances[workflowType]; ok {
		return i.trace
	}
	panic(fmt.Sprintf("Unknown workflowType %v, known types: %v", workflowType, t.instances))
}

func (t *artificialInboundCallsInterceptor) Init(outbound interceptors.WorkflowOutboundCallsInterceptor) error {
	return t.Next.Init(&artificialOutboundCallsInterceptor{
		interceptors.WorkflowOutboundCallsInterceptorBase{Next: outbound}, t})
}

func (t *artificialOutboundCallsInterceptor) ExecuteActivity(ctx workflow.Context, activityType string, args ...interface{}) workflow.Future {
	t.inbound.trace = append(t.inbound.trace, "ExecuteActivity")
	return t.Next.ExecuteActivity(ctx, activityType, args...)
}

func (t *artificialOutboundCallsInterceptor) ExecuteChildWorkflow(ctx workflow.Context, childWorkflowType string, args ...interface{}) workflow.ChildWorkflowFuture {
	t.inbound.trace = append(t.inbound.trace, "ExecuteChildWorkflow")
	return t.Next.ExecuteChildWorkflow(ctx, childWorkflowType, args...)
}

func (t *artificialInboundCallsInterceptor) ExecuteWorkflow(ctx workflow.Context, workflowType string, args ...interface{}) []interface{} {
	t.trace = append(t.trace, "ExecuteWorkflow begin")
	result := t.Next.ExecuteWorkflow(ctx, workflowType, args...)
	t.trace = append(t.trace, "ExecuteWorkflow end")
	return result
}

func (t *artificialInboundCallsInterceptor) ProcessSignal(ctx workflow.Context, signalName string, arg interface{}) error {
	t.trace = append(t.trace, "ProcessSignal")
	return t.Next.ProcessSignal(ctx, signalName, arg)
}

func (t *artificialInboundCallsInterceptor) HandleQuery(ctx workflow.Context, queryType string, args *commonpb.Payloads,
	handler func(*commonpb.Payloads) (*commonpb.Payloads, error)) (*commonpb.Payloads, error) {
	t.trace = append(t.trace, "HandleQuery begin")
	result, err := t.Next.HandleQuery(ctx, queryType, args, handler)
	t.trace = append(t.trace, "HandleQuery end")
	return result, err
}

var _ interceptors.WorkflowInterceptor = (*signalInterceptor)(nil)
var _ interceptors.WorkflowInboundCallsInterceptor = (*signalInboundCallsInterceptor)(nil)
var _ interceptors.WorkflowOutboundCallsInterceptor = (*signalOutboundCallsInterceptor)(nil)

type signalInterceptor struct {
	ReturnErrorTimes int
	TimesInvoked     int
}

func NewSignalInterceptor() *signalInterceptor {
	return &signalInterceptor{}
}

type signalInboundCallsInterceptor struct {
	interceptors.WorkflowInboundCallsInterceptorBase
	control *signalInterceptor
}

func (t *signalInboundCallsInterceptor) ProcessSignal(ctx workflow.Context, signalName string, arg interface{}) error {
	t.control.TimesInvoked++
	if t.control.TimesInvoked <= t.control.ReturnErrorTimes {
		return fmt.Errorf("interceptor induced failure while processing signal %v", signalName)
	}
	return t.Next.ProcessSignal(ctx, signalName, arg)
}

type signalOutboundCallsInterceptor struct {
	interceptors.WorkflowOutboundCallsInterceptorBase
}

func (t *signalInterceptor) InterceptWorkflow(_ *workflow.Info, next interceptors.WorkflowInboundCallsInterceptor) interceptors.WorkflowInboundCallsInterceptor {
	result := &signalInboundCallsInterceptor{}
	result.Next = next
	result.control = t
	return result
}
