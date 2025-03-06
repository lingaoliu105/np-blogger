package scheduler

import (
	"context"
	"fmt"
	"np-blogger/internal/model"
	"np-blogger/internal/repository"
	"np-blogger/internal/service"
	"sync"
	"time"
)

// TaskScheduler 任务调度器
type TaskScheduler struct {
	githubClient *repository.GitHubClient
	blogService  *service.BlogService
	tasks        map[uint]*Task
	mutex        sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
}

// Task 定时任务
type Task struct {
	userID     uint
	interval   time.Duration
	lastRun    time.Time
	isRunning  bool
	repository *model.Repository
}

// NewTaskScheduler 创建新的任务调度器
func NewTaskScheduler(githubClient *repository.GitHubClient, blogService *service.BlogService) *TaskScheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &TaskScheduler{
		githubClient: githubClient,
		blogService:  blogService,
		tasks:        make(map[uint]*Task),
		ctx:          ctx,
		cancel:       cancel,
	}
}

// AddTask 添加定时任务
func (s *TaskScheduler) AddTask(userID uint, repository *model.Repository, interval time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 检查任务是否已存在
	if _, exists := s.tasks[userID]; exists {
		return fmt.Errorf("task already exists for user %d", userID)
	}

	// 创建新任务
	s.tasks[userID] = &Task{
		userID:     userID,
		interval:   interval,
		lastRun:    time.Now(),
		repository: repository,
	}

	// 启动任务协程
	go s.runTask(userID)

	return nil
}

// RemoveTask 移除定时任务
func (s *TaskScheduler) RemoveTask(userID uint) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.tasks, userID)
}

// Stop 停止所有任务
func (s *TaskScheduler) Stop() {
	s.cancel()
}

// runTask 运行任务
func (s *TaskScheduler) runTask(userID uint) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.mutex.RLock()
			task, exists := s.tasks[userID]
			s.mutex.RUnlock()

			if !exists {
				return
			}

			// 检查是否需要执行任务
			if time.Since(task.lastRun) >= task.interval && !task.isRunning {
				task.isRunning = true
				s.processGitCommits(task)
				task.lastRun = time.Now()
				task.isRunning = false
			}
		}
	}
}

// processGitCommits 处理Git提交记录
func (s *TaskScheduler) processGitCommits(task *Task) {
	// 获取最新的提交记录
	commits, err := s.githubClient.GetLatestCommits(
		task.repository.Owner,
		task.repository.Name,
		task.repository.LastProcessedCommit,
	)
	if err != nil {
		fmt.Printf("Error getting commits for user %d: %v\n", task.userID, err)
		return
	}

	// 如果没有新的提交，直接返回
	if len(commits) == 0 {
		return
	}

	// 更新最后处理的提交ID
	task.repository.LastProcessedCommit = commits[0].SHA

	// 处理每个提交
	for _, commit := range commits {
		// 生成并发布博客
		err = s.blogService.GenerateAndPublishBlog(task.userID, commit.Message)
		if err != nil {
			fmt.Printf("Error generating blog for commit %s: %v\n", commit.SHA, err)
			continue
		}
	}
}