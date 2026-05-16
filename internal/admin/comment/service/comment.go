package service

import (
	"context"

	"gotribe/internal/core/constant"
	"gotribe/internal/admin/comment/dto"
	"gotribe/internal/admin/comment/repository"
	"gotribe/internal/model"

	"gotribe/internal/core/database"
)

// Service 评论业务逻辑接口
type Service interface {
	List(ctx context.Context, req *dto.CommentListRequest) ([]*model.Comment, int64, error)
	Update(ctx context.Context, id int64, req *dto.UpdateCommentRequest) error
	Delete(ctx context.Context, id int64) error
}

// service 评论业务逻辑实现
type service struct {
	commentRepo *repository.Repository
}

// NewService 创建评论服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		commentRepo: repository.NewRepository(tx),
	}
}

// List 获取评论列表
func (s *service) List(ctx context.Context, req *dto.CommentListRequest) ([]*model.Comment, int64, error) {
	return s.commentRepo.List(ctx, req)
}

// Update 更新评论
// 如果请求未提供 Status，则根据当前状态自动切换：pending -> pass, pass -> pending
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdateCommentRequest) error {
	oldComment, err := s.commentRepo.Detail(ctx, id)
	if err != nil {
		return err
	}

	var reqStatus uint
	if req.Status != 0 {
		reqStatus = req.Status
	} else {
		if oldComment.Status == constant.AUDIT_STATUS_PENDING {
			reqStatus = constant.AUDIT_STATUS_PASS
		} else {
			reqStatus = constant.AUDIT_STATUS_PENDING
		}
	}
	oldComment.Status = reqStatus
	return s.commentRepo.Update(ctx, &oldComment)
}

// Delete 删除评论。
func (s *service) Delete(ctx context.Context, id int64) error {
	return s.commentRepo.Delete(ctx, id)
}
