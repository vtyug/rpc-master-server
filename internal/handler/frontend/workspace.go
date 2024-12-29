package frontend

import (
	"FastGo/internal/handler"
	"FastGo/internal/model"
	"FastGo/internal/router"
	"FastGo/pkg/response"
	"FastGo/pkg/validator"
	"errors"

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
	routerRegistry.Register("POST", "workspaces", "/create", h.Create, 2, "创建工作区")
	routerRegistry.Register("DELETE", "workspaces", "/delete", h.Delete, 1, "删除工作区")
	routerRegistry.Register("PUT", "workspaces", "/edit", h.Edit, 1, "重命名工作区")
	routerRegistry.Register("GET", "workspaces", "/list", h.List, 2, "获取工作区列表")
}

// Create 创建工作空间
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

	userID, _ := c.Get("user_id")

	workspace := model.Workspace{
		Name:    req.Name,
		OwnerID: cast.ToUint64(userID),
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

func (h *WorkspaceHandler) List(c *gin.Context) {

	result := response.NewResult(c)
	workspaces := []model.Workspace{}
	if err := h.DB.Where("owner_id = ?", 1).Find(&workspaces).Error; err != nil {
		h.Logger.Error("get workspace list failed", zap.Error(err))
		result.FailWithMsg(response.Success, "get workspace list failed")
		return
	}

	result.Success(map[string]interface{}{
		"list": workspaces,
	})
}
func (h *WorkspaceHandler) Edit(c *gin.Context) {
	var req struct {
		ID          string `json:"id" binding:"required,id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		OwnerID     string `json:"owner_id"`
	}

	result := response.NewResult(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("edit workspace failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	// 检查工作空间是否存在
	workspace := model.Workspace{}
	if err := h.DB.Where("id = ?", req.ID).First(&workspace).Error; err != nil {
		h.Logger.Error("workspace not found", zap.Error(err))
		result.FailWithMsg(response.Success, "workspace not found")
		return
	}

	// 检查工作空间是否存在并更新
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.OwnerID != "" {
		updates["owner_id"] = req.OwnerID
	}

	if len(updates) > 0 {
		if err := h.DB.Model(&model.Workspace{}).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
			h.Logger.Error("edit workspace failed", zap.Error(err))
			result.FailWithMsg(response.Success, "edit workspace failed")
			return
		}
	} else {
		h.Logger.Warn("no updates provided for workspace", zap.String("id", req.ID))
		result.FailWithMsg(response.InvalidParams, "no updates provided")
		return
	}

	result.Success(nil)
}

func (h *WorkspaceHandler) Delete(c *gin.Context) {
	id := c.Query("id")

	result := response.NewResult(c)

	if id == "" {
		h.Logger.Error("delete workspace failed due to invalid parameters", zap.Error(errors.New("id is required")))
		result.FailWithError(response.InvalidParams, "id is required")
		return
	}

	var count int64
	if err := h.DB.Model(&model.Workspace{}).Where("id = ?", id).Count(&count).Error; err != nil {
		h.Logger.Error("delete workspace failed", zap.Error(err))
		result.FailWithMsg(response.Success, "delete workspace failed")
		return
	}

	if count == 0 {
		result.FailWithMsg(response.Success, "workspace not found")
		return
	}

	// 执行删除操作
	if err := h.DB.Where("id = ?", id).Delete(&model.Workspace{}).Error; err != nil {
		h.Logger.Error("delete workspace failed", zap.Error(err))
		result.FailWithMsg(response.Success, "delete workspace failed")
		return
	}

	result.Success(nil)
}
