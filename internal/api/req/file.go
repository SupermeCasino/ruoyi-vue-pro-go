package req

import "backend-go/internal/pkg/core"

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
	core.PageParam
	Name       string   `form:"name"`
	Storage    *int32   `form:"storage"`
	CreateTime []string `form:"createTime[]"`
}

// FilePageReq 文件分页 Request
type FilePageReq struct {
	core.PageParam
	Path       string   `form:"path"`
	Type       string   `form:"type"`
	CreateTime []string `form:"createTime[]"`
}

// FileUploadReq 上传文件 Request (无需 JSON binding，直接从 Form 获取)
type FileUploadReq struct {
	Path string `form:"path"` // 自定义上传路径/文件名
}
