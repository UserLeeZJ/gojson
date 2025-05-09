package stream

import (
	"bufio"
	"io"
	"strconv"
	"sync"

	jsonerrors "github.com/UserLeeZJ/gojson/errors"
)

// JSONGenerator 是JSON流式生成器
type JSONGenerator struct {
	writer     *bufio.Writer
	depth      int
	states     []generatorState
	needComma  bool
	err        error
	writeMutex sync.Mutex
}

// generatorState 表示生成器的状态
type generatorState int

const (
	// stateNone 表示初始状态
	stateNone generatorState = iota
	// stateObject 表示正在生成对象
	stateObject
	// stateArray 表示正在生成数组
	stateArray
)

// NewJSONGenerator 创建一个新的JSON流式生成器
func NewJSONGenerator(w io.Writer) *JSONGenerator {
	return &JSONGenerator{
		writer: bufio.NewWriter(w),
		depth:  0,
		states: make([]generatorState, 0, 10),
	}
}

// BeginObject 开始一个新的对象
func (g *JSONGenerator) BeginObject() error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if g.needComma {
		if err := g.writeComma(); err != nil {
			return err
		}
	}

	if err := g.writeByte('{'); err != nil {
		return err
	}

	g.states = append(g.states, stateObject)
	g.depth++
	g.needComma = false

	return nil
}

// EndObject 结束当前对象
func (g *JSONGenerator) EndObject() error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if len(g.states) == 0 || g.states[len(g.states)-1] != stateObject {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "尝试结束不存在的对象")
		return g.err
	}

	if err := g.writeByte('}'); err != nil {
		return err
	}

	g.states = g.states[:len(g.states)-1]
	g.depth--
	g.needComma = true

	return nil
}

// BeginArray 开始一个新的数组
func (g *JSONGenerator) BeginArray() error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if g.needComma {
		if err := g.writeComma(); err != nil {
			return err
		}
	}

	if err := g.writeByte('['); err != nil {
		return err
	}

	g.states = append(g.states, stateArray)
	g.depth++
	g.needComma = false

	return nil
}

// EndArray 结束当前数组
func (g *JSONGenerator) EndArray() error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if len(g.states) == 0 || g.states[len(g.states)-1] != stateArray {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "尝试结束不存在的数组")
		return g.err
	}

	if err := g.writeByte(']'); err != nil {
		return err
	}

	g.states = g.states[:len(g.states)-1]
	g.depth--
	g.needComma = true

	return nil
}

// WriteProperty 写入一个属性名
func (g *JSONGenerator) WriteProperty(name string) error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if len(g.states) == 0 || g.states[len(g.states)-1] != stateObject {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "属性名只能在对象中使用")
		return g.err
	}

	if g.needComma {
		if err := g.writeComma(); err != nil {
			return err
		}
	}

	// 写入属性名
	if err := g.writeString(name); err != nil {
		return err
	}

	// 写入冒号
	if err := g.writeByte(':'); err != nil {
		return err
	}

	g.needComma = false

	return nil
}

// WriteString 写入一个字符串值
func (g *JSONGenerator) WriteString(value string) error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if g.needComma {
		if err := g.writeComma(); err != nil {
			return err
		}
	}

	if err := g.writeString(value); err != nil {
		return err
	}

	g.needComma = true

	return nil
}

// WriteNumber 写入一个数字值
func (g *JSONGenerator) WriteNumber(value float64) error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if g.needComma {
		if err := g.writeComma(); err != nil {
			return err
		}
	}

	// 转换为字符串
	str := strconv.FormatFloat(value, 'f', -1, 64)

	// 写入数字
	if _, err := g.writer.WriteString(str); err != nil {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入数字失败").WithCause(err)
		return g.err
	}

	g.needComma = true

	return nil
}

// WriteBoolean 写入一个布尔值
func (g *JSONGenerator) WriteBoolean(value bool) error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if g.needComma {
		if err := g.writeComma(); err != nil {
			return err
		}
	}

	// 写入布尔值
	var str string
	if value {
		str = "true"
	} else {
		str = "false"
	}

	if _, err := g.writer.WriteString(str); err != nil {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入布尔值失败").WithCause(err)
		return g.err
	}

	g.needComma = true

	return nil
}

// WriteNull 写入一个null值
func (g *JSONGenerator) WriteNull() error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if g.needComma {
		if err := g.writeComma(); err != nil {
			return err
		}
	}

	// 写入null
	if _, err := g.writer.WriteString("null"); err != nil {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入null失败").WithCause(err)
		return g.err
	}

	g.needComma = true

	return nil
}

// Flush 刷新缓冲区
func (g *JSONGenerator) Flush() error {
	g.writeMutex.Lock()
	defer g.writeMutex.Unlock()

	if g.err != nil {
		return g.err
	}

	if err := g.writer.Flush(); err != nil {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "刷新缓冲区失败").WithCause(err)
		return g.err
	}

	return nil
}

// 写入一个字节
func (g *JSONGenerator) writeByte(b byte) error {
	if err := g.writer.WriteByte(b); err != nil {
		g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字节失败").WithCause(err)
		return g.err
	}
	return nil
}

// 写入逗号
func (g *JSONGenerator) writeComma() error {
	return g.writeByte(',')
}

// 写入字符串（带引号和转义）
func (g *JSONGenerator) writeString(s string) error {
	// 写入开始引号
	if err := g.writeByte('"'); err != nil {
		return err
	}

	// 写入字符串内容（需要处理转义）
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '"', '\\', '/':
			if err := g.writeByte('\\'); err != nil {
				return err
			}
			if err := g.writeByte(c); err != nil {
				return err
			}
		case '\b':
			if _, err := g.writer.WriteString("\\b"); err != nil {
				g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字符串失败").WithCause(err)
				return g.err
			}
		case '\f':
			if _, err := g.writer.WriteString("\\f"); err != nil {
				g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字符串失败").WithCause(err)
				return g.err
			}
		case '\n':
			if _, err := g.writer.WriteString("\\n"); err != nil {
				g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字符串失败").WithCause(err)
				return g.err
			}
		case '\r':
			if _, err := g.writer.WriteString("\\r"); err != nil {
				g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字符串失败").WithCause(err)
				return g.err
			}
		case '\t':
			if _, err := g.writer.WriteString("\\t"); err != nil {
				g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字符串失败").WithCause(err)
				return g.err
			}
		default:
			if c < 32 {
				// 控制字符需要使用\uXXXX格式
				if _, err := g.writer.WriteString("\\u00"); err != nil {
					g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字符串失败").WithCause(err)
					return g.err
				}
				if _, err := g.writer.WriteString(strconv.FormatInt(int64(c), 16)); err != nil {
					g.err = jsonerrors.NewJSONError(ErrInvalidJSON, "写入字符串失败").WithCause(err)
					return g.err
				}
			} else {
				if err := g.writeByte(c); err != nil {
					return err
				}
			}
		}
	}

	// 写入结束引号
	if err := g.writeByte('"'); err != nil {
		return err
	}

	return nil
}
