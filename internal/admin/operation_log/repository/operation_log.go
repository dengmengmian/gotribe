package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gotribe/internal/admin/operation_log/dto"
	"gotribe/internal/model"

	"go.uber.org/zap"
	"gotribe/internal/core/database"
)

type Repository struct {
	tx  *database.TransactionManager
	log *zap.SugaredLogger
}

func NewRepository(tx *database.TransactionManager, log *zap.SugaredLogger) *Repository {
	return &Repository{tx: tx, log: log}
}

func buildOperationLogOrder(req *dto.OperationLogListRequest) string {
	sortByMap := map[string]string{
		"username":   "username",
		"ip":         "ip",
		"path":       "path",
		"status":     "status",
		"startTime":  "start_time",
		"start_time": "start_time",
		"timeCost":   "time_cost",
		"time_cost":  "time_cost",
		"desc":       "\"desc\"",
	}

	column, ok := sortByMap[strings.TrimSpace(req.SortBy)]
	if !ok {
		return "start_time DESC"
	}

	direction := "ASC"
	if strings.EqualFold(strings.TrimSpace(req.SortOrder), "desc") {
		direction = "DESC"
	}

	return fmt.Sprintf("%s %s", column, direction)
}

func (r *Repository) List(ctx context.Context, req *dto.OperationLogListRequest) ([]model.OperationLog, int64, error) {
	var list []model.OperationLog
	db := r.tx.DB(ctx).Model(&model.OperationLog{})

	username := strings.TrimSpace(req.Username)
	if username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	ip := strings.TrimSpace(req.Ip)
	if ip != "" {
		db = db.Where("ip LIKE ?", fmt.Sprintf("%%%s%%", ip))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}

	// 分页
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	db = db.Order(buildOperationLogOrder(req))
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}

	return list, total, err
}

func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	err := r.tx.DB(ctx).Where("id IN (?)", ids).Unscoped().Delete(&model.OperationLog{}).Error
	return err
}

func (r *Repository) SaveOperationLogChannel(olc <-chan *model.OperationLog) {
	logs := make([]model.OperationLog, 0, 10)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	flush := func() {
		if len(logs) == 0 {
			return
		}
		if err := r.tx.DB(context.Background()).Create(&logs).Error; err != nil {
			r.log.Errorf("批量写入操作日志失败: %v", err)
		}
		logs = logs[:0]
	}

	for {
		select {
		case log, ok := <-olc:
			if !ok {
				flush()
				return
			}
			logs = append(logs, *log)
			if len(logs) >= 10 {
				flush()
			}
		case <-ticker.C:
			flush()
		}
	}
}
