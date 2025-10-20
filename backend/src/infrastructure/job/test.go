package job

import (
	"context"
	"fmt"
	"time"

	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
)

const (
	FUNCTION_TYPE_CLEAN_UP_OLD_DATA = "clean_up_old_data"
)

func CleanOldData(scheduledTask *domainScheduledTask.ScheduledTask) error {
	fmt.Println("开始执行任务...")
	// 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 模拟耗时操作
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("任务执行完成")
		return nil
	case <-ctx.Done():
		fmt.Println("任务执行超时")
		return ctx.Err()
	}
}
