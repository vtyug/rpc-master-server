package frontend

import (
	"FastGo/internal/handler"
	"FastGo/internal/model"
	"FastGo/internal/router"
	"FastGo/pkg/response"
	"FastGo/pkg/uid"
	"FastGo/pkg/validator"

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
	routerRegistry.Register("GET", "folder", "/list", h.List, 1, "获取文件夹列表")
}

// Create 创建文件夹的处理函数
func (h *FolderHandler) Create(c *gin.Context) {
	var req struct {
		CollectionID string `json:"collection_id" binding:"required,uuid"`
		Name         string `json:"name"`
		FolderID     string `json:"folder_id"`
	}
	result := response.NewResult(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("create folder failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	if req.Name == "" {
		req.Name = "New Folder"
	}

	// 查找集合
	var collection model.Collections
	if err := h.DB.Where("collection_id = ?", req.CollectionID).First(&collection).Error; err != nil {
		h.Logger.Error("get collection failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "get collection failed")
		return
	}

	// 创建新文件夹
	folder := model.Folder{
		CollectionID: collection.CollectionID,
		Name:         req.Name,
		FolderID:     uid.NewUUID(),
	}

	if err := h.DB.Create(&folder).Error; err != nil {
		h.Logger.Error("create folder failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "create folder failed")
		return
	}

	// 更新闭包表
	closure := model.FolderClosure{
		Ancestor:   folder.FolderID,
		Descendant: folder.FolderID,
		Depth:      0,
	}
	if err := h.DB.Create(&closure).Error; err != nil {
		h.Logger.Error("create folder closure failed", zap.Error(err))
		result.FailWithMsg(response.ServerError, "update folder closure failed")
		return
	}

	// 如果提供了 FolderID，插入父子关系
	if req.FolderID != "" {
		var parentClosures []model.FolderClosure
		if err := h.DB.Where("descendant = ?", req.FolderID).Find(&parentClosures).Error; err != nil {
			h.Logger.Error("get parent folder closures failed", zap.Error(err))
			result.FailWithMsg(response.ServerError, "get parent folder closures failed")
			return
		}

		for _, pc := range parentClosures {
			// 检查是否已经存在相同的祖先-后代关系
			var existingClosure model.FolderClosure
			if err := h.DB.Where("ancestor = ? AND descendant = ?", pc.Ancestor, folder.FolderID).First(&existingClosure).Error; err == nil {
				continue // 如果关系已经存在，跳过插入
			}

			newClosure := model.FolderClosure{
				Ancestor:   pc.Ancestor,
				Descendant: folder.FolderID,
				Depth:      pc.Depth + 1,
			}
			if err := h.DB.Create(&newClosure).Error; err != nil {
				h.Logger.Error("create ancestor folder closure failed", zap.Error(err))
				result.FailWithMsg(response.ServerError, "update ancestor folder closure failed")
				return
			}
		}

		// 插入父子关系
		var parentChildClosure model.FolderClosure
		if err := h.DB.Where("ancestor = ? AND descendant = ?", req.FolderID, folder.FolderID).First(&parentChildClosure).Error; err != nil {
			parentChildClosure = model.FolderClosure{
				Ancestor:   req.FolderID,
				Descendant: folder.FolderID,
				Depth:      1,
			}
			if err := h.DB.Create(&parentChildClosure).Error; err != nil {
				h.Logger.Error("create parent-child folder closure failed", zap.Error(err))
				result.FailWithMsg(response.ServerError, "update parent-child folder closure failed")
				return
			}
		}
	}

	result.Success(map[string]interface{}{
		"folder_id": folder.FolderID,
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
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	if err := h.DB.Delete(&model.Folder{}, req.ID).Error; err != nil {
		h.Logger.Error("delete folder failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "delete folder failed")
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
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	if err := h.DB.Model(&model.Folder{}).Where("id = ?", req.ID).Update("name", req.Name).Error; err != nil {
		h.Logger.Error("rename folder failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "rename folder failed")
		return
	}

	result.Success(map[string]interface{}{
		"id":   req.ID,
		"name": req.Name,
	})
}

// List 获取文件夹列表
func (h *FolderHandler) List(c *gin.Context) {
	result := response.NewResult(c)

	collectionID := c.Query("collection_id")

	folders := []model.Folder{}
	if err := h.DB.Where("collection_id = ?", collectionID).Find(&folders).Error; err != nil {
		h.Logger.Error("get folder list failed", zap.Error(err))
		result.FailWithMsg(response.ServerError, "get folder list failed")
		return
	}
	result.Success(map[string]interface{}{
		"list": folders,
	})
}
