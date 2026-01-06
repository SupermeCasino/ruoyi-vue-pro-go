package dto

import (
	"encoding/json"
	"strconv"
	"strings"
)

// Flexible 支持从 JSON 数字或字符串解析的泛型类型
type Flexible[T any] struct {
	Value T
}

func (f *Flexible[T]) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}

	s := strings.Trim(string(data), "\"")
	if s == "" {
		return nil
	}

	var target any = &f.Value
	switch v := target.(type) {
	case *float64:
		val, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		*v = val
	case *int:
		val, err := strconv.Atoi(s)
		if err != nil {
			// 兼容 "100.0" 这种带小数点的字符串转 int
			if fval, ferr := strconv.ParseFloat(s, 64); ferr == nil {
				*v = int(fval)
				return nil
			}
			return err
		}
		*v = val
	case *int64:
		val, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			if fval, ferr := strconv.ParseFloat(s, 64); ferr == nil {
				*v = int64(fval)
				return nil
			}
			return err
		}
		*v = val
	default:
		// 其他类型回退到标准 json 解析
		return json.Unmarshal(data, &f.Value)
	}
	return nil
}

func (f Flexible[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

// ThingModelDataSpecs 物模型数据规范
type ThingModelDataSpecs struct {
	DataType string `json:"dataType"`
	// 数值相关
	Min          *Flexible[float64] `json:"min,omitempty"`
	Max          *Flexible[float64] `json:"max,omitempty"`
	Step         *Flexible[float64] `json:"step,omitempty"`
	Precise      string             `json:"precise,omitempty"`
	DefaultValue string             `json:"defaultValue,omitempty"`
	Unit         string             `json:"unit,omitempty"`
	UnitName     string             `json:"unitName,omitempty"`
	// 文本相关
	MaxLength *Flexible[int] `json:"maxLength,omitempty"`
	// 枚举/布尔相关 (对齐 Java: ThingModelBoolOrEnumDataSpecs)
	Name  string         `json:"name,omitempty"`
	Value *Flexible[int] `json:"value,omitempty"`
	// 结构体成员相关 (对齐 Java: ThingModelStructDataSpecs)
	Identifier string `json:"identifier,omitempty"`
	AccessMode string `json:"accessMode,omitempty"`
	Required   bool   `json:"required,omitempty"`
	// 递归支持 (结构体/数组)
	ChildDataType string                `json:"childDataType,omitempty"`
	Size          *Flexible[int]        `json:"size,omitempty"`
	DataSpecs     *ThingModelDataSpecs  `json:"dataSpecs,omitempty"`
	DataSpecsList []ThingModelDataSpecs `json:"dataSpecsList,omitempty"`
}

// ThingModelProperty 物模型属性
type ThingModelProperty struct {
	Identifier    string                `json:"identifier"`
	Name          string                `json:"name"`
	AccessMode    string                `json:"accessMode"`
	Required      bool                  `json:"required"`
	DataType      string                `json:"dataType"`
	DataSpecs     *ThingModelDataSpecs  `json:"dataSpecs,omitempty"`
	DataSpecsList []ThingModelDataSpecs `json:"dataSpecsList,omitempty"`
}

// ThingModelParam 物模型参数
type ThingModelParam struct {
	Identifier    string                `json:"identifier"`
	Name          string                `json:"name"`
	Direction     string                `json:"direction"`
	ParaOrder     *Flexible[int]        `json:"paraOrder,omitempty"`
	DataType      string                `json:"dataType"`
	DataSpecs     *ThingModelDataSpecs  `json:"dataSpecs,omitempty"`
	DataSpecsList []ThingModelDataSpecs `json:"dataSpecsList,omitempty"`
}

// ThingModelEvent 物模型事件
type ThingModelEvent struct {
	Identifier   string            `json:"identifier"`
	Name         string            `json:"name"`
	Required     bool              `json:"required"`
	Type         string            `json:"type"`
	OutputParams []ThingModelParam `json:"outputParams,omitempty"`
	Method       string            `json:"method,omitempty"`
}

// ThingModelService 物模型服务
type ThingModelService struct {
	Identifier   string            `json:"identifier"`
	Name         string            `json:"name"`
	Required     bool              `json:"required"`
	CallType     string            `json:"callType"`
	InputParams  []ThingModelParam `json:"inputParams,omitempty"`
	OutputParams []ThingModelParam `json:"outputParams,omitempty"`
	Method       string            `json:"method,omitempty"`
}
