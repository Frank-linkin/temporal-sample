package ctxpropagation

import (
	"context"

	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

type (
	// contextKey is an unexported type used as key for items stored in the
	// Context object
	contextKey struct{}

	// propagator implements the custom context propagator
	logIdPropagator struct{}

	// Values is a struct holding values
	Values struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

// PropagateKey is the key used to store the value in the Context object
var LogIdPropagateKey = contextKey{}

// propagationKey is the key used by the propagator to pass values through the
// Temporal server headers
const propagationKey = "custom-header"

// NewContextPropagator returns a context propagator that propagates a set of
// string key-value pairs across a workflow
func NewContextPropagator() workflow.ContextPropagator {
	return &logIdPropagator{}
}

// Inject injects values from context into headers for propagation
func (s *logIdPropagator) Inject(ctx context.Context, writer workflow.HeaderWriter) error {
	logId := ctx.Value(LogIdPropagateKey)
	payload, err := converter.GetDefaultDataConverter().ToPayload(logId)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)
	return nil
}

// InjectFromWorkflow injects values from context into headers for propagation
func (s *logIdPropagator) InjectFromWorkflow(ctx workflow.Context, writer workflow.HeaderWriter) error {
	info:=workflow.GetInfo(ctx)
	var val string
	//如果是 rootWorkflow 就组装LogId
	//如果是 ChildWorkflow 就向Header写入
	if( info.ParentWorkflowExecution == nil) {
		val =info.WorkflowExecution.RunID+"_"+info.WorkflowExecution.ID
	}else{
		val= ctx.Value(LogIdPropagateKey).(string)
	}
	payload, err := converter.GetDefaultDataConverter().ToPayload(val)
	if err != nil {
		return err
	}
	writer.Set(propagationKey, payload)

	return nil
}

// Extract extracts values from headers and puts them into context
func (s *logIdPropagator) Extract(ctx context.Context, reader workflow.HeaderReader) (context.Context, error) {
	if logId, ok := reader.Get(propagationKey); ok {
		var values string
		if err := converter.GetDefaultDataConverter().FromPayload(logId, &values); err != nil {
			return ctx, nil
		}
		ctx = context.WithValue(ctx, LogIdPropagateKey, values)
	}

	return ctx, nil
}

// ExtractToWorkflow extracts values from headers and puts them into context
func (s *logIdPropagator) ExtractToWorkflow(ctx workflow.Context, reader workflow.HeaderReader) (workflow.Context, error) {
	if logId, ok := reader.Get(propagationKey); ok {
		var values string
		if err := converter.GetDefaultDataConverter().FromPayload(logId, &values); err != nil {
			return ctx, nil
		}
		ctx = workflow.WithValue(ctx, LogIdPropagateKey, values)
	}

	return ctx, nil
}
