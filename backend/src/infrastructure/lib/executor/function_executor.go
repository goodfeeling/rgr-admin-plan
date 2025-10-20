// executor/function_executor.go
package executor

import (
	"encoding/json"
	"fmt"

	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
)

// FunctionExecutor 函数任务执行器
type FunctionExecutor struct {
	functions map[string]func(*domainScheduledTask.ScheduledTask) error
	logger    *logger.Logger
}

type FunctionParams struct {
	FunctionName string                 `json:"function_name"`
	Params       map[string]interface{} `json:"params"`
}

func NewFunctionExecutor(logger *logger.Logger) *FunctionExecutor {
	return &FunctionExecutor{
		functions: make(map[string]func(*domainScheduledTask.ScheduledTask) error),
		logger:    logger,
	}
}

// RegisterFunction 注册可执行函数
func (e *FunctionExecutor) RegisterFunction(name string, fn func(*domainScheduledTask.ScheduledTask) error) {
	e.functions[name] = fn
}

func (e *FunctionExecutor) Execute(task *domainScheduledTask.ScheduledTask) error {
	var params FunctionParams
	if err := json.Unmarshal(task.TaskParams, &params); err != nil {
		return fmt.Errorf("failed to parse function params: %w", err)
	}

	function, exists := e.functions[params.FunctionName]
	if !exists {
		return fmt.Errorf("function not found: %s", params.FunctionName)
	}

	e.logger.Info("Executing function task",
		zap.Int("task_id", task.ID),
		zap.String("function_name", params.FunctionName))

	return function(task)
}
