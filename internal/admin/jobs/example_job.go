package jobs

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// ExampleJob 示例任务
type ExampleJob struct {
	*BaseJob
}

// NewExampleJob 创建示例任务
func NewExampleJob(jobConfig JobConfig, log *zap.SugaredLogger) *ExampleJob {
	job := &ExampleJob{}
	job.BaseJob = NewBaseJob(jobConfig, job.execute, log)
	return job
}

// execute 执行示例任务
func (j *ExampleJob) execute(ctx context.Context) error {
	j.Log().Info("Starting example job")

	// 模拟一些工作
	select {
	case <-time.After(2 * time.Second):
		j.Log().Info("Example job completed")
		return nil
	case <-ctx.Done():
		j.Log().Info("Example job cancelled")
		return ctx.Err()
	}
}
