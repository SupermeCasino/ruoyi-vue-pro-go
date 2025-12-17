package handler

import (
	"io/ioutil"
	"strconv"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
	"github.com/wxlbd/ruoyi-mall-go/pkg/response"

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
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	id, err := h.svc.CreateFileConfig(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(id))
}

func (h *FileConfigHandler) UpdateFileConfig(c *gin.Context) {
	var req req.FileConfigSaveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	if err := h.svc.UpdateFileConfig(c, &req); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
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
		c.JSON(400, response.Error(400, "id is required"))
		return
	}

	if err := h.svc.UpdateFileConfigMaster(c, req.ID); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

func (h *FileConfigHandler) DeleteFileConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	if err := h.svc.DeleteFileConfig(c, id); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

func (h *FileConfigHandler) GetFileConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	res, err := h.svc.GetFileConfig(c, id)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

func (h *FileConfigHandler) GetFileConfigPage(c *gin.Context) {
	var req req.FileConfigPageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	res, err := h.svc.GetFileConfigPage(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

func (h *FileConfigHandler) TestFileConfig(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	url, err := h.svc.TestFileConfig(c, id)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(url))
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
		c.JSON(400, response.Error(400, "文件不能为空"))
		return
	}
	path := c.PostForm("path")

	f, err := file.Open()
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	defer f.Close()

	content, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}

	url, err := h.svc.CreateFile(c, file.Filename, path, content)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(url))
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if id == 0 {
		c.JSON(400, response.Error(400, "id is required"))
		return
	}
	if err := h.svc.DeleteFile(c, id); err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(true))
}

func (h *FileHandler) GetFilePage(c *gin.Context) {
	var req req.FilePageReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	res, err := h.svc.GetFilePage(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

func (h *FileHandler) GetFilePresignedUrl(c *gin.Context) {
	path := c.Query("path")
	if path == "" {
		c.JSON(400, response.Error(400, "path is required"))
		return
	}
	res, err := h.svc.GetFilePresignedUrl(c, path)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(res))
}

func (h *FileHandler) CreateFile(c *gin.Context) {
	var req req.FileCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}
	id, err := h.svc.CreateFileCallback(c, &req)
	if err != nil {
		c.JSON(500, response.Error(500, err.Error()))
		return
	}
	c.JSON(200, response.Success(id))
}

func (h *FileHandler) GetFileContent(c *gin.Context) {
	configIdStr := c.Param("configId")
	configId, _ := strconv.ParseInt(configIdStr, 10, 64)
	if configId == 0 {
		c.JSON(400, response.Error(400, "configId is required"))
		return
	}
	// Warning: This implementation might need adjustment depending on how "get/**" wildcard is handled in router
	// For now assuming standard path param
	path := c.Param("path")

	content, err := h.svc.GetFileContent(c, configId, path)
	if err != nil {
		c.JSON(404, response.Error(404, "File not found"))
		return
	}
	c.Data(200, "application/octet-stream", content)
}
