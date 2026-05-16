package service

import (
	"context"
	"strings"

	"gotribe/internal/model"
	"gotribe/internal/admin/user/dto"
	"gotribe/internal/admin/user/repository"
	"gotribe/internal/core/util"

	"gotribe/internal/core/database"
)

// Service 用户业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.User, error)
	List(ctx context.Context, req *dto.UserListRequest) ([]*model.User, int64, error)
	Create(ctx context.Context, req *dto.CreateUserRequest) error
	Update(ctx context.Context, id int64, req *dto.UpdateUserRequest) error
	Delete(ctx context.Context, ids []int64) error
	Search(ctx context.Context, nickname string) ([]*model.User, error)
}

// service 用户业务逻辑实现
type service struct {
	userRepo *repository.Repository
}

// NewService 创建用户服务实例
func NewService(tx *database.TransactionManager) Service {
	return &service{
		userRepo: repository.NewRepository(tx),
	}
}

// Detail 根据ID获取用户
func (s *service) Detail(ctx context.Context, id int64) (model.User, error) {
	return s.userRepo.Detail(ctx, id)
}

// List 获取用户列表
func (s *service) List(ctx context.Context, req *dto.UserListRequest) ([]*model.User, int64, error) {
	return s.userRepo.List(ctx, req)
}

// Create 创建用户
func (s *service) Create(ctx context.Context, req *dto.CreateUserRequest) error {
	encryptedPwd, err := utils.PasswordUtil.Encrypt(req.Password)
	if err != nil {
		return err
	}

	user := model.User{
		Username:  req.Username,
		Nickname:  req.Nickname,
		AvatarURL: req.AvatarURL,
		Phone:     optionalString(req.Phone),
		Email:     optionalString(req.Email),
		ProjectID: req.ProjectID,
		Password:  encryptedPwd,
	}
	return s.userRepo.Create(ctx, &user)
}

// Update 更新用户
func (s *service) Update(ctx context.Context, id int64, req *dto.UpdateUserRequest) error {
	oldUser, err := s.userRepo.Detail(ctx, id)
	if err != nil {
		return err
	}

	oldUser.Nickname = req.Nickname
	oldUser.AvatarURL = req.AvatarURL
	oldUser.Phone = optionalString(req.Phone)
	oldUser.Email = optionalString(req.Email)
	if len(req.Password) > 0 {
		newPassword, _ := utils.PasswordUtil.Encrypt(req.Password)
		oldUser.Password = newPassword
	}

	return s.userRepo.Update(ctx, &oldUser)
}

// Delete 批量删除用户
func (s *service) Delete(ctx context.Context, ids []int64) error {
	return s.userRepo.Delete(ctx, ids)
}

// Search 根据昵称搜索用户
func (s *service) Search(ctx context.Context, nickname string) ([]*model.User, error) {
	return s.userRepo.Search(ctx, nickname)
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}
