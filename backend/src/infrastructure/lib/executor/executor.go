package executor

import (
	"fmt"

	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
)

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	Execute(task *domainScheduledTask.ScheduledTask) error
}

// TaskExecutorManager 任务执行管理器
type TaskExecutorManager struct {
	executors map[string]TaskExecutor
	logger    *logger.Logger
}

// NewTaskExecutorManager 创建任务执行管理器
func NewTaskExecutorManager(logger *logger.Logger) *TaskExecutorManager {
	return &TaskExecutorManager{
		executors: make(map[string]TaskExecutor),
		logger:    logger,
	}
}

// RegisterExecutor 注册任务执行器
func (m *TaskExecutorManager) RegisterExecutor(taskType string, executor TaskExecutor) {
	m.executors[taskType] = executor
}

// Execute 执行任务
func (m *TaskExecutorManager) Execute(task *domainScheduledTask.ScheduledTask) error {
	executor, exists := m.executors[task.TaskType]
	if !exists {
		return fmt.Errorf("no executor found for task type: %s", task.TaskType)
	}

	m.logger.Info("Executing task",
		zap.Int("task_id", task.ID),
		zap.String("task_name", task.TaskName),
		zap.String("task_type", task.TaskType))

	return executor.Execute(task)
}
