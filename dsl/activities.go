package dsl

import (
	"context"
	"fmt"
	"github.com/mitchellh/mapstructure"
)

type SampleActivities struct {
}

//Activity参数表中变量顺序需要与用户在yaml文件中定义的完全一致
func (a *SampleActivities) SampleActivity1(ctx context.Context, arg1 int, arg2 int) (int, error) {
	result1 := arg1 + arg2
	fmt.Println("Activity1 is done")
	return result1, nil
}

func (a *SampleActivities) SampleActivity2(ctx context.Context, arg1 int, result1 int) (int, error) {
	result2 := result1 * arg1
	fmt.Println("Activity2 is done")
	return result2, nil
}

func (a *SampleActivities) SampleActivity3(ctx context.Context, arg2 int, arg3 int, result1 int) (int, error) {
	result3 := arg2 * arg3 * result1
	fmt.Println("Activity3 is done")
	return result3, nil
}

func (a *SampleActivities) SampleActivity4(ctx context.Context, result1 int, result2 int, result3 int) (int, error) {
	result4 := result1 + result2 + result3
	fmt.Printf("\ntotalResult=%v\n", result4)
	return result4, nil
}

type Student struct {
	Age  int    `json:"age"`
	Name string `json:"name"`
}

//没有返回值的Activity
func (a *SampleActivities) SampleActivity5(ctx context.Context) (int, error) {
	fmt.Printf("\nActivity5 is running\n")
	return 3, nil
}

//接受和返回自定义值的Activity
//当参数为对象时，无法判断它是对象还是map，需要用户自己转化成对象
func (a *SampleActivities) SampleActivity6(ctx context.Context, stu map[string]interface{}) (*Student, error) {
	var myStu Student
	if err := mapstructure.Decode(stu, &myStu); err != nil {
		return nil, err
	}
	fmt.Printf("\n stu.Name=%v,stu.age=%v \n", myStu.Name, myStu.Age)
	fmt.Println("111")
	return &myStu, nil
}
