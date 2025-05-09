package stream

import (
	"bytes"
	"encoding/json"
	"sync"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
	"github.com/UserLeeZJ/gojson/types"
)

// IncrementalParser 是增量JSON解析器
type IncrementalParser struct {
	buffer     bytes.Buffer
	complete   bool
	err        error
	result     interface{}
	bufferLock sync.Mutex
}

// NewIncrementalParser 创建一个新的增量JSON解析器
func NewIncrementalParser() *IncrementalParser {
	return &IncrementalParser{
		complete: false,
	}
}

// Feed 向解析器提供更多的JSON数据
func (p *IncrementalParser) Feed(data []byte) error {
	p.bufferLock.Lock()
	defer p.bufferLock.Unlock()

	if p.complete {
		return jsonerrors.NewJSONError(ErrInvalidJSON, "解析已完成，无法提供更多数据")
	}

	if p.err != nil {
		return p.err
	}

	// 将数据添加到缓冲区
	_, err := p.buffer.Write(data)
	if err != nil {
		p.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入缓冲区失败").WithCause(err)
		return p.err
	}

	// 尝试解析完整的JSON
	jsonData := p.buffer.Bytes()
	
	// 检查JSON是否完整
	if isCompleteJSON(jsonData) {
		var result interface{}
		err := json.Unmarshal(jsonData, &result)
		if err != nil {
			// 可能是不完整的JSON，继续等待更多数据
			return nil
		}
		
		// 解析成功，标记为完成
		p.result = result
		p.complete = true
	}
	
	return nil
}

// isCompleteJSON 检查JSON数据是否完整
func isCompleteJSON(data []byte) bool {
	// 跳过前导空白
	i := 0
	for i < len(data) && (data[i] == ' ' || data[i] == '\t' || data[i] == '\n' || data[i] == '\r') {
		i++
	}
	
	// 空数据
	if i >= len(data) {
		return false
	}
	
	// 检查第一个非空白字符
	firstChar := data[i]
	
	// 简单检查：如果是对象或数组，检查是否有匹配的结束括号
	if firstChar == '{' {
		// 对象
		depth := 1
		inString := false
		escaped := false
		
		for i++; i < len(data); i++ {
			c := data[i]
			
			if inString {
				if escaped {
					escaped = false
				} else if c == '\\' {
					escaped = true
				} else if c == '"' {
					inString = false
				}
			} else {
				if c == '"' {
					inString = true
				} else if c == '{' {
					depth++
				} else if c == '}' {
					depth--
					if depth == 0 {
						// 找到匹配的结束括号
						return true
					}
				}
			}
		}
		
		// 没有找到匹配的结束括号
		return false
	} else if firstChar == '[' {
		// 数组
		depth := 1
		inString := false
		escaped := false
		
		for i++; i < len(data); i++ {
			c := data[i]
			
			if inString {
				if escaped {
					escaped = false
				} else if c == '\\' {
					escaped = true
				} else if c == '"' {
					inString = false
				}
			} else {
				if c == '"' {
					inString = true
				} else if c == '[' {
					depth++
				} else if c == ']' {
					depth--
					if depth == 0 {
						// 找到匹配的结束括号
						return true
					}
				}
			}
		}
		
		// 没有找到匹配的结束括号
		return false
	} else if firstChar == '"' {
		// 字符串
		escaped := false
		
		for i++; i < len(data); i++ {
			c := data[i]
			
			if escaped {
				escaped = false
			} else if c == '\\' {
				escaped = true
			} else if c == '"' {
				// 找到结束引号
				return true
			}
		}
		
		// 没有找到结束引号
		return false
	} else if firstChar == 't' || firstChar == 'f' || firstChar == 'n' || 
		(firstChar >= '0' && firstChar <= '9') || firstChar == '-' {
		// 简单值（true, false, null, 数字）
		// 尝试使用标准库解析
		var dummy interface{}
		return json.Unmarshal(data, &dummy) == nil
	}
	
	// 无法识别的JSON开始字符
	return false
}

// Result 返回当前解析结果
func (p *IncrementalParser) Result() (interface{}, error) {
	if p.err != nil {
		return nil, p.err
	}

	if !p.complete {
		return nil, jsonerrors.NewJSONError(ErrInvalidJSON, "解析尚未完成")
	}

	return p.result, nil
}

// ResultValue 返回当前解析结果作为JSONValue
func (p *IncrementalParser) ResultValue() (types.JSONValue, error) {
	result, err := p.Result()
	if err != nil {
		return nil, err
	}

	return types.FromGoValue(result)
}

// IsComplete 返回解析是否已完成
func (p *IncrementalParser) IsComplete() bool {
	return p.complete
}

// Error 返回解析错误
func (p *IncrementalParser) Error() error {
	return p.err
}

// Reset 重置解析器状态
func (p *IncrementalParser) Reset() {
	p.bufferLock.Lock()
	defer p.bufferLock.Unlock()

	p.buffer.Reset()
	p.result = nil
	p.complete = false
	p.err = nil
}
