package request

// 本文件验证请求绑定层的语言识别和字段名处理逻辑。

import "testing"

// TestNormalizeLocale 验证请求绑定相关逻辑是否符合预期。
func TestNormalizeLocale(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "empty falls back later", raw: "", want: ""},
		{name: "zh cn", raw: "zh-CN", want: localeZH},
		{name: "zh quality", raw: "zh-CN,zh;q=0.9,en;q=0.8", want: localeZH},
		{name: "en us", raw: "en-US", want: localeEN},
		{name: "en quality", raw: "en-US,en;q=0.9,zh;q=0.8", want: localeEN},
		{name: "unknown falls back later", raw: "fr-FR", want: ""},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizeLocale(tc.raw); got != tc.want {
				t.Fatalf("normalizeLocale(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}

// TestTagName 验证请求绑定相关逻辑是否符合预期。
func TestTagName(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "json tag", raw: "user_id,omitempty", want: "user_id"},
		{name: "plain", raw: "page", want: "page"},
		{name: "dash ignored", raw: "-", want: ""},
		{name: "empty ignored", raw: "", want: ""},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tagName(tc.raw); got != tc.want {
				t.Fatalf("tagName(%q) = %q, want %q", tc.raw, got, tc.want)
			}
		})
	}
}
