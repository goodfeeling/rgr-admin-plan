package executor

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
)

// ShellCommandParams shell 命令参数结构
type ShellCommandParams struct {
	Command string            `json:"command"`  // 要执行的命令
	Args    []string          `json:"args"`     // 命令参数
	Timeout int               `json:"timeout"`  // 超时时间（秒），默认300秒
	WorkDir string            `json:"work_dir"` // 工作目录
	Env     map[string]string `json:"env"`      // 环境变量
}

// ShellExecutor 执行 shell 脚本的执行器
func ShellExecutor(task *domainScheduledTask.ScheduledTask) error {
	// 解析任务参数
	params, err := parseShellParams(task.TaskParams)
	if err != nil {
		return fmt.Errorf("failed to parse shell params: %w", err)
	}

	// 如果没有指定命令，使用 TaskDescription
	if params.Command == "" {
		params.Command = task.TaskDescription
	}

	if params.Command == "" {
		return fmt.Errorf("no command specified for shell execution")
	}

	// 设置默认超时时间
	if params.Timeout <= 0 {
		params.Timeout = 300 // 默认5分钟
	}

	// 准备命令和参数
	var cmd *exec.Cmd
	var ctx context.Context
	var cancel context.CancelFunc

	// 创建带超时的上下文
	ctx, cancel = context.WithTimeout(context.Background(), time.Duration(params.Timeout)*time.Second)
	defer cancel()

	// 构建命令
	if len(params.Args) > 0 {
		cmd = exec.CommandContext(ctx, params.Command, params.Args...)
	} else {
		// 如果没有明确的参数，尝试分割命令字符串
		args := strings.Fields(params.Command)
		if len(args) > 1 {
			cmd = exec.CommandContext(ctx, args[0], args[1:]...)
		} else {
			cmd = exec.CommandContext(ctx, args[0])
		}
	}

	// 设置工作目录
	if params.WorkDir != "" {
		cmd.Dir = params.WorkDir
	}

	// 设置环境变量
	if len(params.Env) > 0 {
		env := append([]string{}, cmd.Env...) // 复制当前环境变量
		for key, value := range params.Env {
			env = append(env, fmt.Sprintf("%s=%s", key, value))
		}
		cmd.Env = env
	}

	// 执行命令
	output, err := cmd.CombinedOutput()

	// 记录执行结果
	if err != nil {
		return fmt.Errorf("shell command failed: %w, output: %s", err, string(output))
	}

	fmt.Printf("Shell command executed successfully:\n%s\n", string(output))
	return nil
}

// parseShellParams 解析 shell 参数
func parseShellParams(params interface{}) (*ShellCommandParams, error) {
	if params == nil {
		return &ShellCommandParams{}, nil
	}

	// 如果已经是 ShellCommandParams 类型
	if shellParams, ok := params.(*ShellCommandParams); ok {
		return shellParams, nil
	}

	// 如果是 JSON 格式
	if jsonData, ok := params.(json.RawMessage); ok {
		var shellParams ShellCommandParams
		if err := json.Unmarshal(jsonData, &shellParams); err != nil {
			return nil, fmt.Errorf("failed to unmarshal shell params: %w", err)
		}
		return &shellParams, nil
	}

	// 如果是字符串格式的 JSON
	if jsonStr, ok := params.(string); ok {
		var shellParams ShellCommandParams
		if err := json.Unmarshal([]byte(jsonStr), &shellParams); err != nil {
			return nil, fmt.Errorf("failed to unmarshal shell params from string: %w", err)
		}
		return &shellParams, nil
	}

	return &ShellCommandParams{}, nil
}
