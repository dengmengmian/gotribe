package jobs

import "time"

// JobConfig 定义单个任务的配置。
type JobConfig struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Schedule    string        `json:"schedule"`
	Enabled     bool          `json:"enabled"`
	Timeout     time.Duration `json:"timeout"`
	RetryCount  int           `json:"retry_count"`
}

// JobsConfig 定义定时任务集合配置。
type JobsConfig struct {
	Enabled bool                 `json:"enabled"`
	List    map[string]JobConfig `json:"list"`
}

// DefaultJobsConfig 默认任务配置
func DefaultJobsConfig() *JobsConfig {
	return &JobsConfig{
		Enabled: false,
		List: map[string]JobConfig{
			"sitemap": {
				Name:        "sitemap",
				Description: "生成站点地图",
				Schedule:    "@every 1h",
				Enabled:     false,
				Timeout:     5 * time.Minute,
				RetryCount:  3,
			},
			"example": {
				Name:        "example",
				Description: "示例任务",
				Schedule:    "@every 30s",
				Enabled:     false,
				Timeout:     1 * time.Minute,
				RetryCount:  1,
			},
		},
	}
}

// GetJobConfig 获取任务配置
func GetJobConfig(c *JobsConfig, jobName string) (JobConfig, bool) {
	if c == nil || c.List == nil {
		return JobConfig{}, false
	}
	conf, exists := c.List[jobName]
	return conf, exists
}

// IsJobEnabled 检查任务是否启用
func IsJobEnabled(c *JobsConfig, jobName string) bool {
	if c == nil || !c.Enabled {
		return false
	}

	conf, exists := GetJobConfig(c, jobName)
	return exists && conf.Enabled
}
