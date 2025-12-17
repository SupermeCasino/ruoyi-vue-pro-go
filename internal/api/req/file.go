package req

import (
	"github.com/wxlbd/ruoyi-mall-go/pkg/pagination"
)

// FileConfigSaveReq 文件配置创建/修改 Request
type FileConfigSaveReq struct {
	ID      int64                  `json:"id"`
	Name    string                 `json:"name" binding:"required"`
	Storage int32                  `json:"storage" binding:"required"` // 参见 FileStorageEnum
	Config  map[string]interface{} `json:"config" binding:"required"`
	Remark  string                 `json:"remark"`
}

// FileConfigPageReq 文件配置分页 Request
type FileConfigPageReq struct {
	pagination.PageParam
	Name       string   `form:"name"`
	Storage    *int32   `form:"storage"`
	CreateTime []string `form:"createTime[]"`
}

// FilePageReq 文件分页 Request
type FilePageReq struct {
	pagination.PageParam
	Path       string   `form:"path"`
	Type       string   `form:"type"`
	CreateTime []string `form:"createTime[]"`
}

// FileUploadReq 上传文件 Request (无需 JSON binding，直接从 Form 获取)
type FileUploadReq struct {
	Path string `form:"path"` // 自定义上传路径/文件名
}

// FileCreateReq 创建文件 Request (前端直传回调)
type FileCreateReq struct {
	ConfigID int64  `json:"configId" binding:"required"`
	Path     string `json:"path" binding:"required"`
	Name     string `json:"name" binding:"required"`
	URL      string `json:"url" binding:"required"`
	Type     string `json:"type"`
	Size     int    `json:"size" binding:"required"`
}
