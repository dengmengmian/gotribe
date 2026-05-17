package service

import (
	"context"
	"mime/multipart"

	"gotribe/internal/admin/resource/dto"
	"gotribe/internal/admin/resource/repository"
	"gotribe/internal/admin/util"
	"gotribe/internal/admin/util/upload"
	"gotribe/internal/model"

	"go.uber.org/zap"
	"gotribe/internal/core/database"
)

// UploadConfig 上传配置
type UploadConfig struct {
	Provider  string
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	CDNDomain string
}

// Service 资源业务逻辑接口
type Service interface {
	Detail(ctx context.Context, id int64) (model.Resource, error)
	List(ctx context.Context, req *dto.ResourceListRequest) ([]*model.Resource, int64, error)
	Update(ctx context.Context, id int64, req *dto.CreateResourceRequest) error
	Create(ctx context.Context, resource *model.Resource) error
	Delete(ctx context.Context, ids []int64) error
	Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*model.Resource, *upload.UploadResource, error)
}

// service 资源业务逻辑实现
type service struct {
	resourceRepo *repository.Repository
	uploadCfg    UploadConfig
	log          *zap.SugaredLogger
}

// NewService 创建资源服务实例
func NewService(tx *database.TransactionManager, uploadCfg UploadConfig, log *zap.SugaredLogger) Service {
	return &service{
		resourceRepo: repository.NewRepository(tx),
		uploadCfg:    uploadCfg,
		log:          log,
	}
}

// Detail 根据ID获取资源
func (s *service) Detail(ctx context.Context, id int64) (model.Resource, error) {
	return s.resourceRepo.Detail(ctx, id)
}

// List 获取资源列表
func (s *service) List(ctx context.Context, req *dto.ResourceListRequest) ([]*model.Resource, int64, error) {
	return s.resourceRepo.List(ctx, req)
}

// Update 更新资源
func (s *service) Update(ctx context.Context, id int64, req *dto.CreateResourceRequest) error {
	oldResource, err := s.resourceRepo.Detail(ctx, id)
	if err != nil {
		return err
	}
	oldResource.Title = req.Title
	oldResource.Description = req.Description
	return s.resourceRepo.Update(ctx, &oldResource)
}

// Create 创建资源
func (s *service) Create(ctx context.Context, resource *model.Resource) error {
	return s.resourceRepo.Create(ctx, resource)
}

// Delete 批量删除资源
func (s *service) Delete(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		resource, err := s.resourceRepo.Detail(ctx, id)
		if err != nil {
			return err
		}

		if err := s.resourceRepo.Delete(ctx, id); err != nil {
			return err
		}

		// CDN 文件删除（best effort，不影响 DB 删除结果）
		provider := s.uploadCfg.Provider
		uploadSvc, err := upload.NewService(
			provider,
			s.uploadCfg.Endpoint,
			s.uploadCfg.AccessKey,
			s.uploadCfg.SecretKey,
			s.uploadCfg.Bucket,
			s.uploadCfg.Region,
		)
		if err != nil {
			s.log.Errorf("Failed to create upload service for deleting resource %d: %v", id, err)
			continue
		}
		if err := uploadSvc.DeleteFile(resource.Path); err != nil {
			s.log.Errorf("Failed to delete CDN file for resource %d: %v", id, err)
		}
	}
	return nil
}

// Upload 上传资源到 CDN 并入库
func (s *service) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*model.Resource, *upload.UploadResource, error) {
	provider := s.uploadCfg.Provider
	uploadSvc, err := upload.NewService(
		provider,
		s.uploadCfg.Endpoint,
		s.uploadCfg.AccessKey,
		s.uploadCfg.SecretKey,
		s.uploadCfg.Bucket,
		s.uploadCfg.Region,
	)
	if err != nil {
		return nil, nil, err
	}

	fileRes, err := uploadSvc.UploadFile(fileHeader)
	if err != nil {
		return nil, nil, err
	}

	fileType, err := util.FileUtil.GetFileType(fileHeader)
	if err != nil {
		// 文件类型检测失败不影响上传，记为未知类型
		fileType = 0
	}

	resource := &model.Resource{
		Title:         fileHeader.Filename,
		Path:          fileRes.Key,
		URL:           s.uploadCfg.CDNDomain,
		FileExtension: fileRes.FileExt,
		Size:          fileHeader.Size,
		FileType:      int64(fileType),
	}

	if err := s.resourceRepo.Create(ctx, resource); err != nil {
		// DB 入库失败，回滚 CDN 文件
		if delErr := uploadSvc.DeleteFile(fileRes.Key); delErr != nil {
			s.log.Errorf("Failed to rollback CDN file %s after DB failure: %v", fileRes.Key, delErr)
		}
		return nil, nil, err
	}

	return resource, &fileRes, nil
}
