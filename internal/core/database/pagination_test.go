package database

// 本文件验证分页参数规范化逻辑。

import "testing"

// TestNormalizePagination 验证数据库相关逻辑是否符合预期。
func TestNormalizePagination(t *testing.T) {
	page, perPage := NormalizePagination(0, 999)
	if page != 1 {
		t.Fatalf("page = %d, want 1", page)
	}
	if perPage != 100 {
		t.Fatalf("perPage = %d, want 100", perPage)
	}
}
