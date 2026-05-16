package repository

import (
	"context"
	"fmt"
	"time"

	"gotribe/internal/admin/point/dto"
	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	"gotribe/internal/core/util"
	"gotribe/internal/model"
)

type Repository struct {
	tx *database.TransactionManager
}

// NewPointLogRepository 创建积分仓库实例
func NewRepository(tx *database.TransactionManager) *Repository {
	return &Repository{tx: tx}
}

// List 获取积分列表
func (r *Repository) List(ctx context.Context, req *dto.PointLogListRequest) ([]*model.PointLog, int64, error) {
	var list []*model.PointLog
	db := r.tx.DB(ctx).Model(&model.PointLog{})

	if req.ProjectID > 0 {
		db = db.Where("project_id = ?", req.ProjectID)
	}
	if req.UserID > 0 {
		db = db.Where("user_id =  ?", req.UserID)
	}
	if !utils.IsEmpty(req.Nickname) {
		// 查出用户 ID。再用用户 ID 去筛选
		var user model.User
		if result := r.tx.DB(ctx).Model(&model.User{}).Where("nickname like ?", fmt.Sprintf("%%%s%%", req.Nickname)).First(&user); result.Error != nil {
			return nil, 0, errs.NotFoundWithKey(errs.MsgUserNotFound, nil, nil)
		}
		db = db.Where("user_id = ?", user.ID)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	db = db.Order("created_at DESC")

	page, perPage := database.NormalizePagination(int(req.PageNum), int(req.PageSize))
	err = db.Offset((page - 1) * perPage).Limit(perPage).Find(&list).Error

	list, err = r.GetPointLogOther(ctx, list)
	return list, total, err
}

// GetPointLogOther 获取其他信息
func (r *Repository) GetPointLogOther(ctx context.Context, pointLogs []*model.PointLog) ([]*model.PointLog, error) {
	if len(pointLogs) == 0 {
		return pointLogs, nil
	}

	// 收集所有 UserID
	userIDs := make([]int64, 0, len(pointLogs))
	for _, m := range pointLogs {
		if m.UserID > 0 {
			userIDs = append(userIDs, m.UserID)
		}
	}

	if len(userIDs) == 0 {
		return pointLogs, nil
	}

	// 批量查询
	var users []*model.User
	if err := r.tx.DB(ctx).Where("id IN ?", userIDs).Find(&users).Error; err != nil {
		return pointLogs, err
	}

	// 建立映射
	userMap := make(map[int64]*model.User)
	for _, user := range users {
		userMap[user.ID] = user
	}

	// 赋值
	for _, m := range pointLogs {
		if user, ok := userMap[m.UserID]; ok {
			m.User = user
		}
	}
	return pointLogs, nil
}

// Detail 获取单条积分记录。
func (r *Repository) Detail(ctx context.Context, id int64) (*model.PointLog, error) {
	var log model.PointLog
	if err := r.tx.DB(ctx).Where("id = ?", id).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

// Update 更新积分记录。
func (r *Repository) Update(ctx context.Context, id int64, log *model.PointLog) error {
	return r.tx.DB(ctx).Model(&model.PointLog{}).Where("id = ?", id).Updates(log).Error
}

// Delete 删除积分记录（批量）。
func (r *Repository) Delete(ctx context.Context, ids []int64) error {
	return r.tx.DB(ctx).Where("id IN ?", ids).Delete(&model.PointLog{}).Error
}

// Create 创建积分（事务保护双表写入）。
func (r *Repository) Create(ctx context.Context, req *dto.CreatePointLogRequest) error {
	pointsCents := int64(req.Point * 100)
	db := r.tx.DB(ctx)

	pointLog := &model.PointLog{
		UserID:    req.UserID,
		Type:      "admin",
		Reason:    "后台添加",
		EventID:   "0",
		Points:    pointsCents,
		ProjectID: req.ProjectID,
	}
	if err := db.Create(pointLog).Error; err != nil {
		return err
	}

	userPoint := &model.PointAvailable{
		ProjectID:      req.ProjectID,
		UserID:         req.UserID,
		Points:         pointsCents,
		PointsLogID:    pointLog.ID,
		ExpirationDate: time.Now().AddDate(1, 0, 0),
	}
	return db.Create(userPoint).Error
}
