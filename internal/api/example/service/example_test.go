package service

// 本文件展示 example service 层的单元测试写法：纯函数测试 + 接口 mock 测试。

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	examplemodel "gotribe/internal/model"
	exampleview "gotribe/internal/api/example/view"
	postview "gotribe/internal/api/post/view"
)

// mockPostReader 是 PostSummaryReader 的 mock 实现。
type mockPostReader struct {
	mock.Mock
}

func (m *mockPostReader) GetSummaries(ctx context.Context, projectID string, postIDs []string) (map[string]postview.Summary, error) {
	args := m.Called(ctx, projectID, postIDs)
	return args.Get(0).(map[string]postview.Summary), args.Error(1)
}

func TestNormalizePostIDs(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		want    []string
		wantErr bool
	}{
		{
			name:  "去重并保留顺序",
			input: []string{"  post-a  ", "post-b", "post-a", "post-c"},
			want:  []string{"post-a", "post-b", "post-c"},
		},
		{
			name:    "空数组返回错误",
			input:   []string{},
			wantErr: true,
		},
		{
			name:    "全空格返回错误",
			input:   []string{"  ", ""},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizePostIDs(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestService_resolvePosts(t *testing.T) {
	mockPosts := &mockPostReader{}
	svc := &Service{
		posts: mockPosts,
	}

	t.Run("成功解析文章引用", func(t *testing.T) {
		mockPosts.On("GetSummaries", mock.Anything, "1", []string{"post-a", "post-b"}).
			Return(map[string]postview.Summary{
				"post-a": {PostID: "post-a", Title: "Title A", Type: 1, Status: 1},
				"post-b": {PostID: "post-b", Title: "Title B", Type: 2, Status: 1},
			}, nil).Once()

		rows, refs, err := svc.resolvePosts(context.Background(), "1", "post-a", []string{"post-a", "post-b"}, 0, 42)

		assert.NoError(t, err)
		assert.Len(t, rows, 2)
		assert.Len(t, refs, 2)
		assert.Equal(t, "post-a", refs[0].PostID)
		assert.Equal(t, "Title A", refs[0].Title)
		mockPosts.AssertExpectations(t)
	})

	t.Run("主文章不在 post_ids 中返回错误", func(t *testing.T) {
		_, _, err := svc.resolvePosts(context.Background(), "1", "post-c", []string{"post-a", "post-b"}, 0, 42)
		assert.Error(t, err)
	})

	t.Run("主文章为空返回错误", func(t *testing.T) {
		_, _, err := svc.resolvePosts(context.Background(), "1", "", []string{"post-a"}, 0, 42)
		assert.Error(t, err)
	})
}

func TestToView(t *testing.T) {
	entity := examplemodel.Example{
		ExampleID: "ex_123",
		Name:      "Test Example",
		Status:    1,
	}
	posts := []exampleview.PostRef{
		{PostID: "p1", Title: "Primary", Type: 1, Status: 1},
		{PostID: "p2", Title: "Secondary", Type: 2, Status: 1},
	}

	view := toView(entity, posts)

	assert.Equal(t, "ex_123", view.ExampleID)
	assert.Equal(t, "Test Example", view.Name)
	assert.Len(t, view.Posts, 2)
}
