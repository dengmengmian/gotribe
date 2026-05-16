// Package service implements example business object lifecycle management.
package service

// 本文件实现 example 模块的完整业务示例。

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gotribe/internal/core/database"
	"gotribe/internal/core/errs"
	exampledto "gotribe/internal/api/example/dto"
	examplemodel "gotribe/internal/model"
	examplerepo "gotribe/internal/api/example/repository"
	exampleview "gotribe/internal/api/example/view"
	postview "gotribe/internal/api/post/view"

	"gorm.io/gorm"
)

// PostSummaryReader 定义 example 模块依赖的文章读取能力契约。
type PostSummaryReader interface {
	GetSummaries(ctx context.Context, projectID string, postIDs []string) (map[string]postview.Summary, error)
}

// Service 负责封装示例业务单的业务逻辑。
type Service struct {
	repo  *examplerepo.Repository
	tx    *database.TransactionManager
	posts PostSummaryReader
}

// NewService 创建示例业务单服务实例。
func NewService(repo *examplerepo.Repository, tx *database.TransactionManager, posts PostSummaryReader) *Service {
	return &Service{
		repo:  repo,
		tx:    tx,
		posts: posts,
	}
}

// Create 创建一张示例业务单，并在事务中落主表和关联表。
func (s *Service) Create(ctx context.Context, projectID string, actor exampleview.Actor, req exampledto.CreateRequest) (*exampleview.Example, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errs.BadRequest("name is required", nil)
	}

	description := strings.TrimSpace(req.Description)
	status := int16(1)
	if req.Status != nil {
		status = *req.Status
	}

	projectIDNum, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		return nil, errs.BadRequest("invalid project_id", err)
	}

	postRows, postRefs, err := s.resolvePosts(ctx, projectID, req.PrimaryPostID, req.PostIDs, 0, int64(actor.UserID))
	if err != nil {
		return nil, err
	}

	entity := &examplemodel.Example{
		ExampleID:     newExampleID(),
		ProjectID:     projectIDNum,
		UserID:        int64(actor.UserID),
		OwnerUsername: actor.Username,
		OwnerNickname: actor.Nickname,
		Name:          name,
		Description:   description,
		Status:        uint8(status),
		PrimaryPostID: strings.TrimSpace(req.PrimaryPostID),
	}

	if err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.Create(txCtx, entity); err != nil {
			return errs.Internal("create example", err)
		}
		for i := range postRows {
			postRows[i].ExampleRecordID = entity.ID
		}
		if err := s.repo.CreatePosts(txCtx, postRows); err != nil {
			return errs.Internal("create example posts", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	view := toView(*entity, postRefs)
	return &view, nil
}

// List 查询当前用户的示例业务单列表。
func (s *Service) List(ctx context.Context, projectID string, actor exampleview.Actor, query exampledto.ListQuery) ([]exampleview.Example, database.Pagination, error) {
	filter := examplerepo.ListFilter{
		Page:    query.Page,
		PerPage: query.PerPage,
		Keyword: query.Keyword,
		Status:  query.Status,
	}

	items, total, err := s.repo.List(ctx, projectID, int64(actor.UserID), filter)
	if err != nil {
		return nil, database.Pagination{}, errs.Internal("list examples", err)
	}

	postMap, err := s.loadPostsByExampleRecordID(ctx, items)
	if err != nil {
		return nil, database.Pagination{}, errs.Internal("list example posts", err)
	}

	page, perPage := database.NormalizePagination(query.Page, query.PerPage)
	result := make([]exampleview.Example, 0, len(items))
	for _, item := range items {
		result = append(result, toView(item, postMap[item.ID]))
	}
	return result, database.Pagination{Page: page, PerPage: perPage, Total: total}, nil
}

// Detail 读取当前用户的一张示例业务单详情。
func (s *Service) Detail(ctx context.Context, projectID string, actor exampleview.Actor, exampleID string) (*exampleview.Example, error) {
	entity, err := s.repo.GetByExampleID(ctx, projectID, int64(actor.UserID), exampleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("example not found", err)
		}
		return nil, errs.Internal("get example", err)
	}

	postMap, err := s.repo.ListPostsByExampleRecordIDs(ctx, []int64{entity.ID})
	if err != nil {
		return nil, errs.Internal("get example posts", err)
	}

	view := toView(*entity, modelRowsToRefs(postMap[entity.ID]))
	return &view, nil
}

// Update 更新示例业务单，并在事务中完成主记录和关联记录替换。
func (s *Service) Update(ctx context.Context, projectID string, actor exampleview.Actor, exampleID string, req exampledto.UpdateRequest) (*exampleview.Example, error) {
	entity, err := s.repo.GetByExampleID(ctx, projectID, int64(actor.UserID), exampleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errs.NotFound("example not found", err)
		}
		return nil, errs.Internal("get example", err)
	}

	currentPostMap, err := s.repo.ListPostsByExampleRecordIDs(ctx, []int64{entity.ID})
	if err != nil {
		return nil, errs.Internal("get example posts", err)
	}

	finalName := entity.Name
	if req.Name != nil {
		finalName = strings.TrimSpace(*req.Name)
		if finalName == "" {
			return nil, errs.BadRequest("name cannot be empty", nil)
		}
	}

	finalDescription := entity.Description
	if req.Description != nil {
		finalDescription = strings.TrimSpace(*req.Description)
	}

	finalStatus := entity.Status
	if req.Status != nil {
		finalStatus = uint8(*req.Status)
	}

	finalPrimaryPostID := entity.PrimaryPostID
	if req.PrimaryPostID != nil {
		finalPrimaryPostID = strings.TrimSpace(*req.PrimaryPostID)
		if finalPrimaryPostID == "" {
			return nil, errs.BadRequest("primary_post_id cannot be empty", nil)
		}
	}

	finalPostIDs := postIDsFromModelRows(currentPostMap[entity.ID])
	if req.PostIDs != nil {
		finalPostIDs = *req.PostIDs
	}

	postRows, postRefs, err := s.resolvePosts(ctx, projectID, finalPrimaryPostID, finalPostIDs, entity.ID, int64(actor.UserID))
	if err != nil {
		return nil, err
	}

	if err := s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		updates := map[string]any{
			"name":            finalName,
			"description":     finalDescription,
			"status":          finalStatus,
			"primary_post_id": finalPrimaryPostID,
		}
		if err := s.repo.UpdateByID(txCtx, entity.ID, updates); err != nil {
			return errs.Internal("update example", err)
		}

		if err := s.repo.DeletePostsByExampleRecordID(txCtx, entity.ID); err != nil {
			return errs.Internal("delete example posts", err)
		}
		if err := s.repo.CreatePosts(txCtx, postRows); err != nil {
			return errs.Internal("create example posts", err)
		}
		return nil
	}); err != nil {
		return nil, err
	}

	entity.Name = finalName
	entity.Description = finalDescription
	entity.Status = finalStatus
	entity.PrimaryPostID = finalPrimaryPostID
	entity.UpdatedAt = time.Now().UTC()

	view := toView(*entity, postRefs)
	return &view, nil
}

// Delete 删除示例业务单，并在事务中级联删除关联记录。
func (s *Service) Delete(ctx context.Context, projectID string, actor exampleview.Actor, exampleID string) error {
	entity, err := s.repo.GetByExampleID(ctx, projectID, int64(actor.UserID), exampleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errs.NotFound("example not found", err)
		}
		return errs.Internal("get example", err)
	}

	return s.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
		if err := s.repo.DeletePostsByExampleRecordID(txCtx, entity.ID); err != nil {
			return errs.Internal("delete example posts", err)
		}
		if err := s.repo.DeleteByID(txCtx, entity.ID); err != nil {
			return errs.Internal("delete example", err)
		}
		return nil
	})
}

func (s *Service) resolvePosts(ctx context.Context, projectID, primaryPostID string, rawPostIDs []string, exampleRecordID int64, userID int64) ([]examplemodel.ExamplePost, []exampleview.PostRef, error) {
	postIDs, err := normalizePostIDs(rawPostIDs)
	if err != nil {
		return nil, nil, err
	}

	primaryPostID = strings.TrimSpace(primaryPostID)
	if primaryPostID == "" {
		return nil, nil, errs.BadRequest("primary_post_id is required", nil)
	}
	if !contains(postIDs, primaryPostID) {
		return nil, nil, errs.BadRequest("primary_post_id must be included in post_ids", nil)
	}

	summaries, err := s.posts.GetSummaries(ctx, projectID, postIDs)
	if err != nil {
		return nil, nil, err
	}
	if len(summaries) != len(postIDs) {
		for _, postID := range postIDs {
			if _, ok := summaries[postID]; !ok {
				return nil, nil, errs.NotFound(fmt.Sprintf("post not found: %s", postID), nil)
			}
		}
	}

	projectIDNum, err := strconv.ParseInt(projectID, 10, 64)
	if err != nil {
		return nil, nil, errs.BadRequest("invalid project_id", err)
	}

	rows := make([]examplemodel.ExamplePost, 0, len(postIDs))
	refs := make([]exampleview.PostRef, 0, len(postIDs))
	for index, postID := range postIDs {
		summary := summaries[postID]
		rows = append(rows, examplemodel.ExamplePost{
			ExampleRecordID: exampleRecordID,
			ProjectID:      projectIDNum,
			UserID:          userID,
			PostID:          summary.PostID,
			PostTitle:       summary.Title,
			PostType:        summary.Type,
			PostStatus:      summary.Status,
			Sort:            index + 1,
		})
		refs = append(refs, exampleview.PostRef{
			PostID: summary.PostID,
			Title:  summary.Title,
			Type:   summary.Type,
			Status: summary.Status,
		})
	}
	return rows, refs, nil
}

func (s *Service) loadPostsByExampleRecordID(ctx context.Context, items []examplemodel.Example) (map[int64][]exampleview.PostRef, error) {
	recordIDs := make([]int64, 0, len(items))
	for _, item := range items {
		recordIDs = append(recordIDs, item.ID)
	}

	rowsByExampleID, err := s.repo.ListPostsByExampleRecordIDs(ctx, recordIDs)
	if err != nil {
		return nil, err
	}

	result := make(map[int64][]exampleview.PostRef, len(rowsByExampleID))
	for exampleRecordID, rows := range rowsByExampleID {
		refs := make([]exampleview.PostRef, 0, len(rows))
		for _, row := range rows {
			refs = append(refs, exampleview.PostRef{
				PostID: row.PostID,
				Title:  row.PostTitle,
				Type:   row.PostType,
				Status: row.PostStatus,
			})
		}
		result[exampleRecordID] = refs
	}
	return result, nil
}

func toView(entity examplemodel.Example, posts []exampleview.PostRef) exampleview.Example {
	primary := exampleview.PostRef{PostID: entity.PrimaryPostID}
	for _, post := range posts {
		if post.PostID == entity.PrimaryPostID {
			primary = post
			break
		}
	}

	return exampleview.Example{
		ExampleID:   entity.ExampleID,
		Name:        entity.Name,
		Description: entity.Description,
		Status:      int16(entity.Status),
		Owner: exampleview.Owner{
			UserID:   int64(entity.UserID),
			Username: entity.OwnerUsername,
			Nickname: entity.OwnerNickname,
		},
		PrimaryPost: primary,
		Posts:       posts,
		CreatedAt:   entity.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   entity.UpdatedAt.Format(time.RFC3339),
	}
}

func normalizePostIDs(raw []string) ([]string, error) {
	if len(raw) == 0 {
		return nil, errs.BadRequest("post_ids is required", nil)
	}

	result := make([]string, 0, len(raw))
	seen := make(map[string]struct{}, len(raw))
	for _, item := range raw {
		postID := strings.TrimSpace(item)
		if postID == "" {
			return nil, errs.BadRequest("post_ids cannot contain empty value", nil)
		}
		if _, exists := seen[postID]; exists {
			continue
		}
		seen[postID] = struct{}{}
		result = append(result, postID)
	}
	if len(result) == 0 {
		return nil, errs.BadRequest("post_ids is required", nil)
	}
	return result, nil
}

func postIDsFromRows(rows []exampleview.PostRef) []string {
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.PostID)
	}
	return result
}

func postIDsFromModelRows(rows []examplemodel.ExamplePost) []string {
	result := make([]string, 0, len(rows))
	for _, row := range rows {
		result = append(result, row.PostID)
	}
	return result
}

func modelRowsToRefs(rows []examplemodel.ExamplePost) []exampleview.PostRef {
	result := make([]exampleview.PostRef, 0, len(rows))
	for _, row := range rows {
		result = append(result, exampleview.PostRef{
			PostID: row.PostID,
			Title:  row.PostTitle,
			Type:   row.PostType,
			Status: row.PostStatus,
		})
	}
	return result
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}

func newExampleID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("ex_%d", time.Now().UnixNano())
	}
	return "ex_" + hex.EncodeToString(buf)
}
