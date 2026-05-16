// Package constant 定义全局应用常量。
package constant

// 上下文键
const (
	// ContextKeyUserID 是 context 中用户 ID 的键。
	ContextKeyUserID = "user_id"
	// ContextKeyUsername 是 context 中用户名的键。
	ContextKeyUsername = "username"
	// ContextKeyEmail 是 context 中邮箱的键。
	ContextKeyEmail = "email"
	// ContextKeyRequestID 是 context 中请求 ID 的键。
	ContextKeyRequestID = "request_id"
	// ContextKeyTraceID 是 context 中追踪 ID 的键。
	ContextKeyTraceID = "trace_id"
	// ContextKeyRole 是 context 中用户角色的键。
	ContextKeyRole = "role"
)

// 状态值
const (
	// StatusActive 表示激活状态。
	StatusActive = 1
	// StatusInactive 表示未激活状态。
	StatusInactive = 0
	// StatusDeleted 表示已删除状态。
	StatusDeleted = -1
	// StatusPending 表示待处理状态。
	StatusPending = 2
	// StatusDisabled 表示已禁用状态。
	StatusDisabled = 3
)

// 日期格式化
const (
	TIME_FORMAT_DAY   = "20060102"
	TIME_FORMAT       = "2006-01-02 15:04:05"
	TIME_FORMAT_SHORT = "2006-01-02"
)

// 默认值
const (
	DEFAULT_ID           = 1
	DEFAULT_UPLOAD_SIZE  = 200 * 1024 * 1024
)

// 文件类型
const (
	FILE_TYPE_IMAGE    = 1
	FILE_TYPE_VIDEO    = 2
	FILE_TYPE_AUDIO    = 3
	FILE_TYPE_ARCHIVE  = 4
	FILE_TYPE_DOCUMENT = 5
	FILE_TYPE_FONT     = 6
	FILE_TYPE_APP      = 7
	FILE_TYPE_UNKNOWN  = 8
)

// 审核状态
const (
	AUDIT_STATUS_PENDING = 1
	AUDIT_STATUS_PASS    = 2
)

// 文章状态
const (
	POST_STATUS_DRAFT  = 1
	POST_STATUS_PUBLIC = 2
)

// 文章类型
const (
	POST_TYPE_POST = 1
	POST_TYPE_PAGE = 2
)

// 订单状态
const (
	OrderStatusPendingPayment = 1
	OrderStatusPaid           = 2
	OrderStatusShipped        = 3
	OrderStatusReceived       = 4
	OrderStatusCanceled       = 5
	OrderStatusRefunding      = 6
	OrderStatusRefunded       = 7
	OrderStatusRefundFailed   = 8
)

// 订单类型
const (
	OrderTypeNormal = 1
	OrderTypePoint  = 2
)

// 支付方式
const (
	PaymentMethodWechatPay = 1
	PaymentMethodAlipay    = 2
	PaymentMethodPoint     = 3
	PaymentMethodBalance   = 4
)
