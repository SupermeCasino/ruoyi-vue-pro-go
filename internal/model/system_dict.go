package model

import (
	"time"
)

// SystemDictType 字典类型表
type SystemDictType struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"column:name;not null;default:''" json:"name"`
	Type      string    `gorm:"column:type;not null;default:''" json:"type"`
	Status    int32     `gorm:"column:status;not null;default:0" json:"status"` // 0:开启, 1:禁用
	Remark    string    `gorm:"column:remark;default:''" json:"remark"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
}

func (SystemDictType) TableName() string {
	return "system_dict_type"
}

// SystemDictData 字典数据表
type SystemDictData struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Sort      int32     `gorm:"column:sort;not null;default:0" json:"sort"`
	Label     string    `gorm:"column:label;not null;default:''" json:"label"`
	Value     string    `gorm:"column:value;not null;default:''" json:"value"`
	DictType  string    `gorm:"column:dict_type;not null;default:''" json:"dictType"`
	Status    int32     `gorm:"column:status;not null;default:0" json:"status"` // 0:开启, 1:禁用
	ColorType string    `gorm:"column:color_type;default:''" json:"colorType"`
	CssClass  string    `gorm:"column:css_class;default:''" json:"cssClass"`
	Remark    string    `gorm:"column:remark;default:''" json:"remark"`
	Creator   string    `gorm:"column:creator;default:''" json:"creator"`
	Updater   string    `gorm:"column:updater;default:''" json:"updater"`
	CreatedAt time.Time `gorm:"column:create_time;autoCreateTime" json:"createTime"`
	UpdatedAt time.Time `gorm:"column:update_time;autoUpdateTime" json:"updateTime"`
	Deleted   BitBool   `gorm:"column:deleted;softDelete:flag" json:"-"`
}

func (SystemDictData) TableName() string {
	return "system_dict_data"
}
