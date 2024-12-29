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

type RequestHandler struct {
	*handler.CommonHandler
}

func NewRequestHandler() *RequestHandler {
	return &RequestHandler{
		CommonHandler: handler.NewCommonHandler(),
	}
}

func (h *RequestHandler) RegisterRoutes(routerRegistry *router.RouteRegistry) {
	routerRegistry.Register("POST", "request", "/create", h.Create, 2, "创建请求")
}

func (h *RequestHandler) Create(c *gin.Context) {
	var req struct {
		Name         string `json:"name"`
		CollectionID string `json:"collection_id" binding:"required,uuid"`
		FolderID     string `json:"folder_id"`
		Type         string `json:"type" binding:"required,oneof=HTTP WebSocket GRPC"`
		Method       string `json:"method" binding:"required,oneof=GET POST PUT DELETE"`
	}
	result := response.NewResult(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("请求参数错误", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	if req.Name == "" {
		req.Name = "New Request"
	}

	request := model.Request{
		Name:         req.Name,
		CollectionID: req.CollectionID,
		FolderID:     req.FolderID,
		RequestID:    uid.NewUUID(),
		Type:         model.RequestType(req.Type),
		Method:       model.RequestMethod(req.Method),
	}

	if err := h.DB.Model(&model.Request{}).Create(&request).Error; err != nil {
		h.Logger.Error("创建请求失败", zap.Error(err))
		result.FailWithError(response.Success, "创建请求失败")
		return
	}

	result.Success(nil)
}
