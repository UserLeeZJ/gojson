// Package errors 提供gojson库的错误处理功能。
package errors

import (
	"fmt"
)

// ErrorCode 表示错误类型的枚举。
type ErrorCode string

// 定义错误代码常量。
const (
	// 通用错误。
	ErrInvalidJSON ErrorCode = "INVALID_JSON"
	ErrEmptyInput  ErrorCode = "EMPTY_INPUT"

	// 类型错误。
	ErrInvalidType    ErrorCode = "INVALID_TYPE"
	ErrTypeConversion ErrorCode = "TYPE_CONVERSION"

	// 路径错误。
	ErrPathNotFound ErrorCode = "PATH_NOT_FOUND"
	ErrInvalidPath  ErrorCode = "INVALID_PATH"

	// 索引错误。
	ErrIndexOutOfRange ErrorCode = "INDEX_OUT_OF_RANGE"
	ErrInvalidIndex    ErrorCode = "INVALID_INDEX"

	// 操作错误。
	ErrOperationFailed ErrorCode = "OPERATION_FAILED"
	ErrNotSupported    ErrorCode = "NOT_SUPPORTED"

	// Patch 错误。
	ErrInvalidPatch ErrorCode = "INVALID_PATCH"
	ErrPatchFailed  ErrorCode = "PATCH_FAILED"
	ErrTestFailed   ErrorCode = "TEST_FAILED"
)

// JSONError 表示JSON操作中的错误。
type JSONError struct {
	Code    ErrorCode // 错误代码。
	Message string    // 错误消息。
	Path    string    // 错误发生的路径（如果适用）。
	Cause   error     // 原始错误（如果有）。
}

// Error 实现error接口。
func (e *JSONError) Error() string {
	if e.Path != "" {
		if e.Cause != nil {
			return fmt.Sprintf("%s: %s (path: %s): %v", e.Code, e.Message, e.Path, e.Cause)
		}
		return fmt.Sprintf("%s: %s (path: %s)", e.Code, e.Message, e.Path)
	}

	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap 返回底层错误，支持Go 1.13+的错误链。
func (e *JSONError) Unwrap() error {
	return e.Cause
}

// NewJSONError 创建一个新的JSONError。
func NewJSONError(code ErrorCode, message string) *JSONError {
	return &JSONError{
		Code:    code,
		Message: message,
	}
}

// WithPath 添加路径信息到错误。
func (e *JSONError) WithPath(path string) *JSONError {
	e.Path = path
	return e
}

// WithCause 添加原始错误。
func (e *JSONError) WithCause(err error) *JSONError {
	e.Cause = err
	return e
}

// 以下是常用错误创建函数。

// ErrInvalidTypeWithDetails 创建类型错误详情。
func ErrInvalidTypeWithDetails(expected, actual string) *JSONError {
	return NewJSONError(ErrInvalidType,
		fmt.Sprintf("类型错误: 期望 %s, 实际 %s", expected, actual))
}

// ErrPathNotFoundWithDetails 创建路径不存在错误详情。
func ErrPathNotFoundWithDetails(path string) *JSONError {
	return NewJSONError(ErrPathNotFound,
		fmt.Sprintf("路径不存在")).WithPath(path)
}

// ErrIndexOutOfRangeWithDetails 创建索引超出范围错误详情。
func ErrIndexOutOfRangeWithDetails(index, size int) *JSONError {
	return NewJSONError(ErrIndexOutOfRange,
		fmt.Sprintf("索引超出范围: %d (大小: %d)", index, size))
}

// ErrInvalidPathWithDetails 创建无效路径错误详情。
func ErrInvalidPathWithDetails(path, reason string) *JSONError {
	return NewJSONError(ErrInvalidPath,
		fmt.Sprintf("无效的路径: %s", reason)).WithPath(path)
}

// ErrInvalidJSONWithDetails 创建无效JSON错误详情。
func ErrInvalidJSONWithDetails(reason string) *JSONError {
	return NewJSONError(ErrInvalidJSON,
		fmt.Sprintf("无效的JSON: %s", reason))
}

// ErrPatchFailedWithDetails 创建补丁操作失败错误详情。
func ErrPatchFailedWithDetails(op, path, reason string) *JSONError {
	return NewJSONError(ErrPatchFailed,
		fmt.Sprintf("补丁操作 '%s' 失败: %s", op, reason)).WithPath(path)
}

// ErrTestFailedWithDetails 创建测试失败错误详情。
func ErrTestFailedWithDetails(path string, expected, actual interface{}) *JSONError {
	return NewJSONError(ErrTestFailed,
		fmt.Sprintf("测试失败: 期望 %v, 实际 %v", expected, actual)).WithPath(path)
}
