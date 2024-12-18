package frontend

import (
	"FastGo/internal/handler"
	"FastGo/internal/model"
	"FastGo/internal/router"
	"FastGo/pkg/response"
	"FastGo/pkg/uid"
	"FastGo/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
)

type CollectionHandler struct {
	*handler.CommonHandler
}

func NewCollectionHandler() *CollectionHandler {
	return &CollectionHandler{
		CommonHandler: handler.NewCommonHandler(),
	}
}

func (h *CollectionHandler) RegisterRoutes(routerRegistry *router.RouteRegistry) {
	routerRegistry.Register("POST", "collections", "/create", h.Create, 1, "创建收藏夹")
	routerRegistry.Register("POST", "collections", "/delete", h.Delete, 1, "删除收藏夹")
	routerRegistry.Register("POST", "collections", "/rename", h.Rename, 1, "重命名收藏夹")
	routerRegistry.Register("GET", "collections", "/list", h.GetList, 1, "获取收藏夹列表")
}

// create 创建收藏夹
func (h *CollectionHandler) Create(c *gin.Context) {

	var req struct {
		Name        string `json:"name"`
		WorkspaceID string `json:"workspace_id" binding:"required,workspace_id"`
	}

	result := response.NewResult(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("create collection failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	if req.Name == "" {
		req.Name = "New Collection"
	}

	// 默认创建一个名为“New collection”的收藏夹
	collection := model.Collections{
		Name:         req.Name,
		CollectionID: uid.NewUUID(),
		WorkspaceID:  cast.ToUint64(req.WorkspaceID),
	}

	err := h.DB.Create(&collection).Error
	if err != nil {
		result.FailWithMsg(response.Success, "create collection failed")
		return
	}

	result.Success(nil)
}

// rename 重命名收藏夹
func (h *CollectionHandler) Rename(c *gin.Context) {
	var req struct {
		ID   string `json:"id" binding:"required,id"`
		Name string `json:"name" binding:"required,name"`
	}
	result := response.NewResult(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("rename collection failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	collection := model.Collections{
		ID:   cast.ToUint64(req.ID),
		Name: req.Name,
	}

	err := h.DB.Model(&collection).Where("id = ?", req.ID).Update("name", req.Name).Error
	if err != nil {
		h.Logger.Error("rename collection failed due to database error", zap.Error(err))
		result.FailWithError(response.InvalidParams, "rename collection failed")
		return
	}

	result.Success(map[string]interface{}{
		"id":   req.ID,
		"name": req.Name,
	})
}

// delete 删除收藏夹
func (h *CollectionHandler) Delete(c *gin.Context) {
	var req struct {
		ID string `json:"id" binding:"required,id"`
	}

	result := response.NewResult(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("delete collection failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	// 查询收藏夹是否存在
	collection := model.Collections{ID: cast.ToUint64(req.ID)}
	err := h.DB.First(&collection).Error
	if err != nil {
		h.Logger.Error("collection not found", zap.Error(err))
		result.FailWithMsg(response.Success, "collection not found")
		return
	}

	//记录删除个数
	err = h.DB.Delete(&model.Collections{ID: cast.ToUint64(req.ID)}).Error
	if err != nil {
		h.Logger.Error("delete collection failed due to database error", zap.Error(err))
		result.FailWithMsg(response.Success, "delete collection failed")
		return
	}

	result.Success(nil)
}

// getList 获取收藏夹列表
func (h *CollectionHandler) GetList(c *gin.Context) {

	result := response.NewResult(c)
	workspaceID := c.Query("workspace_id")

	collections := []model.Collections{}
	err := h.DB.Where("workspace_id = ?", workspaceID).Find(&collections).Error
	if err != nil {
		h.Logger.Error("get collection list failed due to database error", zap.Error(err))
		result.FailWithMsg(response.Success, "get collection list failed")
		return
	}

	result.Success(map[string]interface{}{
		"list": collections,
	})
}
