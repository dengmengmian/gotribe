package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"gotribe/internal/admin/index/dto"
	"gotribe/internal/admin/index/repository"
	"gotribe/internal/core/database"
)

// Service 仪表盘业务逻辑接口。
type Service interface {
	Dashboard(ctx context.Context, projectID string) (*dto.IndexResponse, error)
	CacheClear(ctx context.Context) error
}

type service struct {
	indexRepo *repository.Repository
	tx        *database.TransactionManager
	redis     redis.UniversalClient
}

// NewService 创建仪表盘服务实例。
func NewService(tx *database.TransactionManager, redis redis.UniversalClient) Service {
	return &service{
		indexRepo: repository.NewRepository(tx),
		tx:        tx,
		redis:     redis,
	}
}

// Dashboard 并发聚合仪表盘全量数据。
func (s *service) Dashboard(ctx context.Context, projectID string) (*dto.IndexResponse, error) {
	requestCtx := ctx
	var (
		stats          dto.Stats
		visitTrend     []dto.VisitPoint
		recentPosts    []dto.PostSummary
		recentComments []dto.CommentSummary
		popularPosts   []dto.PostSummary
		pending        dto.Pending
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error { var err error; stats, err = s.indexRepo.Stats(ctx, projectID); return err })
	g.Go(func() error { var err error; visitTrend, err = s.indexRepo.VisitTrend(ctx, projectID); return err })
	g.Go(func() error { var err error; recentPosts, err = s.indexRepo.RecentPosts(ctx, projectID, 5); return err })
	g.Go(func() error {
		var err error
		recentComments, err = s.indexRepo.RecentComments(ctx, projectID, 5)
		return err
	})
	g.Go(func() error {
		var err error
		popularPosts, err = s.indexRepo.PopularPosts(ctx, projectID, 5)
		return err
	})
	g.Go(func() error { var err error; pending, err = s.indexRepo.PendingCounts(ctx, projectID); return err })

	if err := g.Wait(); err != nil {
		return nil, err
	}

	sysStatus := dto.SystemStatus{DBStatus: "ok", RedisStatus: "ok"}
	if s.redis != nil {
		if _, err := s.redis.Ping(requestCtx).Result(); err != nil {
			sysStatus.RedisStatus = "error"
		}
	} else {
		sysStatus.RedisStatus = "未连接"
	}

	cacheStatus := s.cacheInfo(requestCtx)
	seoAlerts := s.indexRepo.SeoAlerts(requestCtx, projectID)

	return &dto.IndexResponse{
		Stats:          stats,
		VisitTrend:     visitTrend,
		RecentPosts:    ensureSlice(recentPosts),
		RecentComments: ensureSlice(recentComments),
		PopularPosts:   ensureSlice(popularPosts),
		Pending:        pending,
		SystemStatus:   sysStatus,
		CacheStatus:    cacheStatus,
		SeoAlerts:      seoAlerts,
	}, nil
}

func (s *service) cacheInfo(ctx context.Context) dto.CacheStatus {
	if s.redis == nil {
		return dto.CacheStatus{UsedMemory: "未知", UsedPercent: 0}
	}

	info, err := s.redis.Info(ctx, "memory").Result()
	if err != nil {
		return dto.CacheStatus{UsedMemory: "未知", UsedPercent: 0}
	}

	usedMem := parseRedisInfo(info, "used_memory_human")
	usedBytes := parseRedisInfoInt(info, "used_memory")
	maxBytes := parseRedisInfoInt(info, "maxmemory")

	percent := 0
	if usedBytes > 0 && maxBytes > 0 {
		percent = int(float64(usedBytes) / float64(maxBytes) * 100)
	} else if usedBytes > 0 {
		percent = 45 // 无 maxmemory 限制时给一个估算值
	}

	return dto.CacheStatus{
		UsedMemory:  usedMem,
		UsedPercent: percent,
	}
}

// CacheClear 清空当前 Redis DB 的缓存。
func (s *service) CacheClear(ctx context.Context) error {
	if s.redis == nil {
		return fmt.Errorf("redis 未连接")
	}
	return s.redis.FlushDB(ctx).Err()
}

func parseRedisInfo(info, key string) string {
	for _, line := range strings.Split(info, "\r\n") {
		if strings.HasPrefix(line, key+":") {
			return strings.TrimPrefix(line, key+":")
		}
	}
	return ""
}

func parseRedisInfoInt(info, key string) int64 {
	v := parseRedisInfo(info, key)
	if v == "" {
		return 0
	}
	n, _ := strconv.ParseInt(v, 10, 64)
	return n
}

func ensureSlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
