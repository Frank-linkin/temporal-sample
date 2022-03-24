package dsl

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
)

type SampleActivities struct {
}

func (a *SampleActivities) SampleActivity1(ctx context.Context, input []int) (int, error) {
	/*name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil*/
	result1 := 0
	for _, num := range input {
		fmt.Printf("\nActivity1:num = %v ", num)
		result1 += num
	}
	return result1, nil
}

func (a *SampleActivities) SampleActivity2(ctx context.Context, input []int) (int, error) {
	/*name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil*/

	result2 := 1
	for _, num := range input {
		fmt.Printf("\nActivity2:num = %v ", num)
		result2 = result2 * num
	}
	return result2, nil
}

func (a *SampleActivities) SampleActivity3(ctx context.Context, input []int) (int, error) {
	/*	name := activity.GetInfo(ctx).ActivityType.Name
		fmt.Printf("Run %s with input %v \n", name, input)
		return "Result_" + name, nil*/

	result3 := 1
	for _, num := range input {
		fmt.Printf("\nActivity3:num = %v ", num)
		result3 = result3 * num
	}
	return result3, nil
}

func (a *SampleActivities) SampleActivity4(ctx context.Context, input []int) (int, error) {
	/*	name := activity.GetInfo(ctx).ActivityType.Name
		fmt.Printf("Run %s with input %v \n", name, input)
		return "Result_" + name, nil*/

	result4 := 0
	for _, num := range input {
		fmt.Printf("\nActivity4:num = %v ", num)
		result4 += num
	}
	fmt.Printf("\nresult4 = %v\n", result4)
	return result4, nil
}

func (a *SampleActivities) SampleActivity5(ctx context.Context, input []string) (string, error) {
	name := activity.GetInfo(ctx).ActivityType.Name
	fmt.Printf("Run %s with input %v \n", name, input)
	return "Result_" + name, nil
}
