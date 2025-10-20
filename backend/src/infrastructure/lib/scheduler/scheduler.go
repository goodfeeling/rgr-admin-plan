// scheduler/scheduler.go
package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"sync"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	scheduleTaskConstants "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task/constants"
	domainTaskExecutionLog "github.com/gbrayhan/microservices-go/src/domain/sys/task_execution_log"
	"github.com/gbrayhan/microservices-go/src/infrastructure/lib/executor"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/scheduled_task"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/task_execution_log"
	wsHandler "github.com/gbrayhan/microservices-go/src/infrastructure/ws/handler/task_execution_log"
	"github.com/go-co-op/gocron"
	"go.uber.org/zap"
)

type TaskScheduler struct {
	scheduler            *gocron.Scheduler
	repo                 scheduled_task.IScheduledTaskRepository
	logger               *logger.Logger
	tasks                map[int]*gocron.Job
	executor             *executor.TaskExecutorManager
	mutex                sync.RWMutex
	taskWg               map[int]*sync.WaitGroup
	taskExecutionLogRepo task_execution_log.ITaskExecutionLogRepository
	wsHandler            *wsHandler.LogHandler
}

func NewTaskScheduler(
	repo scheduled_task.IScheduledTaskRepository,
	logger *logger.Logger,
	executor *executor.TaskExecutorManager,
	taskExecutionLogRepo task_execution_log.ITaskExecutionLogRepository,
) *TaskScheduler {
	// 创建支持秒级的调度器
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.SetMaxConcurrentJobs(10, gocron.RescheduleMode)

	return &TaskScheduler{
		scheduler:            scheduler,
		repo:                 repo,
		logger:               logger,
		tasks:                make(map[int]*gocron.Job),
		executor:             executor,
		taskExecutionLogRepo: taskExecutionLogRepo,
	}
}

func (s *TaskScheduler) SetWsHandler(handler *wsHandler.LogHandler) {
	s.wsHandler = handler
}

func (s *TaskScheduler) Start() {
	s.loadTasks()
	s.scheduler.StartAsync()
	s.logger.Info("Task scheduler started")
}

func (s *TaskScheduler) Stop() {
	s.scheduler.Stop()
	s.logger.Info("Task scheduler stopped")
}
func (s *TaskScheduler) loadTasks() {
	filters := domain.DataFilters{
		Matches: map[string][]string{
			"status": {scheduleTaskConstants.TaskStatusEnabled}, // 只加载启用的任务
		},
	}

	result, err := s.repo.SearchPaginated(filters)
	if err != nil {
		s.logger.Error("Failed to load tasks", zap.Error(err))
		return
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, task := range *result.Data {
		// 为每个任务创建本地副本以避免闭包问题
		taskCopy := task
		s.addTaskToScheduleInternal(&taskCopy)
	}
}

// 内部方法，不加锁
func (s *TaskScheduler) addTaskToScheduleInternal(task *domainScheduledTask.ScheduledTask) {
	// 创建一个闭包来捕获当前任务
	taskFunc := func() {
		s.executeTask(task)
	}

	// 使用gocron解析cron表达式并调度任务
	// 如果表达式包含6个字段（秒级），则使用 WithSeconds 选项
	var job *gocron.Job
	var err error

	// 检查cron表达式的字段数
	fields := strings.Fields(task.CronExpression)
	if len(fields) == 6 {
		// 6字段表达式，使用秒级解析
		job, err = s.scheduler.CronWithSeconds(task.CronExpression).Do(taskFunc)
	} else {
		// 标准5字段表达式
		job, err = s.scheduler.Cron(task.CronExpression).Do(taskFunc)
	}

	if err != nil {
		s.logger.Error("Failed to schedule task",
			zap.Int("task_id", task.ID),
			zap.String("task_name", task.TaskName),
			zap.String("cron", task.CronExpression),
			zap.Error(err))
		return
	}

	s.tasks[task.ID] = job
	s.logger.Info("Task scheduled",
		zap.Int("task_id", task.ID),
		zap.String("task_name", task.TaskName),
		zap.String("cron", task.CronExpression))
}

// 公共方法，加锁
func (s *TaskScheduler) addTaskToSchedule(task *domainScheduledTask.ScheduledTask) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.addTaskToScheduleInternal(task)
	return nil
}

func (s *TaskScheduler) executeTask(task *domainScheduledTask.ScheduledTask) {
	// 初始化 WaitGroup
	s.mutex.Lock()
	if s.taskWg == nil {
		s.taskWg = make(map[int]*sync.WaitGroup)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	s.taskWg[task.ID] = wg
	s.mutex.Unlock()

	// 确保在函数结束时调用 Done
	defer func() {
		s.mutex.Lock()
		if wg, exists := s.taskWg[task.ID]; exists {
			wg.Done()
			delete(s.taskWg, task.ID)
		}
		s.mutex.Unlock()
	}()

	s.logger.Info("Executing task",
		zap.Int("task_id", task.ID),
		zap.String("task_name", task.TaskName))

	// 更新任务状态为"运行中"
	now := time.Now()
	updateData := map[string]interface{}{
		"status":            scheduleTaskConstants.TaskStatusRunning, // "2" 表示运行中
		"last_execute_time": &now,
	}

	_, err := s.repo.Update(task.ID, updateData)
	if err != nil {
		s.logger.Error("Failed to update task status to running",
			zap.Int("task_id", task.ID),
			zap.Error(err))
	}
	// 执行任务
	err = s.executor.Execute(task)

	// 执行完成后，根据任务类型更新状态
	// 对于周期性任务，执行完成后恢复为"启用"状态
	// 对于一次性任务，可以设置为"已完成"状态
	finalStatus := scheduleTaskConstants.TaskStatusEnabled // 默认恢复为启用状态

	// 判断是否为一次性任务
	isOneTimeTask := task.ExecType == scheduleTaskConstants.TaskExecOnetime

	// 如果是一次性任务，则可以设置为已完成
	if isOneTimeTask && err == nil { // 假设有这样的字段标识一次性任务
		finalStatus = scheduleTaskConstants.TaskStatusCompleted // "5" 表示已完成
	}

	var taskLogData *domainTaskExecutionLog.TaskExecutionLog
	duration := int(time.Since(now).Seconds())
	if err != nil {
		s.logger.Error("Task execution failed",
			zap.Int("task_id", task.ID),
			zap.String("task_name", task.TaskName),
			zap.Error(err))

		// 执行失败，更新状态为"错误"
		finalStatus = scheduleTaskConstants.TaskStatusError // "4" 表示错误

		taskLogData, err = s.taskExecutionLogRepo.Create(&domainTaskExecutionLog.TaskExecutionLog{
			TaskID:          uint(task.ID),
			ExecuteTime:     now,
			ExecuteDuration: &duration,
			ErrorMessage:    err.Error(),
			ExecuteResult:   0,
		})
		if err != nil {
			s.logger.Error("Error creating task execution log",
				zap.Int("task_id", task.ID),
				zap.Error(err))
		}
	} else {
		s.logger.Info("Task executed successfully",
			zap.Int("task_id", task.ID),
			zap.String("task_name", task.TaskName))

		taskLogData, err = s.taskExecutionLogRepo.Create(&domainTaskExecutionLog.TaskExecutionLog{
			TaskID:          uint(task.ID),
			ExecuteTime:     now,
			ExecuteDuration: &duration,
			ErrorMessage:    "",
			ExecuteResult:   1,
		})
		if err != nil {
			s.logger.Error("Failed to create task execution log",
				zap.Int("task_id", task.ID),
				zap.String("task_name", task.TaskName),
				zap.Error(err))
		}
	}

	// update result
	updateData = map[string]interface{}{
		"status": finalStatus,
	}

	_, err = s.repo.Update(task.ID, updateData)
	if err != nil {
		s.logger.Error("Failed to update task execution result",
			zap.Int("task_id", task.ID),
			zap.Error(err))
	}

	// If it is a one-time task, remove it from the scheduler after execution is completed.
	if isOneTimeTask {
		s.mutex.Lock()
		if job, exists := s.tasks[task.ID]; exists {
			s.scheduler.RemoveByReference(job)
			delete(s.tasks, task.ID)
			s.logger.Info("One-time task removed from scheduler after execution",
				zap.Int("task_id", task.ID),
				zap.String("task_name", task.TaskName))
		}
		s.mutex.Unlock()
	}
	// Notify subscribers about the task execution log
	s.wsHandler.NotifyLogToTaskSubscribers(task.ID, taskLogData)
}

// ReloadTasks 重新加载所有任务（用于运行时刷新）
func (s *TaskScheduler) ReloadTasks() error {
	s.logger.Info("Reloading all tasks")

	// 停止所有当前任务
	s.StopAllTasks()

	// 清空任务映射
	s.mutex.Lock()
	s.tasks = make(map[int]*gocron.Job)
	s.mutex.Unlock()

	// 重新加载任务
	s.loadTasks()

	s.logger.Info("Tasks reloaded successfully")
	return nil
}

// StopAllTasks 停止所有任务
func (s *TaskScheduler) StopAllTasks() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, job := range s.tasks {
		s.scheduler.RemoveByReference(job)
	}
	s.logger.Info("All tasks stopped")
}

// StartTask 启动单个任务
func (s *TaskScheduler) StartTask(taskID int) error {
	// 先从数据库获取任务（在锁外面进行数据库操作）
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		s.logger.Error("Failed to get task by ID",
			zap.Int("task_id", taskID),
			zap.Error(err))
		return err
	}

	// 检查任务状态是否为启用
	if strconv.Itoa(task.Status) != scheduleTaskConstants.TaskStatusEnabled {
		s.logger.Warn("Task is not enabled, cannot start", zap.Int("task_id", taskID))
		return fmt.Errorf("task is not enabled")
	}

	// 然后获取锁进行调度操作
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果任务已经在运行，先停止它
	if job, exists := s.tasks[taskID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, taskID)
	}

	// 调度任务
	s.addTaskToScheduleInternal(task)
	return nil
}

// StopTask 停止单个任务
func (s *TaskScheduler) StopTask(taskID int) error {
	s.mutex.Lock()

	job, exists := s.tasks[taskID]
	if !exists {
		s.mutex.Unlock()
		s.logger.Warn("Task not found", zap.Int("task_id", taskID))
		return fmt.Errorf("task not found")
	}

	// 检查是否有 WaitGroup
	var wg *sync.WaitGroup
	if s.taskWg != nil {
		wg = s.taskWg[taskID]
	}

	// 先释放锁，避免死锁
	s.mutex.Unlock()

	// 如果任务正在运行且有 WaitGroup，等待完成
	if job.IsRunning() && wg != nil {
		s.logger.Info("Task is currently running, waiting for completion", zap.Int("task_id", taskID))
		wg.Wait() // 等待任务完成
		s.logger.Info("Task completed, now stopping", zap.Int("task_id", taskID))
	}

	// 重新获取锁来移除任务
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 再次检查任务是否存在（可能在等待期间已被删除）
	if job, exists := s.tasks[taskID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, taskID)
		s.logger.Info("Task stopped", zap.Int("task_id", taskID))
		return nil
	}

	s.logger.Warn("Task not found after waiting", zap.Int("task_id", taskID))
	return fmt.Errorf("task not found")
}

// AddTask 添加新任务
func (s *TaskScheduler) AddTask(task *domainScheduledTask.ScheduledTask) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果任务已存在，先移除
	if job, exists := s.tasks[task.ID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, task.ID)
	}

	// 如果任务是启用状态，则调度它
	if strconv.Itoa(task.Status) != scheduleTaskConstants.TaskStatusEnabled {
		return s.addTaskToSchedule(task)
	}

	s.logger.Info("Task added but not scheduled (disabled)", zap.Int("task_id", task.ID))
	return nil
}

// RemoveTask 移除任务
func (s *TaskScheduler) RemoveTask(taskID int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 停止任务调度
	if job, exists := s.tasks[taskID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, taskID)
	}

	s.logger.Info("Task removed", zap.Int("task_id", taskID))
	return nil
}

// UpdateTask 更新任务
func (s *TaskScheduler) UpdateTask(task *domainScheduledTask.ScheduledTask) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 停止现有任务
	if job, exists := s.tasks[task.ID]; exists {
		s.scheduler.RemoveByReference(job)
		delete(s.tasks, task.ID)
	}

	// 如果任务启用，则重新调度
	if strconv.Itoa(task.Status) != scheduleTaskConstants.TaskStatusEnabled {
		return s.addTaskToSchedule(task)
	}

	s.logger.Info("Task updated", zap.Int("task_id", task.ID))
	return nil
}

// GetTaskStatus 获取任务状态
func (s *TaskScheduler) GetTaskStatus(taskID int) (bool, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	_, exists := s.tasks[taskID]
	return exists, nil
}

// ListAllTasks 列出所有任务及其状态
func (s *TaskScheduler) ListAllTasks() map[int]bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	status := make(map[int]bool)
	for id, job := range s.tasks {
		status[id] = job.IsRunning()
	}
	return status
}
