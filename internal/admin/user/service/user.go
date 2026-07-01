package service

import (
	"context"
	"fmt"
	"strings"

	"gotribe/internal/model"
	"gotribe/internal/admin/user/dto"
	"gotribe/internal/admin/user/repository"
	"gotribe/internal/core/cache"
	"gotribe/internal/core/util"

	"go.uber.org/zap"
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
	cache    *cache.Store
	log      *zap.SugaredLogger
}

// NewService 创建用户服务实例。
// cache 用于失效 ToC 端用户资料缓存（profile），可为 nil（Redis 不可用时降级为不失效）。
func NewService(tx *database.TransactionManager, store *cache.Store, log *zap.SugaredLogger) Service {
	return &service{
		userRepo: repository.NewRepository(tx),
		cache:    store,
		log:      log,
	}
}

// clearProfileCache 失效 ToC 端用户资料缓存（best effort）。
// ToC 的 profile key 使用数字字符串形式的 project_id（与 PostResponse.ProjectID 一致）。
func (s *service) clearProfileCache(ctx context.Context, projectID, userID int64) {
	if s.cache == nil || userID <= 0 {
		return
	}
	key := s.cache.ProfileKey(fmt.Sprintf("%d", projectID), userID)
	if err := s.cache.Delete(ctx, key); err != nil && s.log != nil {
		s.log.Warnf("清除用户 %d 资料缓存失败: %v", userID, err)
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

	if err := s.userRepo.Update(ctx, &oldUser); err != nil {
		return err
	}
	s.clearProfileCache(ctx, oldUser.ProjectID, oldUser.ID)
	return nil
}

// Delete 批量删除用户
func (s *service) Delete(ctx context.Context, ids []int64) error {
	// 删除前取出失效缓存所需字段；查询失败则不继续删除，避免留下无法清理的脏缓存。
	refs, err := s.userRepo.ListProfileRefsByIDs(ctx, ids)
	if err != nil {
		return err
	}
	if err := s.userRepo.Delete(ctx, ids); err != nil {
		return err
	}
	for _, ref := range refs {
		s.clearProfileCache(ctx, ref.ProjectID, ref.ID)
	}
	return nil
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
