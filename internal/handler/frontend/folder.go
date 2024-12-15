package frontend

import (
	"FastGo/internal/handler"
	"FastGo/internal/model"
	"FastGo/internal/router"
	"FastGo/pkg/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type FolderHandler struct {
	*handler.CommonHandler
}

func NewFolderHandler() *FolderHandler {
	return &FolderHandler{
		CommonHandler: handler.NewCommonHandler(),
	}
}

func (h *FolderHandler) RegisterRoutes(routerRegistry *router.RouteRegistry) {
	routerRegistry.Register("POST", "folder", "/create", h.Create, 1, "创建文件夹")
	routerRegistry.Register("POST", "folder", "/delete", h.Delete, 1, "删除文件夹")
	routerRegistry.Register("POST", "folder", "/rename", h.Rename, 1, "重命名文件夹")
}

// Create 创建文件夹的处理函数
func (h *FolderHandler) Create(c *gin.Context) {
	var req struct {
		CollectionID uint64 `json:"collection_id" binding:"required,id"`
		ParentID     uint64 `json:"parent_id" binding:"omitempty,id"`
		Name         string `json:"name"`
	}
	result := response.NewResult(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("create folder failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, err.Error())
		return
	}

	folder := model.Folder{
		CollectionID: req.CollectionID,
		Name:         req.Name,
	}

	if err := h.DB.Create(&folder).Error; err != nil {
		h.Logger.Error("create folder failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "创建文件夹失败")
		return
	}

	result.Success(map[string]interface{}{
		"id": folder.ID,
	})
}

// Delete 删除文件夹的处理函数
func (h *FolderHandler) Delete(c *gin.Context) {
	var req struct {
		ID uint64 `json:"id" binding:"required"`
	}
	result := response.NewResult(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("delete folder failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, err.Error())
		return
	}

	if err := h.DB.Delete(&model.Folder{}, req.ID).Error; err != nil {
		h.Logger.Error("delete folder failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "删除文件夹失败")
		return
	}

	result.Success(nil)
}

// Rename 重命名文件夹的处理函数
func (h *FolderHandler) Rename(c *gin.Context) {
	var req struct {
		ID   uint64 `json:"id" binding:"required"`
		Name string `json:"name" binding:"required"`
	}
	result := response.NewResult(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("rename folder failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, err.Error())
		return
	}

	if err := h.DB.Model(&model.Folder{}).Where("id = ?", req.ID).Update("name", req.Name).Error; err != nil {
		h.Logger.Error("rename folder failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "重命名文件夹失败")
		return
	}

	result.Success(map[string]interface{}{
		"id":   req.ID,
		"name": req.Name,
	})
}
