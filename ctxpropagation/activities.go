package ctxpropagation

import (
	"context"
	"fmt"
	"reflect"
)

func SampleActivity(ctx context.Context) (*Values, error) {
	ctxType:=reflect.TypeOf(ctx)
	str:=ctxType.String()
	name:=ctxType.Name()
	fmt.Println(str)
	fmt.Println(name)
	if val := ctx.Value(LogIdPropagateKey); val != nil {
		fmt.Println(val)
	}
	fmt.Println("this his it")

	return nil, nil
}

func SampleActivity2(ctx context.Context) (*Values, error) {
	//logger := activity.GetLogger(ctx)
	//if val := ctx.Value(PropagateKey); val != nil {
	//	vals := val.(Values)
	//	logger.Info("Activity PropagateKey", "key1", vals.Key,"value1",vals.Value)
	//	return &vals, nil
	//}
	fmt.Println("this his it")
	return nil, nil
}
