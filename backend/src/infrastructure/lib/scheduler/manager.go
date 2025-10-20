// scheduler/manager.go
package scheduler

import (
	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
)

type TaskManager interface {
	// 任务控制
	StartTask(taskID int) error
	StopTask(taskID int) error
	ReloadTasks() error
	StopAllTasks()

	// 任务管理
	AddTask(task *domainScheduledTask.ScheduledTask) error
	RemoveTask(taskID int) error
	UpdateTask(task *domainScheduledTask.ScheduledTask) error

	// 状态查询
	GetTaskStatus(taskID int) (bool, error)
	ListAllTasks() map[int]bool
}
