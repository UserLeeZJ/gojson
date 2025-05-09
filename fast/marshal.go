// Package fast 提供高性能的JSON序列化和反序列化功能。
package fast

import (
	"bytes"
	"encoding/json"
	"strconv"
	"sync"
	"unsafe"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
)

// 预分配的缓冲区大小常量。
const (
	defaultBufSize = 4096        // 默认缓冲区大小。
	maxBufSize     = 1024 * 1024 // 最大缓冲区大小，防止内存泄漏。
)

// 错误代码常量。
const (
	ErrInvalidJSON = jsonerrors.ErrInvalidJSON
	ErrEmptyInput  = jsonerrors.ErrEmptyInput
)

// bufferPool 是用于JSON编码/解码的缓冲池。
var bufferPool = sync.Pool{
	New: func() interface{} {
		// 预分配一个合理大小的缓冲区，减少扩容操作。
		return bytes.NewBuffer(make([]byte, 0, defaultBufSize))
	},
}

// getBuffer 获取缓冲区。
func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// releaseBuffer 释放缓冲区。
func releaseBuffer(buf *bytes.Buffer) {
	// 如果缓冲区容量过大，不放回池中，让GC回收。
	if buf.Cap() > maxBufSize {
		return
	}
	bufferPool.Put(buf)
}

// digitStrings 用于存储预计算的数字字符串值。
var digitStrings [1000]string

// init 初始化预计算的数字字符串。
func init() {
	// 预计算常用数字的字符串表示，避免频繁转换。
	for i := 0; i < 1000; i++ {
		digitStrings[i] = strconv.Itoa(i)
	}
}

// fastItoa 快速获取数字的字符串表示。
func fastItoa(i int) string {
	if i >= 0 && i < 1000 {
		return digitStrings[i]
	}
	return strconv.Itoa(i)
}

// bytesToString 安全地将[]byte转换为string，不需要内存复制。
// 注意：返回的字符串不应该在输入字节切片被修改后使用。
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// stringToBytes 安全地将string转换为[]byte，不需要内存复制。
// 注意：返回的字节切片不应该被修改。
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// Marshal 是一个优化的JSON序列化函数。
func Marshal(v interface{}) ([]byte, error) {
	// 对于nil值，直接返回"null"。
	if v == nil {
		return []byte("null"), nil
	}

	// 快速路径：处理简单类型。
	switch val := v.(type) {
	case string:
		// 字符串需要特殊处理，添加引号和转义。
		return marshalString(val)
	case bool:
		if val {
			return []byte("true"), nil
		}
		return []byte("false"), nil
	case int:
		if val >= 0 && val < 1000 {
			return stringToBytes(digitStrings[val]), nil
		}
		return stringToBytes(strconv.FormatInt(int64(val), 10)), nil
	case int64:
		if val >= 0 && val < 1000 {
			return stringToBytes(digitStrings[val]), nil
		}
		return stringToBytes(strconv.FormatInt(val, 10)), nil
	case float64:
		return stringToBytes(strconv.FormatFloat(val, 'f', -1, 64)), nil
	case []byte:
		// 对于[]byte，我们需要base64编码，使用标准库。
		return json.Marshal(val)
	case map[string]interface{}:
		// 对于小型map，使用优化的方法。
		if len(val) < 10 {
			return marshalSmallMap(val)
		}
	case []interface{}:
		// 对于小型数组，使用优化的方法。
		if len(val) < 10 {
			return marshalSmallArray(val)
		}
	}

	// 获取缓冲区。
	buf := getBuffer()
	defer releaseBuffer(buf)

	// 创建一个新的编码器。
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false) // 避免HTML转义，保持JSON原样。
	enc.SetIndent("", "")    // 默认不缩进。

	// 编码数据。
	if err := enc.Encode(v); err != nil {
		return nil, jsonerrors.NewJSONError(ErrInvalidJSON, "序列化失败").WithCause(err)
	}

	// 获取缓冲区内容。
	bytes := buf.Bytes()

	// 处理末尾的换行符（json.Encoder总是添加一个换行符）。
	if len(bytes) > 0 && bytes[len(bytes)-1] == '\n' {
		bytes = bytes[:len(bytes)-1]
	}

	// 创建结果，避免后续修改缓冲区影响结果。
	result := make([]byte, len(bytes))
	copy(result, bytes)

	return result, nil
}

// marshalString 将字符串转换为JSON字符串。
func marshalString(s string) ([]byte, error) {
	// 快速路径：空字符串。
	if s == "" {
		return []byte(`""`), nil
	}

	// 快速路径：检查是否需要转义。
	needEscape := false
	for i := 0; i < len(s); i++ {
		// ASCII控制字符或需要转义的字符。
		if s[i] < 32 || s[i] == '"' || s[i] == '\\' {
			needEscape = true
			break
		}
	}

	// 如果不需要转义，直接添加引号。
	if !needEscape {
		result := make([]byte, len(s)+2)
		result[0] = '"'
		copy(result[1:], s)
		result[len(result)-1] = '"'
		return result, nil
	}

	// 需要转义，使用标准库。
	return json.Marshal(s)
}

// marshalSmallMap 优化小型map的序列化。
func marshalSmallMap(m map[string]interface{}) ([]byte, error) {
	if len(m) == 0 {
		return []byte("{}"), nil
	}

	buf := getBuffer()
	defer releaseBuffer(buf)

	buf.WriteByte('{')
	first := true

	for k, v := range m {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		// 写入键。
		keyBytes, err := marshalString(k)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)
		buf.WriteByte(':')

		// 写入值。
		valBytes, err := Marshal(v)
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}

	buf.WriteByte('}')

	// 复制结果。
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())

	return result, nil
}

// marshalSmallArray 优化小型数组的序列化。
func marshalSmallArray(arr []interface{}) ([]byte, error) {
	if len(arr) == 0 {
		return []byte("[]"), nil
	}

	buf := getBuffer()
	defer releaseBuffer(buf)

	buf.WriteByte('[')

	for i, v := range arr {
		if i > 0 {
			buf.WriteByte(',')
		}

		// 写入值。
		valBytes, err := Marshal(v)
		if err != nil {
			return nil, err
		}
		buf.Write(valBytes)
	}

	buf.WriteByte(']')

	// 复制结果。
	result := make([]byte, buf.Len())
	copy(result, buf.Bytes())

	return result, nil
}
