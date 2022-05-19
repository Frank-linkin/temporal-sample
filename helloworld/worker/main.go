package main

import (
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"time"

	"github.com/temporalio/samples-go/helloworld"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.NewClient(client.Options{
		HostPort: "106.13.193.55:7233",
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world", worker.Options{
		WorkerStopTimeout: 10 * time.Second,
	})

	w.RegisterWorkflow(helloworld.Workflow)
	var a helloworld.Act
	w.RegisterActivityWithOptions(&a, activity.RegisterOptions{
		Name: "prefix_",
	})
	//w.RegisterActivity(helloworld.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

}

//func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
//	if parent == nil {
//		panic("cannot create context from nil parent")
//	}
//	// 如果parent的截止时间更早，直接返回一个cancelCtx即可
//	if cur, ok := parent.Deadline(); ok && cur.Before(d) {
//		return WithCancel(parent)
//	}
//	c := &timerCtx{
//		cancelCtx: newCancelCtx(parent),
//		deadline:  d,
//	}
//	// 建立新建context与可取消context祖先节点的取消关联关系
//	propagateCancel(parent, c)
//	dur := time.Until(d)
//	if dur <= 0 { //当前时间已经超过了截止时间，直接cancel
//		c.cancel(true, DeadlineExceeded)
//		return c, func() { c.cancel(false, Canceled) }
//	}
//	c.mu.Lock()
//	defer c.mu.Unlock()
//	if c.err == nil {
//		// 设置一个定时器，到截止时间后取消
//		c.timer = time.AfterFunc(dur, func() {
//			c.cancel(true, DeadlineExceeded)
//		})
//	}
//	return c, func() { c.cancel(true, Canceled) }
//}
