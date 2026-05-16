// Package request provides unified request binding, validation, and multilingual error translation.
package request

// 本文件封装统一的请求绑定、参数校验和中英文错误翻译逻辑。

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	enlocale "github.com/go-playground/locales/en"
	zhlocale "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"regexp"

	"gotribe/internal/core/errs"
)

const (
	localeEN = "en"
	localeZH = "zh"
)

var (
	initOnce sync.Once
	initErr  error
	uni      *ut.UniversalTranslator
)

// Initialize 初始化请求绑定和多语言校验器。
func Initialize() error {
	initOnce.Do(func() {
		engine, ok := binding.Validator.Engine().(*validator.Validate)
		if !ok {
			initErr = errors.New("gin validator engine is not available")
			return
		}

		engine.RegisterTagNameFunc(func(field reflect.StructField) string {
			for _, key := range []string{"json", "form"} {
				if name := tagName(field.Tag.Get(key)); name != "" {
					return name
				}
			}
			return field.Name
		})

		en := enlocale.New()
		zh := zhlocale.New()
		uni = ut.New(en, en, zh)

		enTranslator, found := uni.GetTranslator(localeEN)
		if !found {
			initErr = errors.New("english translator is not available")
			return
		}
		if err := entranslations.RegisterDefaultTranslations(engine, enTranslator); err != nil {
			initErr = fmt.Errorf("register english validation translations: %w", err)
			return
		}

		zhTranslator, found := uni.GetTranslator(localeZH)
		if !found {
			initErr = errors.New("chinese translator is not available")
			return
		}
		if err := zhtranslations.RegisterDefaultTranslations(engine, zhTranslator); err != nil {
			initErr = fmt.Errorf("register chinese validation translations: %w", err)
			return
		}

		if err := engine.RegisterValidation("checkMobile", checkMobile); err != nil {
			initErr = fmt.Errorf("register checkMobile validation: %w", err)
			return
		}
	})
	return initErr
}

// MaxBodyBytes 限制 JSON 请求体的最大字节数，默认 10 MB。
// 在应用启动时通过 bootstrap 设置为配置值。
var MaxBodyBytes int64 = 10 << 20 // 10 MB

// BindJSON 绑定 JSON 请求体并统一处理校验错误。
func BindJSON(c *gin.Context, target any) error {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	if err := c.ShouldBindJSON(target); err != nil {
		return wrapBindError(c, err, messageInvalidBody)
	}
	return nil
}

// BindQuery 绑定查询参数并统一处理校验错误。
func BindQuery(c *gin.Context, target any) error {
	if err := c.ShouldBindQuery(target); err != nil {
		return wrapBindError(c, err, messageInvalidQuery)
	}
	return nil
}

// wrapBindError 将绑定或校验失败统一转换为业务错误。
func wrapBindError(c *gin.Context, err error, key messageKey) error {
	locale := localeForRequest(c)
	details := translateValidationErrors(err, locale)
	return errs.BadRequestWithDetails(localize(locale, key), details, err)
}

// translateValidationErrors 将校验错误翻译为可直接返回给客户端的字段提示。
func translateValidationErrors(err error, locale string) map[string]any {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) || uni == nil {
		return nil
	}

	translator, found := uni.GetTranslator(locale)
	if !found {
		translator, _ = uni.GetTranslator(localeEN)
	}

	fields := make(map[string]string, len(validationErrors))
	for _, item := range validationErrors {
		fields[item.Field()] = item.Translate(translator)
	}
	if len(fields) == 0 {
		return nil
	}
	return map[string]any{
		"fields": fields,
		"locale": locale,
	}
}

// localeForRequest 从请求头中解析当前请求使用的语言。
func localeForRequest(c *gin.Context) string {
	for _, raw := range []string{c.GetHeader("X-Language"), c.GetHeader("Accept-Language")} {
		if normalized := normalizeLocale(raw); normalized != "" {
			return normalized
		}
	}
	return localeEN
}

// normalizeLocale 将语言标识规整为系统支持的语言代码。
func normalizeLocale(raw string) string {
	raw = strings.TrimSpace(strings.ToLower(raw))
	if raw == "" {
		return ""
	}

	for _, part := range strings.Split(raw, ",") {
		value := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		switch {
		case strings.HasPrefix(value, "zh"):
			return localeZH
		case strings.HasPrefix(value, "en"):
			return localeEN
		}
	}
	return ""
}

// tagName 提取用于返回给客户端的字段名。
func tagName(raw string) string {
	if raw == "" || raw == "-" {
		return ""
	}
	name := strings.TrimSpace(strings.SplitN(raw, ",", 2)[0])
	if name == "" || name == "-" {
		return ""
	}
	return name
}

// messageKey 表示请求绑定错误的消息键。
type messageKey string

const (
	messageInvalidBody  messageKey = "invalid_body"
	messageInvalidQuery messageKey = "invalid_query"
)

// localize 根据语言和消息键返回对应提示文案。
func localize(locale string, key messageKey) string {
	switch key {
	case messageInvalidBody:
		if locale == localeZH {
			return "请求体参数错误"
		}
		return "request body is invalid"
	case messageInvalidQuery:
		if locale == localeZH {
			return "查询参数错误"
		}
		return "query parameters are invalid"
	default:
		if locale == localeZH {
			return "请求参数错误"
		}
		return "request parameters are invalid"
	}
}

func checkMobile(fl validator.FieldLevel) bool {
	reg := `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
	rgx := regexp.MustCompile(reg)
	return rgx.MatchString(fl.Field().String())
}
