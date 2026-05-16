package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"gotribe/internal/admin/ai/dto"
	"gotribe/internal/core/config"
	"gotribe/internal/core/errs"
)

// Service 提供通用 AI 生成能力。
type Service interface {
	Generate(ctx context.Context, req *dto.GenerateRequest) (*dto.GenerateResponse, error)
}

type service struct {
	cfg    config.AIConfig
	client *http.Client
}

// NewService 创建 AI 服务。
func NewService(cfg config.AIConfig) Service {
	timeout := time.Duration(cfg.TimeoutSeconds) * time.Second
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &service{
		cfg: cfg,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *service) Generate(ctx context.Context, req *dto.GenerateRequest) (*dto.GenerateResponse, error) {
	if strings.TrimSpace(s.cfg.APIKey) == "" {
		return nil, errs.ServiceUnavailable("AI 未配置，请设置 ai.api_key 或 GOTRIBE_AI_API_KEY", nil)
	}

	switch strings.TrimSpace(req.Task) {
	case "post_metadata":
		return s.generatePostMetadata(ctx, req)
	case "post_slug":
		return s.generatePostSlug(ctx, req)
	case "post_description":
		return s.generatePostDescription(ctx, req)
	case "post_tags":
		return s.generatePostTags(ctx, req)
	default:
		return nil, errs.BadRequest("不支持的 AI 任务", nil)
	}
}

func (s *service) generatePostSlug(ctx context.Context, req *dto.GenerateRequest) (*dto.GenerateResponse, error) {
	title := stringInput(req.Input, "title")
	if strings.TrimSpace(title) == "" {
		return nil, errs.BadRequest("标题不能为空", nil)
	}

	raw, err := s.chat(
		ctx,
		"你是 CMS 的 SEO URL 编辑助手。只输出 JSON，不要输出 Markdown。",
		fmt.Sprintf(`请根据文章标题生成英文 URL slug。

要求：
1. slug 必须是英文小写，只允许 a-z、0-9、连字符。
2. 不能使用拼音，要翻译/概括成自然英文。
3. 长度 20-90 字符。

标题：%s

请严格返回 JSON：
{"slug":"..."}`, title),
	)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Slug string `json:"slug"`
	}
	if err := json.Unmarshal([]byte(extractJSONObject(raw)), &parsed); err != nil {
		return nil, errs.Internal("AI 返回解析失败", err)
	}
	return &dto.GenerateResponse{Result: map[string]any{"slug": normalizeSlug(parsed.Slug)}}, nil
}

func (s *service) generatePostDescription(ctx context.Context, req *dto.GenerateRequest) (*dto.GenerateResponse, error) {
	title := stringInput(req.Input, "title")
	content := stringInput(req.Input, "content")
	if strings.TrimSpace(title) == "" {
		return nil, errs.BadRequest("标题不能为空", nil)
	}
	if len(content) > 6000 {
		content = content[:6000]
	}

	language := strings.TrimSpace(req.Language)
	if language == "" {
		language = "zh-CN"
	}

	raw, err := s.chat(
		ctx,
		"你是 CMS 的文章简介编辑助手。只输出 JSON，不要输出 Markdown。",
		fmt.Sprintf(`请根据文章标题和正文提炼文章简介。

要求：
1. 使用 %s。
2. 80-160 个中文字符或等量长度。
3. 具体、克制、有信息量。
4. 不要编造正文没有的信息。

标题：%s

正文：
%s

请严格返回 JSON：
{"description":"..."}`, language, title, content),
	)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(extractJSONObject(raw)), &parsed); err != nil {
		return nil, errs.Internal("AI 返回解析失败", err)
	}
	return &dto.GenerateResponse{Result: map[string]any{"description": strings.TrimSpace(parsed.Description)}}, nil
}

func (s *service) generatePostMetadata(ctx context.Context, req *dto.GenerateRequest) (*dto.GenerateResponse, error) {
	title := stringInput(req.Input, "title")
	content := stringInput(req.Input, "content")
	if strings.TrimSpace(title) == "" {
		return nil, errs.BadRequest("标题不能为空", nil)
	}
	if len(content) > 6000 {
		content = content[:6000]
	}

	language := strings.TrimSpace(req.Language)
	if language == "" {
		language = "zh-CN"
	}

	systemPrompt := "你是 CMS 的 SEO 编辑助手。只输出 JSON，不要输出 Markdown。"
	userPrompt := fmt.Sprintf(`请根据文章标题和正文生成内容元数据。

要求：
1. slug 必须是英文小写 URL slug，只允许 a-z、0-9、连字符，不能用拼音，长度 20-90 字符。
2. description 使用 %s，提炼文章简介，80-160 个中文字符或等量长度，具体、克制、有信息量。
3. 不要编造正文没有的信息。

标题：%s

正文：
%s

请严格返回 JSON：
{"slug":"...","description":"..."}`, language, title, content)

	raw, err := s.chat(ctx, systemPrompt, userPrompt)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Slug        string `json:"slug"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal([]byte(extractJSONObject(raw)), &parsed); err != nil {
		return nil, errs.Internal("AI 返回解析失败", err)
	}

	result := map[string]any{
		"slug":        normalizeSlug(parsed.Slug),
		"description": strings.TrimSpace(parsed.Description),
	}
	return &dto.GenerateResponse{Result: result}, nil
}

func (s *service) generatePostTags(ctx context.Context, req *dto.GenerateRequest) (*dto.GenerateResponse, error) {
	title := stringInput(req.Input, "title")
	description := stringInput(req.Input, "description")
	content := stringInput(req.Input, "content")
	if strings.TrimSpace(title) == "" && strings.TrimSpace(description) == "" && strings.TrimSpace(content) == "" {
		return nil, errs.BadRequest("标题、简介或正文不能为空", nil)
	}
	if len(content) > 6000 {
		content = content[:6000]
	}

	language := strings.TrimSpace(req.Language)
	if language == "" {
		language = "zh-CN"
	}

	raw, err := s.chat(
		ctx,
		"你是 CMS 的内容分类编辑助手。只输出 JSON，不要输出 Markdown。",
		fmt.Sprintf(`请根据文章内容生成 1-3 个标签。

要求：
1. 标签使用 %s，短而具体，每个 2-10 个字符或等量长度。
2. 不要超过 3 个标签，不要输出重复标签。
3. color 必须是适合标签展示的十六进制颜色，格式为 #RRGGBB。
4. 不要编造正文没有的信息。

标题：%s

简介：
%s

正文：
%s

请严格返回 JSON：
{"tags":[{"title":"...","color":"#3366CC"}]}`, language, title, description, content),
	)
	if err != nil {
		return nil, err
	}

	var parsed struct {
		Tags []struct {
			Title string `json:"title"`
			Color string `json:"color"`
		} `json:"tags"`
	}
	if err := json.Unmarshal([]byte(extractJSONObject(raw)), &parsed); err != nil {
		return nil, errs.Internal("AI 返回解析失败", err)
	}

	tags := make([]map[string]string, 0, 3)
	seen := map[string]struct{}{}
	for _, item := range parsed.Tags {
		title := strings.TrimSpace(item.Title)
		if title == "" {
			continue
		}
		key := strings.ToLower(title)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		tags = append(tags, map[string]string{
			"title": title,
			"color": normalizeHexColor(item.Color),
		})
		if len(tags) == 3 {
			break
		}
	}
	return &dto.GenerateResponse{Result: map[string]any{"tags": tags}}, nil
}

func (s *service) chat(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	baseURL := strings.TrimRight(s.cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = "https://api.deepseek.com"
	}
	model := strings.TrimSpace(s.cfg.Model)
	if model == "" {
		model = "deepseek-v4-flash"
	}

	body := map[string]any{
		"model": model,
		"messages": []map[string]string{
			{"role": "system", "content": systemPrompt},
			{"role": "user", "content": userPrompt},
		},
		"temperature":     0.2,
		"max_tokens":      300,
		"response_format": map[string]string{"type": "json_object"},
		"thinking":        map[string]string{"type": "disabled"},
		"stream":          false,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", errs.Internal("AI 请求构造失败", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", errs.Internal("AI 请求构造失败", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(s.cfg.APIKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(httpReq)
	if err != nil {
		return "", errs.ServiceUnavailable("AI 服务请求失败", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", errs.ServiceUnavailable("AI 服务返回异常", fmt.Errorf("status=%d body=%s", resp.StatusCode, string(respBody)))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", errs.Internal("AI 响应解析失败", err)
	}
	if len(result.Choices) == 0 || strings.TrimSpace(result.Choices[0].Message.Content) == "" {
		return "", errs.ServiceUnavailable("AI 未返回内容", nil)
	}
	return result.Choices[0].Message.Content, nil
}

func stringInput(input map[string]any, key string) string {
	value, ok := input[key]
	if !ok || value == nil {
		return ""
	}
	if s, ok := value.(string); ok {
		return s
	}
	return fmt.Sprint(value)
}

func extractJSONObject(raw string) string {
	raw = strings.TrimSpace(raw)
	start := strings.Index(raw, "{")
	end := strings.LastIndex(raw, "}")
	if start >= 0 && end >= start {
		return raw[start : end+1]
	}
	return raw
}

func normalizeSlug(raw string) string {
	slug := strings.ToLower(strings.TrimSpace(raw))
	var b strings.Builder
	lastDash := false
	for _, r := range slug {
		isAlphaNum := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if isAlphaNum {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func normalizeHexColor(raw string) string {
	color := strings.ToUpper(strings.TrimSpace(raw))
	if len(color) == 6 {
		color = "#" + color
	}
	if len(color) != 7 || color[0] != '#' {
		return "#64748B"
	}
	for _, r := range color[1:] {
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'F')) {
			return "#64748B"
		}
	}
	return color
}
