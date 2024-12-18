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

type WorkspaceHandler struct {
	*handler.CommonHandler
}

func NewWorkspaceHandler() *WorkspaceHandler {
	return &WorkspaceHandler{
		CommonHandler: handler.NewCommonHandler(),
	}
}

func (h *WorkspaceHandler) RegisterRoutes(routerRegistry *router.RouteRegistry) {
	routerRegistry.Register("POST", "workspaces", "/create", h.Create, 1, "创建工作区")
	// routerRegistry.Register("POST", "workspaces", "/delete", h.Delete, 1, "删除工作区")
	// routerRegistry.Register("POST", "workspaces", "/rename", h.Rename, 1, "重命名工作区")
}

// Create 创建工作区
func (h *WorkspaceHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required,name"`
	}

	result := response.NewResult(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("create workspace failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	workspace := model.Workspace{
		Name: req.Name,
	}

	if err := h.DB.Create(&workspace).Error; err != nil {
		h.Logger.Error("create workspace failed", zap.Error(err))
		result.FailWithMsg(response.Success, "create workspace failed")
		return
	}

	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	resp.ID = cast.ToString(workspace.ID)
	resp.Name = workspace.Name

	result.Success(resp)
}
