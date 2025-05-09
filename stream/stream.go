// Package stream 提供gojson库的流式JSON处理功能
package stream

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
)

// 错误代码
const (
	ErrInvalidJSON = jsonerrors.ErrInvalidJSON
	ErrEmptyInput  = jsonerrors.ErrEmptyInput
)

// 缓冲区大小
const (
	defaultBufSize = 4096 // 默认缓冲区大小
)

// JSONTokenType 表示JSON令牌的类型
type JSONTokenType int

const (
	// TokenError 表示解析错误
	TokenError JSONTokenType = iota
	// TokenObjectStart 表示对象开始 {
	TokenObjectStart
	// TokenObjectEnd 表示对象结束 }
	TokenObjectEnd
	// TokenArrayStart 表示数组开始 [
	TokenArrayStart
	// TokenArrayEnd 表示数组结束 ]
	TokenArrayEnd
	// TokenPropertyName 表示属性名
	TokenPropertyName
	// TokenString 表示字符串值
	TokenString
	// TokenNumber 表示数字值
	TokenNumber
	// TokenBoolean 表示布尔值
	TokenBoolean
	// TokenNull 表示null值
	TokenNull
	// TokenEOF 表示输入结束
	TokenEOF
)

// JSONToken 表示JSON令牌
type JSONToken struct {
	// Type 是令牌的类型
	Type JSONTokenType
	// Value 是令牌的值
	Value interface{}
	// Depth 是令牌的深度
	Depth int
	// Path 是令牌的路径
	Path string
	// Error 是解析错误
	Error error
}

// JSONTokenizer 是JSON流式解析器
type JSONTokenizer struct {
	reader    *bufio.Reader
	buffer    bytes.Buffer
	depth     int
	path      []string
	lastToken JSONToken
	err       error
}

// NewJSONTokenizer 创建一个新的JSON流式解析器
func NewJSONTokenizer(r io.Reader) *JSONTokenizer {
	return &JSONTokenizer{
		reader: bufio.NewReaderSize(r, defaultBufSize),
		depth:  0,
		path:   make([]string, 0),
	}
}

// Next 返回下一个JSON令牌
func (t *JSONTokenizer) Next() JSONToken {
	// 如果已经有错误，直接返回错误令牌
	if t.err != nil {
		return JSONToken{Type: TokenError, Error: t.err}
	}

	// 读取下一个非空白字符
	c, err := t.readNonWhitespace()
	if err != nil {
		if err == io.EOF {
			return JSONToken{Type: TokenEOF}
		}
		t.err = err
		return JSONToken{Type: TokenError, Error: err}
	}

	// 根据字符类型解析令牌
	switch c {
	case '{':
		t.depth++
		return JSONToken{Type: TokenObjectStart, Depth: t.depth, Path: t.currentPath()}
	case '}':
		t.depth--
		return JSONToken{Type: TokenObjectEnd, Depth: t.depth, Path: t.currentPath()}
	case '[':
		t.depth++
		return JSONToken{Type: TokenArrayStart, Depth: t.depth, Path: t.currentPath()}
	case ']':
		t.depth--
		return JSONToken{Type: TokenArrayEnd, Depth: t.depth, Path: t.currentPath()}
	case ',':
		// 跳过逗号，读取下一个令牌
		return t.Next()
	case ':':
		// 跳过冒号，读取下一个令牌
		return t.Next()
	case '"':
		// 解析字符串
		value, err := t.parseString()
		if err != nil {
			t.err = err
			return JSONToken{Type: TokenError, Error: err}
		}

		// 检查是否为属性名
		nextChar, err := t.peekNextNonWhitespace()
		if err == nil && nextChar == ':' {
			// 消耗冒号
			_, _ = t.readNonWhitespace()
			return JSONToken{Type: TokenPropertyName, Value: value, Depth: t.depth, Path: t.currentPath()}
		}
		return JSONToken{Type: TokenString, Value: value, Depth: t.depth, Path: t.currentPath()}
	case 't':
		// 解析true
		if err := t.expectString("rue"); err != nil {
			t.err = err
			return JSONToken{Type: TokenError, Error: err}
		}
		return JSONToken{Type: TokenBoolean, Value: true, Depth: t.depth, Path: t.currentPath()}
	case 'f':
		// 解析false
		if err := t.expectString("alse"); err != nil {
			t.err = err
			return JSONToken{Type: TokenError, Error: err}
		}
		return JSONToken{Type: TokenBoolean, Value: false, Depth: t.depth, Path: t.currentPath()}
	case 'n':
		// 解析null
		if err := t.expectString("ull"); err != nil {
			t.err = err
			return JSONToken{Type: TokenError, Error: err}
		}
		return JSONToken{Type: TokenNull, Depth: t.depth, Path: t.currentPath()}
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		// 解析数字
		value, err := t.parseNumber(c)
		if err != nil {
			t.err = err
			return JSONToken{Type: TokenError, Error: err}
		}
		return JSONToken{Type: TokenNumber, Value: value, Depth: t.depth, Path: t.currentPath()}
	default:
		t.err = jsonerrors.NewJSONError(ErrInvalidJSON, "无效的JSON字符")
		return JSONToken{Type: TokenError, Error: t.err}
	}
}

// 读取下一个非空白字符
func (t *JSONTokenizer) readNonWhitespace() (byte, error) {
	for {
		c, err := t.reader.ReadByte()
		if err != nil {
			return 0, err
		}
		if !isWhitespace(c) {
			return c, nil
		}
	}
}

// 解析字符串
func (t *JSONTokenizer) parseString() (string, error) {
	// 直接使用标准库的方式解析JSON字符串
	var sb bytes.Buffer
	sb.WriteByte('"') // 添加开始引号

	escaped := false
	for {
		c, err := t.reader.ReadByte()
		if err != nil {
			return "", jsonerrors.NewJSONError(ErrInvalidJSON, "解析字符串时遇到EOF")
		}

		// 添加字符到缓冲区
		sb.WriteByte(c)

		// 处理转义字符
		if escaped {
			escaped = false
		} else if c == '\\' {
			escaped = true
		} else if c == '"' {
			break
		}
	}

	// 使用标准库解析JSON字符串
	var result string
	err := json.Unmarshal(sb.Bytes(), &result)
	if err != nil {
		return "", jsonerrors.NewJSONError(ErrInvalidJSON, "解析字符串失败").WithCause(err)
	}

	return result, nil
}

// 解析布尔值
func (t *JSONTokenizer) parseBoolean(first byte) (bool, error) {
	if first == 't' {
		// 期望 "true"
		expected := "rue"
		for i := 0; i < len(expected); i++ {
			c, err := t.reader.ReadByte()
			if err != nil {
				return false, jsonerrors.NewJSONError(ErrInvalidJSON, "解析布尔值时遇到EOF")
			}
			if c != expected[i] {
				return false, jsonerrors.NewJSONError(ErrInvalidJSON, "无效的布尔值")
			}
		}
		return true, nil
	} else {
		// 期望 "false"
		expected := "alse"
		for i := 0; i < len(expected); i++ {
			c, err := t.reader.ReadByte()
			if err != nil {
				return false, jsonerrors.NewJSONError(ErrInvalidJSON, "解析布尔值时遇到EOF")
			}
			if c != expected[i] {
				return false, jsonerrors.NewJSONError(ErrInvalidJSON, "无效的布尔值")
			}
		}
		return false, nil
	}
}

// 解析null
func (t *JSONTokenizer) parseNull() error {
	// 期望 "null"
	expected := "ull"
	for i := 0; i < len(expected); i++ {
		c, err := t.reader.ReadByte()
		if err != nil {
			return jsonerrors.NewJSONError(ErrInvalidJSON, "解析null时遇到EOF")
		}
		if c != expected[i] {
			return jsonerrors.NewJSONError(ErrInvalidJSON, "无效的null值")
		}
	}
	return nil
}

// 解析数字
func (t *JSONTokenizer) parseNumber(first byte) (json.Number, error) {
	var sb bytes.Buffer
	sb.WriteByte(first)

	for {
		c, err := t.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", jsonerrors.NewJSONError(ErrInvalidJSON, "解析数字时遇到错误").WithCause(err)
		}

		if isDigit(c) || c == '.' || c == 'e' || c == 'E' || c == '+' || c == '-' {
			sb.WriteByte(c)
		} else {
			// 将字符放回缓冲区
			t.reader.UnreadByte()
			break
		}
	}

	// 验证数字格式
	numStr := sb.String()
	if !isValidNumber(numStr) {
		return "", jsonerrors.NewJSONError(ErrInvalidJSON, "无效的数字格式")
	}

	return json.Number(numStr), nil
}

// 检查是否为属性名
func (t *JSONTokenizer) isPropertyName() bool {
	// 读取下一个非空白字符
	c, err := t.readNonWhitespace()
	if err != nil {
		return false
	}

	// 将字符放回缓冲区
	t.reader.UnreadByte()

	// 如果下一个字符是冒号，则当前字符串是属性名
	return c == ':'
}

// peekNextNonWhitespace 查看下一个非空白字符但不消耗它
func (t *JSONTokenizer) peekNextNonWhitespace() (byte, error) {
	// 读取下一个非空白字符
	c, err := t.readNonWhitespace()
	if err != nil {
		return 0, err
	}

	// 将字符放回缓冲区
	t.reader.UnreadByte()

	return c, nil
}

// expectString 期望读取指定的字符串
func (t *JSONTokenizer) expectString(expected string) error {
	for i := 0; i < len(expected); i++ {
		c, err := t.reader.ReadByte()
		if err != nil {
			return jsonerrors.NewJSONError(ErrInvalidJSON, "读取字符时遇到EOF")
		}
		if c != expected[i] {
			return jsonerrors.NewJSONError(ErrInvalidJSON, "无效的字符序列")
		}
	}
	return nil
}

// 获取当前路径
func (t *JSONTokenizer) currentPath() string {
	if len(t.path) == 0 {
		return "$"
	}
	var sb bytes.Buffer
	sb.WriteString("$")
	for _, p := range t.path {
		sb.WriteString(p)
	}
	return sb.String()
}

// 检查字符是否为空白
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// 检查字符是否为数字
func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

// 检查字符是否为十六进制数字
func isHexDigit(c byte) bool {
	return isDigit(c) || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
}

// 检查字符串是否为有效的数字
func isValidNumber(s string) bool {
	// 简单检查，可以使用更复杂的正则表达式
	_, err := json.Number(s).Float64()
	return err == nil
}
