package frontend

import (
	"FastGo/internal/handler"
	"FastGo/internal/model"
	"FastGo/internal/router"
	"FastGo/pkg/response"
	"FastGo/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
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
		CollectionID string `json:"collection_id" binding:"required,collection_id"`
		Name         string `json:"name"`
		ParentID     string `json:"parent_id" binding:"required,parent_id"`
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
	if err := h.DB.First(&collection, req.CollectionID).Error; err != nil {
		h.Logger.Error("get collection failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "获取集合失败")
		return
	}

	// 创建新文件夹
	folder := model.Folder{
		CollectionID: collection.ID,
		Name:         req.Name,
	}

	if err := h.DB.Create(&folder).Error; err != nil {
		h.Logger.Error("create folder failed due to database error", zap.Error(err))
		result.FailWithMsg(response.ServerError, "创建文件夹失败")
		return
	}

	// 更新闭包表
	// 插入新文件夹与其自身的关系
	closure := model.FolderClosure{
		Ancestor:   folder.ID,
		Descendant: folder.ID,
		Depth:      0,
	}
	if err := h.DB.Create(&closure).Error; err != nil {
		h.Logger.Error("create folder closure failed", zap.Error(err))
		result.FailWithMsg(response.ServerError, "更新文件夹关系失败")
		return
	}

	// 如果提供了 ParentID，插入父子关系
	if cast.ToUint64(req.ParentID) != 0 {
		// 获取父节点的所有祖先关系
		var parentClosures []model.FolderClosure
		if err := h.DB.Where("descendant = ?", cast.ToUint64(req.ParentID)).Find(&parentClosures).Error; err != nil {
			h.Logger.Error("get parent folder closures failed", zap.Error(err))
			result.FailWithMsg(response.ServerError, "获取父文件夹关系失败")
			return
		}

		// 为新节点插入所有祖先关系
		for _, pc := range parentClosures {
			newClosure := model.FolderClosure{
				Ancestor:   pc.Ancestor,
				Descendant: folder.ID,
				Depth:      pc.Depth + 1,
			}
			if err := h.DB.Create(&newClosure).Error; err != nil {
				h.Logger.Error("create ancestor folder closure failed", zap.Error(err))
				result.FailWithMsg(response.ServerError, "更新祖先文件夹关系失败")
				return
			}
		}
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
