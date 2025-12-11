package handler

import (
	"backend-go/internal/api/req"
	"backend-go/internal/pkg/core"
	"backend-go/internal/service"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileConfigHandler struct {
	svc *service.FileConfigService
}

func NewFileConfigHandler(svc *service.FileConfigService) *FileConfigHandler {
	return &FileConfigHandler{svc: svc}
}

func (h *FileConfigHandler) CreateFileConfig(c *gin.Context) {
	var req req.FileConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	id, err := h.svc.CreateFileConfig(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(id))
}

func (h *FileConfigHandler) UpdateFileConfig(c *gin.Context) {
	var req req.FileConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	if err := h.svc.UpdateFileConfig(c, &req); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *FileConfigHandler) UpdateFileConfigMaster(c *gin.Context) {
	var req struct {
		ID int64 `json:"id"`
	}
	// Support both JSON body and Query param
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == 0 {
		idStr := c.Query("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)
		req.ID = id
	}

	if req.ID == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}

	if err := h.svc.UpdateFileConfigMaster(c, req.ID); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *FileConfigHandler) DeleteFileConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	if err := h.svc.DeleteFileConfig(c, id); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *FileConfigHandler) GetFileConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	res, err := h.svc.GetFileConfig(c, id)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}

func (h *FileConfigHandler) GetFileConfigPage(c *gin.Context) {
	var req req.FileConfigPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	res, err := h.svc.GetFileConfigPage(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}

// File Handler

type FileHandler struct {
	svc *service.FileService
}

func NewFileHandler(svc *service.FileService) *FileHandler {
	return &FileHandler{svc: svc}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, core.Error(400, "文件不能为空"))
		return
	}
	path := c.PostForm("path")

	f, err := file.Open()
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}

	url, err := h.svc.CreateFile(c, file.Filename, path, content)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(url))
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, core.Error(400, "id is required"))
		return
	}
	if err := h.svc.DeleteFile(c, id); err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *FileHandler) GetFilePage(c *gin.Context) {
	var req req.FilePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, core.Error(400, err.Error()))
		return
	}
	res, err := h.svc.GetFilePage(c, &req)
	if err != nil {
		c.JSON(500, core.Error(500, err.Error()))
		return
	}
	c.JSON(200, core.Success(res))
}
