package resp

import "time"

// DictTypeSimpleRespVO 字典类型精简信息
type DictTypeSimpleRespVO struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

// DictTypeRespVO 字典类型详细信息
type DictTypeRespVO struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Type       string    `json:"type"`
	Status     int32     `json:"status"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}

// DictDataSimpleRespVO 字典数据精简信息
type DictDataSimpleRespVO struct {
	DictType  string `json:"dictType"`
	Value     string `json:"value"`
	Label     string `json:"label"`
	ColorType string `json:"colorType"`
	CssClass  string `json:"cssClass"`
}

// DictDataRespVO 字典数据详细信息
type DictDataRespVO struct {
	ID         int64     `json:"id"`
	Sort       int32     `json:"sort"`
	Label      string    `json:"label"`
	Value      string    `json:"value"`
	DictType   string    `json:"dictType"`
	Status     int32     `json:"status"`
	ColorType  string    `json:"colorType"`
	CssClass   string    `json:"cssClass"`
	Remark     string    `json:"remark"`
	CreateTime time.Time `json:"createTime"`
}
