package types

import (
	"encoding/json"
	"errors"
	"strconv"
)

// FlexInt64 支持从数字和字符串反序列化
// 当前端发送 "123" 而不是 123 时很有用
// 同时在序列化为 JSON 时会转为字符串，防止前端精度丢失
type FlexInt64 int64

func (i *FlexInt64) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		*i = 0
		return nil
	}
	// 尝试直接作为数字解析
	var n int64
	if err := json.Unmarshal(data, &n); err == nil {
		*i = FlexInt64(n)
		return nil
	}

	// 尝试作为字符串解析
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		p, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		*i = FlexInt64(p)
		return nil
	}

	return errors.New("FlexInt64: 无法解析为 int64 或字符串")
}

func (i FlexInt64) MarshalJSON() ([]byte, error) {
	// 转为字符串输出，防止前端 ID 溢出
	return json.Marshal(strconv.FormatInt(int64(i), 10))
}
