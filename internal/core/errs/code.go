// Package errs defines unified business error codes, error structures, and helper constructors.
package errs

// Code 表示统一业务错误码。
type Code string

const (
	// CodeBadRequest indicates the request parameters are invalid.
	CodeBadRequest Code = "bad_request"
	// CodeUnauthorized indicates the request lacks valid authentication credentials.
	CodeUnauthorized Code = "unauthorized"
	// CodeForbidden indicates the authenticated user does not have permission.
	CodeForbidden Code = "forbidden"
	// CodeNotFound indicates the requested resource could not be found.
	CodeNotFound Code = "not_found"
	// CodeConflict indicates a resource conflict occurred, such as a unique constraint violation.
	CodeConflict Code = "conflict"
	// CodeRateLimited indicates the client has exceeded the rate limit.
	CodeRateLimited Code = "rate_limit_exceeded"
	// CodeUnavailable indicates the service is temporarily unavailable.
	CodeUnavailable Code = "service_unavailable"
	// CodeInternal indicates an unexpected internal server error occurred.
	CodeInternal Code = "internal_error"
	// CodeAccountLocked indicates the account is temporarily locked due to failed login attempts.
	CodeAccountLocked Code = "account_locked"
	// CodeTOTPRequired indicates the caller must complete TOTP second-factor verification.
	CodeTOTPRequired Code = "totp_required"
	// CodeTOTPInvalid indicates the submitted TOTP code or recovery code is wrong.
	CodeTOTPInvalid Code = "totp_invalid"
	// CodeTOTPAlreadyBound indicates the caller already bound a TOTP secret.
	CodeTOTPAlreadyBound Code = "totp_already_bound"
	// CodeTOTPNotBound indicates the caller has no active TOTP secret.
	CodeTOTPNotBound Code = "totp_not_bound"
)
