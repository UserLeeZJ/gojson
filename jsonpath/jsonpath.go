// Package jsonpath 提供gojson库的JSON Path查询功能
package jsonpath

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/parser"
	"github.com/UserLeeZJ/gojson/types"
)

// JSONPath 表示一个JSON Path表达式
type JSONPath struct {
	segments []pathSegment
	original string
}

// pathSegment 表示JSON Path的一个段
type pathSegment interface {
	// 应用段到JSON值，返回匹配的值
	apply(value types.JSONValue) ([]types.JSONValue, error)
	// 返回段的字符串表示
	String() string
}

// rootSegment 表示根节点 $
type rootSegment struct{}

func (s *rootSegment) apply(value types.JSONValue) ([]types.JSONValue, error) {
	return []types.JSONValue{value}, nil
}

func (s *rootSegment) String() string {
	return "$"
}

// propertySegment 表示属性访问 .property 或 ['property']
type propertySegment struct {
	name string
}

func (s *propertySegment) apply(value types.JSONValue) ([]types.JSONValue, error) {
	if !value.IsObject() {
		return nil, jsonerrors.ErrInvalidTypeWithDetails("object", value.Type())
	}

	obj, _ := value.AsObject()
	prop := obj.Get(s.name)
	if prop.IsNull() && !obj.Has(s.name) {
		return []types.JSONValue{}, nil
	}
	return []types.JSONValue{prop}, nil
}

func (s *propertySegment) String() string {
	if isValidIdentifier(s.name) {
		return "." + s.name
	}
	return "['" + s.name + "']"
}

// indexSegment 表示数组索引访问 [0]
type indexSegment struct {
	index int
}

func (s *indexSegment) apply(value types.JSONValue) ([]types.JSONValue, error) {
	if !value.IsArray() {
		return nil, jsonerrors.ErrInvalidTypeWithDetails("array", value.Type())
	}

	arr, _ := value.AsArray()
	if s.index < 0 || s.index >= arr.Size() {
		return []types.JSONValue{}, nil
	}
	return []types.JSONValue{arr.Get(s.index)}, nil
}

func (s *indexSegment) String() string {
	return fmt.Sprintf("[%d]", s.index)
}

// wildcardSegment 表示通配符 .* 或 [*]
type wildcardSegment struct{}

func (s *wildcardSegment) apply(value types.JSONValue) ([]types.JSONValue, error) {
	if value.IsObject() {
		obj, _ := value.AsObject()
		result := make([]types.JSONValue, 0, obj.Size())
		for _, key := range obj.Keys() {
			result = append(result, obj.Get(key))
		}
		return result, nil
	} else if value.IsArray() {
		arr, _ := value.AsArray()
		result := make([]types.JSONValue, 0, arr.Size())
		for i := 0; i < arr.Size(); i++ {
			result = append(result, arr.Get(i))
		}
		return result, nil
	}
	return nil, jsonerrors.ErrInvalidTypeWithDetails("object or array", value.Type())
}

func (s *wildcardSegment) String() string {
	return "[*]"
}

// sliceSegment 表示数组切片 [start:end]
type sliceSegment struct {
	start    int
	end      int
	hasStart bool
	hasEnd   bool
}

func (s *sliceSegment) apply(value types.JSONValue) ([]types.JSONValue, error) {
	if !value.IsArray() {
		return nil, jsonerrors.ErrInvalidTypeWithDetails("array", value.Type())
	}

	arr, _ := value.AsArray()
	size := arr.Size()

	// 计算实际的起始和结束索引
	start, end := s.start, s.end

	// 如果没有指定起始索引，默认为0
	if !s.hasStart {
		start = 0
	}

	// 如果没有指定结束索引，默认为数组长度
	if !s.hasEnd {
		end = size
	}

	// 处理负索引
	if start < 0 {
		start = size + start
	}
	if end < 0 {
		end = size + end
	}

	// 确保索引在有效范围内
	if start < 0 {
		start = 0
	}
	if end > size {
		end = size
	}

	// 如果起始索引大于等于结束索引或超出数组范围，返回空数组
	if start >= size || start >= end {
		return []types.JSONValue{}, nil
	}

	// 创建结果数组
	result := make([]types.JSONValue, 0, end-start)
	for i := start; i < end; i++ {
		result = append(result, arr.Get(i))
	}

	return result, nil
}

func (s *sliceSegment) String() string {
	startStr := ""
	if s.hasStart {
		startStr = fmt.Sprintf("%d", s.start)
	}

	endStr := ""
	if s.hasEnd {
		endStr = fmt.Sprintf("%d", s.end)
	}

	return fmt.Sprintf("[%s:%s]", startStr, endStr)
}

// ParseJSONPath 解析JSON Path表达式
func ParseJSONPath(path string) (*JSONPath, error) {
	if path == "" {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPath, "JSON Path不能为空")
	}

	// 创建JSON Path对象
	jp := &JSONPath{
		segments: make([]pathSegment, 0),
		original: path,
	}

	// 解析根节点
	if !strings.HasPrefix(path, "$") {
		return nil, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPath, "JSON Path必须以$开头")
	}

	jp.segments = append(jp.segments, &rootSegment{})
	path = path[1:]

	// 解析剩余部分
	for len(path) > 0 {
		var segment pathSegment
		var consumed int
		var err error

		// 解析下一个段
		if segment, consumed, err = parseNextSegment(path); err != nil {
			return nil, err
		}

		jp.segments = append(jp.segments, segment)
		path = path[consumed:]
	}

	return jp, nil
}

// 解析下一个路径段
func parseNextSegment(path string) (pathSegment, int, error) {
	// 属性访问 .property
	if strings.HasPrefix(path, ".") {
		if len(path) == 1 {
			return nil, 0, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPath, "属性名不能为空")
		}

		// 通配符 .*
		if path[1] == '*' {
			return &wildcardSegment{}, 2, nil
		}

		// 提取属性名
		match := regexp.MustCompile(`^\.([a-zA-Z_][a-zA-Z0-9_]*)`).FindStringSubmatch(path)
		if match == nil {
			return nil, 0, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPath, "无效的属性名")
		}

		return &propertySegment{name: match[1]}, len(match[0]), nil
	}

	// 括号表达式 [...]
	if strings.HasPrefix(path, "[") {
		// 查找匹配的右括号
		depth := 1
		end := 1
		for end < len(path) && depth > 0 {
			if path[end] == '[' {
				depth++
			} else if path[end] == ']' {
				depth--
			}
			end++
		}

		if depth > 0 {
			return nil, 0, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPath, "括号不匹配")
		}

		bracketContent := path[1 : end-1]

		// 通配符 [*]
		if bracketContent == "*" {
			return &wildcardSegment{}, end, nil
		}

		// 数字索引 [0]
		if index, err := strconv.Atoi(bracketContent); err == nil {
			return &indexSegment{index: index}, end, nil
		}

		// 切片 [start:end]
		if strings.Contains(bracketContent, ":") {
			parts := strings.Split(bracketContent, ":")
			if len(parts) == 2 {
				var startIdx, endIdx int
				hasStart, hasEnd := false, false

				// 解析起始索引
				if parts[0] != "" {
					hasStart = true
					startIdx, _ = strconv.Atoi(parts[0])
				}

				// 解析结束索引
				if parts[1] != "" {
					hasEnd = true
					endIdx, _ = strconv.Atoi(parts[1])
				}

				return &sliceSegment{
					start:    startIdx,
					end:      endIdx,
					hasStart: hasStart,
					hasEnd:   hasEnd,
				}, end, nil // 这里的end是指右括号的位置
			}
		}

		// 字符串属性 ['property'] 或 ["property"]
		if (strings.HasPrefix(bracketContent, "'") && strings.HasSuffix(bracketContent, "'")) ||
			(strings.HasPrefix(bracketContent, "\"") && strings.HasSuffix(bracketContent, "\"")) {
			propName := bracketContent[1 : len(bracketContent)-1]
			return &propertySegment{name: propName}, end, nil
		}

		return nil, 0, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPath, "无效的括号表达式")
	}

	return nil, 0, jsonerrors.NewJSONError(jsonerrors.ErrInvalidPath, "无效的路径段")
}

// 检查字符串是否为有效的标识符
func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	if !isLetter(s[0]) && s[0] != '_' {
		return false
	}
	for i := 1; i < len(s); i++ {
		if !isLetter(s[i]) && !isDigit(s[i]) && s[i] != '_' {
			return false
		}
	}
	return true
}

// 检查字符是否为字母
func isLetter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

// 检查字符是否为数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// Query 使用JSON Path查询JSON值
func (jp *JSONPath) Query(value types.JSONValue) ([]types.JSONValue, error) {
	current := []types.JSONValue{value}

	for _, segment := range jp.segments {
		if len(current) == 0 {
			break
		}

		var nextCurrent []types.JSONValue

		for _, val := range current {
			results, err := segment.apply(val)
			if err != nil {
				return nil, err
			}

			nextCurrent = append(nextCurrent, results...)
		}

		current = nextCurrent
	}

	return current, nil
}

// String 返回JSON Path的字符串表示
func (jp *JSONPath) String() string {
	var sb strings.Builder

	for _, segment := range jp.segments {
		sb.WriteString(segment.String())
	}

	return sb.String()
}

// QueryJSONPath 使用JSON Path查询JSON值
func QueryJSONPath(value types.JSONValue, pathExpr string) ([]types.JSONValue, error) {
	path, err := ParseJSONPath(pathExpr)
	if err != nil {
		return nil, err
	}

	return path.Query(value)
}

// QueryJSONPathString 使用JSON Path查询JSON字符串
func QueryJSONPathString(jsonStr string, pathExpr string) ([]types.JSONValue, error) {
	value, err := parser.ParseToValue(jsonStr)
	if err != nil {
		return nil, err
	}

	return QueryJSONPath(value, pathExpr)
}
