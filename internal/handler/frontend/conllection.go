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
}

// create 创建收藏夹
func (h *CollectionHandler) Create(c *gin.Context) {
	// 默认创建一个名为“New collection”的收藏夹
	collection := model.Collections{
		Name: "New collection",
	}

	result := response.NewResult(c)

	err := h.DB.Create(&collection).Error
	if err != nil {
		result.FailWithMsg(response.SUCCESS, "create collection failed")
		return
	}

	result.Success(map[string]interface{}{
		"id": collection.ID,
	})
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
		result.FailWithError(response.SUCCESS, "rename collection failed")
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
		result.FailWithMsg(response.SUCCESS, "collection not found")
		return
	}

	//记录删除个数
	err = h.DB.Delete(&model.Collections{ID: cast.ToUint64(req.ID)}).Error
	if err != nil {
		h.Logger.Error("delete collection failed due to database error", zap.Error(err))
		result.FailWithMsg(response.SUCCESS, "delete collection failed")
		return
	}

	result.Success(nil)
}
