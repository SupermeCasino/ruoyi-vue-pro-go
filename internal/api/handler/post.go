package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wxlbd/ruoyi-mall-go/internal/api/req"
	"github.com/wxlbd/ruoyi-mall-go/internal/pkg/core"
	"github.com/wxlbd/ruoyi-mall-go/internal/service"
)

type PostHandler struct {
	svc *service.PostService
}

func NewPostHandler(svc *service.PostService) *PostHandler {
	return &PostHandler{
		svc: svc,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var r req.PostSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	id, err := h.svc.CreatePost(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(id))
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	var r req.PostSaveReq
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	if err := h.svc.UpdatePost(c.Request.Context(), &r); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	if err := h.svc.DeletePost(c.Request.Context(), id); err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(true))
}

func (h *PostHandler) GetPost(c *gin.Context) {
	idStr := c.Query("id")
	id, _ := strconv.ParseInt(idStr, 10, 64)
	item, err := h.svc.GetPost(c.Request.Context(), id)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(item))
}

func (h *PostHandler) GetPostPage(c *gin.Context) {
	var r req.PostPageReq
	if err := c.ShouldBindQuery(&r); err != nil {
		c.JSON(200, core.ErrParam)
		return
	}
	page, err := h.svc.GetPostPage(c.Request.Context(), &r)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(page))
}

func (h *PostHandler) GetSimplePostList(c *gin.Context) {
	list, err := h.svc.GetSimplePostList(c.Request.Context())
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, core.Success(list))
}
