// executor/http_executor.go
package executor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
)

// HTTPExecutor HTTP请求任务执行器
type HTTPExecutor struct {
	client *http.Client
	logger *logger.Logger
}

type HTTPParams struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Body    interface{}       `json:"body"`
	Timeout int               `json:"timeout"` // 秒
}

func NewHTTPExecutor(logger *logger.Logger) *HTTPExecutor {
	return &HTTPExecutor{
		client: &http.Client{},
		logger: logger,
	}
}

func (e *HTTPExecutor) Execute(task *domainScheduledTask.ScheduledTask) error {
	var params HTTPParams
	if err := json.Unmarshal(task.TaskParams, &params); err != nil {
		return fmt.Errorf("failed to parse HTTP params: %w", err)
	}

	// 设置超时
	timeout := 30 // 默认30秒
	if params.Timeout > 0 {
		timeout = params.Timeout
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 准备请求体
	var bodyBytes []byte
	if params.Body != nil {
		var err error
		bodyBytes, err = json.Marshal(params.Body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, params.Method, params.URL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	if params.Headers != nil {
		for key, value := range params.Headers {
			req.Header.Set(key, value)
		}
	}

	if params.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 发送请求
	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute HTTP request: %w", err)
	}
	defer resp.Body.Close()

	e.logger.Info("HTTP task executed successfully",
		zap.Int("task_id", task.ID),
		zap.String("url", params.URL),
		zap.Int("status_code", resp.StatusCode))

	return nil
}
