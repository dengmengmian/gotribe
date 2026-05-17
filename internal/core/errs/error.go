package errs

// 本文件定义统一错误结构和常用错误构造函数。

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// AppError 表示统一业务错误结构，便于跨层传递和输出。
type AppError struct {
	Code        Code           `json:"code"`
	Message     string         `json:"message"`
	MessageKey  MsgKey         `json:"-"`
	MessageArgs []any          `json:"-"`
	Status      int            `json:"-"`
	Details     map[string]any `json:"details,omitempty"`
	Err         error          `json:"-"`
}

// Error 返回业务错误的字符串表示。
func (e *AppError) Error() string {
	msg := e.Message
	if msg == "" && e.MessageKey != "" {
		msg = T("zh", e.MessageKey)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", msg, e.Err)
	}
	return msg
}

// Unwrap 返回被包装的底层错误。
func (e *AppError) Unwrap() error {
	return e.Err
}

// Localize 根据语言设置翻译后的 Message。
func (e *AppError) Localize(lang string) {
	if e.MessageKey == "" {
		return
	}
	msg := T(lang, e.MessageKey)
	if len(e.MessageArgs) > 0 {
		msg = fmt.Sprintf(msg, e.MessageArgs...)
	}
	e.Message = msg
}

// New 创建统一业务错误对象。
func New(code Code, status int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Err:     err,
	}
}

// NewWithKey 创建携带消息键的错误（由 response 层按当前语言翻译）。
func NewWithKey(code Code, status int, key MsgKey, args []any, err error) *AppError {
	return &AppError{
		Code:        code,
		MessageKey:  key,
		MessageArgs: args,
		Status:      status,
		Err:         err,
	}
}

// BadRequestWithKey 创建参数错误（消息键版）。
func BadRequestWithKey(key MsgKey, args []any, err error) error {
	return NewWithKey(CodeBadRequest, http.StatusBadRequest, key, args, err)
}

// NotFoundWithKey 创建资源不存在错误（消息键版）。
func NotFoundWithKey(key MsgKey, args []any, err error) error {
	return NewWithKey(CodeNotFound, http.StatusNotFound, key, args, err)
}

// InternalWithKey 创建内部错误（消息键版）。
func InternalWithKey(key MsgKey, args []any, err error) error {
	return NewWithKey(CodeInternal, http.StatusInternalServerError, key, args, err)
}

// ConflictWithKey 创建资源冲突错误（消息键版）。
func ConflictWithKey(key MsgKey, args []any, err error) error {
	return NewWithKey(CodeConflict, http.StatusConflict, key, args, err)
}

// BadRequest 创建参数错误类型的业务错误。
func BadRequest(message string, err error) error {
	return New(CodeBadRequest, http.StatusBadRequest, message, err)
}

// BadRequestWithDetails 创建附带字段详情的参数错误。
func BadRequestWithDetails(message string, details map[string]any, err error) error {
	appErr := New(CodeBadRequest, http.StatusBadRequest, message, err)
	appErr.Details = details
	return appErr
}

// Unauthorized 创建未认证错误。
func Unauthorized(message string) error {
	return New(CodeUnauthorized, http.StatusUnauthorized, message, nil)
}

// Forbidden 创建无权限错误。
func Forbidden(message string) error {
	return New(CodeForbidden, http.StatusForbidden, message, nil)
}

// NotFound 创建资源不存在错误。
func NotFound(message string, err error) error {
	return New(CodeNotFound, http.StatusNotFound, message, err)
}

// Conflict 创建资源冲突错误。
func Conflict(message string, err error) error {
	return New(CodeConflict, http.StatusConflict, message, err)
}

// TooManyRequests 创建请求过于频繁错误。
func TooManyRequests(message string) error {
	return New(CodeRateLimited, http.StatusTooManyRequests, message, nil)
}

// ServiceUnavailable 创建服务不可用错误。
func ServiceUnavailable(message string, err error) error {
	return New(CodeUnavailable, http.StatusServiceUnavailable, message, err)
}

// Internal 创建内部错误。
func Internal(message string, err error) error {
	return New(CodeInternal, http.StatusInternalServerError, message, err)
}

// AccountLocked 创建账户已锁定错误，Details 携带 locked_until / remaining_seconds 字段。
func AccountLocked(message string, lockedUntilUnix int64, remainingSeconds int64) error {
	appErr := New(CodeAccountLocked, http.StatusTooManyRequests, message, nil)
	appErr.Details = map[string]any{
		"locked_until":      lockedUntilUnix,
		"remaining_seconds": remainingSeconds,
	}
	return appErr
}

// TOTPRequired 创建需要 TOTP 二次校验的"业务态错误"，HTTP 200 由 handler 包装为正常响应。
// 此函数仅用于内部传递阶段信息，handler 层不会直接把它转成错误响应。
func TOTPRequired(stepToken string, stepExpiresUnix int64) error {
	appErr := New(CodeTOTPRequired, http.StatusOK, "需要 TOTP 二次校验", nil)
	appErr.Details = map[string]any{
		"step_token":   stepToken,
		"step_expires": stepExpiresUnix,
	}
	return appErr
}

// TOTPInvalid 创建 TOTP 校验失败错误。
func TOTPInvalid(message string) error {
	return New(CodeTOTPInvalid, http.StatusUnauthorized, message, nil)
}

// TOTPAlreadyBound 创建已绑定 TOTP 的冲突错误。
func TOTPAlreadyBound(message string) error {
	return New(CodeTOTPAlreadyBound, http.StatusConflict, message, nil)
}

// TOTPNotBound 创建未绑定 TOTP 的错误。
func TOTPNotBound(message string) error {
	return New(CodeTOTPNotBound, http.StatusBadRequest, message, nil)
}

// IsUniqueViolation 判断错误是否为数据库唯一约束冲突。
func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

// ConstraintName 提取数据库约束冲突对应的约束名。
func ConstraintName(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.ConstraintName
	}
	return ""
}

// As 将任意错误转换为统一业务错误结构。
func As(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return New(CodeNotFound, http.StatusNotFound, "resource not found", err)
	}
	return New(CodeInternal, http.StatusInternalServerError, "internal server error", err)
}
