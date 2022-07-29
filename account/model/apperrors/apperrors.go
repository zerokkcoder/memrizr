package apperrors

import (
	"errors"
	"fmt"
	"net/http"
)

// Type 保存错误的类型字符串和整数代码
type Type string

// 定义错误类型
const (
	Authorization        Type = "AUTHORIZATION"        // Authentication Failures -
	BadRequest           Type = "BADREQUEST"           // Validation errors / BadInput
	Conflict             Type = "CONFLICT"             // Already exists (eg, create account with existent email) - 409
	Internal             Type = "INTERNAL"             // Server (500) and fallback errors
	NotFound             Type = "NOTFOUND"             // For not finding resource
	PayloadTooLarge      Type = "PAYLOADTOOLARGE"      // for uploading tons of JSON, or an image over the limit - 413
	UnsupportedMediaType Type = "UNSUPPORTEDMEDIATYPE" // for http 415
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"  // For long running handlers
)

// Error 应用程序的自定义错误
// 这有助于从API端点返回一致的错误类型/消息
type Error struct {
	Type    Type   `json:"type"`
	Message string `json:message`
}

// Error 实现标准的错误接口，我们可以从这个包返回一个普通的旧的 go _error_ 错误
func (e Error) Error() string {
	return e.Message
}

// Status 错误到状态码的映射
// 当然，这有些多余，因为我们的错误已经映射了http状态代码
func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// Status 检查错误的运行时类型
// 如果错误是 model.Error ，则返回对应错误 http状态码
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

// Error 工厂函数

// NewAuthorization 创建 401 error
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

// NewBadRequest 创建 400 error
func NewBadRequest(reaseon string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Bad request. Reason: %v", reaseon),
	}
}

// NewConflict 创建 409 error
func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists.", name, value),
	}
}

// NewInternal 创建 500 error 或者 unknow error
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: fmt.Sprintf("Internal server error."),
	}
}

// NewNotFound 创建 404 error
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("resource: %v with value: %v not found", name, value),
	}
}

// NewPayloadTooLarge 创建 403 error
func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
	}
}

// NewUnsupportedMediaType 创建 415 error
func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
	}
}

// NewServiceUnavailable 创建 503 error
func NewServiceUnavailable() *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: fmt.Sprintf("Service unavailable or timed out"),
	}
}
