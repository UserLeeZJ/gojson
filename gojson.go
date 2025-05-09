// Package gojson 提供类似JavaScript JSON接口的Go JSON工具库。
// 该库专注于提供面向对象的JSON操作体验，同时保持高性能和可靠性。
package gojson

import (
	"github.com/UserLeeZJ/gojson/diff"
	"github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/fast"
	"github.com/UserLeeZJ/gojson/jsonpath"
	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/patch"
	"github.com/UserLeeZJ/gojson/stream"
	"github.com/UserLeeZJ/gojson/types"
	"github.com/UserLeeZJ/gojson/utils"
)

// 版本信息常量。
const (
	// Version 是gojson库的版本号。
	Version = "1.0.2"

	// APIVersion 是gojson API的版本号。
	APIVersion = "1.0.2"
)

// JSON Path相关常量。
const (
	// JSONPathRoot 是JSON Path的根节点表示。
	JSONPathRoot = "$"
)

// 重新导出的类型。
type (
	JSONValue   = types.JSONValue
	JSONObject  = types.JSONObject
	JSONArray   = types.JSONArray
	JSONString  = types.JSONString
	JSONNumber  = types.JSONNumber
	JSONBool    = types.JSONBool
	JSONNull    = types.JSONNull
	JSONError   = errors.JSONError
	ErrorCode   = errors.ErrorCode
	DiffType    = diff.DiffType
	Diff        = diff.Diff
	DiffOptions = diff.DiffOptions

	// 流式处理相关类型
	JSONTokenType     = stream.JSONTokenType
	JSONToken         = stream.JSONToken
	JSONTokenizer     = stream.JSONTokenizer
	JSONGenerator     = stream.JSONGenerator
	IncrementalParser = stream.IncrementalParser
)

// 重新导出的错误代码常量。
const (
	ErrInvalidJSON     = errors.ErrInvalidJSON
	ErrEmptyInput      = errors.ErrEmptyInput
	ErrInvalidType     = errors.ErrInvalidType
	ErrTypeConversion  = errors.ErrTypeConversion
	ErrPathNotFound    = errors.ErrPathNotFound
	ErrInvalidPath     = errors.ErrInvalidPath
	ErrIndexOutOfRange = errors.ErrIndexOutOfRange
	ErrInvalidIndex    = errors.ErrInvalidIndex
	ErrOperationFailed = errors.ErrOperationFailed
	ErrNotSupported    = errors.ErrNotSupported
	ErrInvalidPatch    = errors.ErrInvalidPatch
	ErrPatchFailed     = errors.ErrPatchFailed
	ErrTestFailed      = errors.ErrTestFailed
)

// 重新导出的流式处理常量。
const (
	TokenError        = stream.TokenError
	TokenObjectStart  = stream.TokenObjectStart
	TokenObjectEnd    = stream.TokenObjectEnd
	TokenArrayStart   = stream.TokenArrayStart
	TokenArrayEnd     = stream.TokenArrayEnd
	TokenPropertyName = stream.TokenPropertyName
	TokenString       = stream.TokenString
	TokenNumber       = stream.TokenNumber
	TokenBoolean      = stream.TokenBoolean
	TokenNull         = stream.TokenNull
	TokenEOF          = stream.TokenEOF
)

// 重新导出的构造函数。
var (
	NewJSONObject          = types.NewJSONObject
	NewJSONArray           = types.NewJSONArray
	NewJSONArrayFromValues = types.NewJSONArrayFromValues
	NewJSONString          = types.NewJSONString
	NewJSONNumber          = types.NewJSONNumber
	NewJSONBool            = types.NewJSONBool
	NewJSONNull            = types.NewJSONNull
	NewJSONError           = errors.NewJSONError
)

// 重新导出的解析函数。
var (
	ParseToValue      = parser.ParseToValue
	ParseBytesToValue = parser.ParseBytesToValue
	Parse             = parser.Parse
	ParseBytes        = parser.ParseBytes
	Stringify         = parser.Stringify
	StringifyBytes    = parser.StringifyBytes
	StringifyIndent   = parser.StringifyIndent
)

// 重新导出的JSON Path函数。
var (
	ParseJSONPath       = jsonpath.ParseJSONPath
	QueryJSONPath       = jsonpath.QueryJSONPath
	QueryJSONPathString = jsonpath.QueryJSONPathString
)

// 重新导出的JSON Diff函数。
var (
	DiffJSON           = diff.DiffJSON
	DiffJSONStrings    = diff.DiffJSONStrings
	DefaultDiffOptions = diff.DefaultDiffOptions
)

// 重新导出的JSON Patch函数。
var (
	ApplyPatch    = patch.ApplyPatch
	GeneratePatch = diff.GeneratePatch
)

// 重新导出的类型转换函数。
var (
	ValueToInterface = types.ValueToInterface
)

// 重新导出的性能优化函数。
var (
	// FastMarshal 是一个优化的JSON序列化函数。
	FastMarshal = fast.Marshal
	// FastUnmarshal 是一个优化的JSON反序列化函数。
	FastUnmarshal = fast.Unmarshal
	// CacheFragment 缓存JSON片段。
	CacheFragment = fast.CacheFragment
	// GetCachedFragment 获取缓存的JSON片段。
	GetCachedFragment = fast.GetCachedFragment
	// ClearFragmentCache 清空片段缓存。
	ClearFragmentCache = fast.ClearFragmentCache
)

// 重新导出的流式处理函数。
var (
	// NewJSONTokenizer 创建一个新的JSON流式解析器。
	NewJSONTokenizer = stream.NewJSONTokenizer
	// NewJSONGenerator 创建一个新的JSON流式生成器。
	NewJSONGenerator = stream.NewJSONGenerator
	// NewIncrementalParser 创建一个新的增量JSON解析器。
	NewIncrementalParser = stream.NewIncrementalParser
)

// 注意：泛型函数不能直接导出，需要在使用时导入generic包并指定类型参数
// 例如：generic.NewJSONObject[map[string]interface{}]()
// 例如：generic.GetTyped[string](obj, "key")

// 重新导出的工具函数。
var (
	// FormatJSON 格式化JSON字符串。
	FormatJSON = utils.FormatJSON
	// CompactJSON 压缩JSON字符串。
	CompactJSON = utils.CompactJSON
	// SortJSONKeys 对JSON对象的键进行排序。
	SortJSONKeys = utils.SortJSONKeys
	// ValidateJSON 验证JSON字符串是否有效。
	ValidateJSON = utils.ValidateJSON
	// MergeJSON 合并两个JSON对象。
	MergeJSON = utils.MergeJSON
	// DeepCopy 深度复制JSON值。
	DeepCopy = utils.DeepCopy
)
