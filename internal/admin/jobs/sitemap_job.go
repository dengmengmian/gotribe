package jobs

import (
	"context"

	"gotribe/internal/model"
	projectRepo "gotribe/internal/admin/project/repository"

	"github.com/douyacun/gositemap"
	"go.uber.org/zap"
	"gotribe/internal/core/database"
)

// SitemapJob 站点地图生成任务
type SitemapJob struct {
	*BaseJob
	tx *database.TransactionManager
}

// NewSitemapJob 创建站点地图任务
func NewSitemapJob(jobConfig JobConfig, tx *database.TransactionManager, log *zap.SugaredLogger) *SitemapJob {
	job := &SitemapJob{tx: tx}
	job.BaseJob = NewBaseJob(jobConfig, job.execute, log)
	return job
}

// execute 执行站点地图生成
func (j *SitemapJob) execute(ctx context.Context) error {
	j.Log().Info("Starting sitemap generation job")

	// 查出 project 信息
	projects, err := projectRepo.NewRepository(j.tx).GetProjectsBySitemap(ctx)
	if err != nil {
		return err
	}

	var posts []*model.Post
	st := gositemap.NewSiteMap()
	st.SetPretty(true)
	st.SetPublicPath("public")

	for _, project := range projects {
		// 检查上下文是否被取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 使用固定的文件名，确保每次都是覆盖
		st.SetFilename(project.Name + ".xml")

		if err := j.tx.DB(ctx).Model(&model.Post{}).Where("status = ? and type != ? and project_id = ?", 2, 2, project.Name).Find(&posts).Error; err != nil {
			j.Log().Errorf("Failed to query posts for project %s: %v", project.Name, err)
			continue
		}

		for _, post := range posts {
			url := gositemap.NewUrl()
			url.SetLoc(project.PostURL + post.Slug)
			url.SetLastmod(post.UpdatedAt)
			url.SetChangefreq(gositemap.Daily)
			url.SetPriority(1)
			st.AppendUrl(url)
		}
	}

	if _, err := st.Storage(); err != nil {
		return err
	}

	j.Log().Infof("Sitemap generation completed for %d projects", len(projects))
	return nil
}
