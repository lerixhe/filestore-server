package util

import (
	"encoding/json"
	"fmt"
)

// ResMsg http响应的通用结构
type ResMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// NewResMsg 生成response对象
func NewResMsg(code int, msg string, data interface{}) *ResMsg {
	return &ResMsg{code, msg, data}
}

// JSONString 对象转json string
func (r *ResMsg) JSONString() string {
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(data)
}

// JSONBytes 对象转json bytes
func (r *ResMsg) JSONBytes() []byte {
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return data
}

// GenSimpleResStream 生成仅含code和msg的响应体，返回[]byte
func GenSimpleResStream(code int, msg string) []byte {
	return []byte(fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg))
}

// GenSimpleResString 生成仅含code和msg的响应体，返回[]byte
func GenSimpleResString(code int, msg string) string {
	return fmt.Sprintf(`{"code":%d,"msg":"%s"}`, code, msg)
}
