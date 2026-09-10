package apperror

// Kind 表示业务错误的应用层分类，web 层会据此转换响应状态。
type Kind uint8

const (
	// Unknown 表示未指定分类，通常应按服务器内部错误处理。
	Unknown Kind = iota
	// InvalidInput 表示请求参数或输入值不合法。
	InvalidInput
	// NotFound 表示请求的资源不存在。
	NotFound
	// Conflict 表示请求与资源当前状态冲突，例如重复创建。
	Conflict
	// Unauthenticated 表示请求没有有效身份凭证。
	Unauthenticated
	// Forbidden 表示请求者已认证，但没有执行操作的权限。
	Forbidden
	// BusinessResult 表示业务流程已处理，但需要通过 HTTP 200 返回业务结果。
	BusinessResult
)

// Error 是稳定的、与传输层解耦的业务错误。
type Error struct {
	kind    Kind
	code    int
	message string
}

// New 使用公开的应用元数据创建业务错误。
func New(kind Kind, code int, message string) *Error {
	return &Error{kind: kind, code: code, message: message}
}

// Error 返回客户端可见的错误消息。
func (e *Error) Error() string {
	return e.message
}

// Kind 返回业务错误的应用层分类。
func (e *Error) Kind() Kind {
	return e.kind
}

// Code 返回稳定的应用错误码。
func (e *Error) Code() int {
	return e.code
}

// Message 返回面向客户端的应用消息。
func (e *Error) Message() string {
	return e.message
}
