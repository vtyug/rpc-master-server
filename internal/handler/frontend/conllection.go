package frontend

import (
	"FastGo/internal/handler"
	"FastGo/internal/model"
	"FastGo/internal/router"
	"FastGo/pkg/response"
	"FastGo/pkg/uid"
	"FastGo/pkg/validator"
	"errors"

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
	routerRegistry.Register("POST", "collections", "/create", h.Create, 2, "创建集合")
	routerRegistry.Register("DELETE", "collections", "/delete", h.Delete, 2, "删除集合")
	routerRegistry.Register("POST", "collections", "/edit", h.Edit, 2, "编辑集合")
	routerRegistry.Register("GET", "collections", "/list", h.GetList, 2, "获取集合列表")
}

// create 创建收藏夹
func (h *CollectionHandler) Create(c *gin.Context) {

	var req struct {
		WorkspaceID string `json:"workspace_id" binding:"required,workspace_id"`
		Name        string `json:"name" binding:"required,name"`
		Protocol    string `json:"protocol"`
		Description string `json:"description"`
	}

	result := response.NewResult(c)

	if err := c.ShouldBindJSON(&req); err != nil {
		h.Logger.Error("create collection failed due to invalid parameters", zap.Error(err))
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	collection := model.Collections{
		Name:         req.Name,
		CollectionID: uid.NewUUID(),
		WorkspaceID:  cast.ToUint64(req.WorkspaceID),
		Protocol:     model.FromString(req.Protocol),
		Description:  req.Description,
	}

	err := h.DB.Create(&collection).Error
	if err != nil {
		result.FailWithMsg(response.Success, "create collection failed")
		return
	}

	result.Success(nil)
}

// rename 重命名收藏夹
func (h *CollectionHandler) Edit(c *gin.Context) {
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

	id := c.Query("id")

	result := response.NewResult(c)

	if id == "" {
		h.Logger.Error("delete collection failed due to invalid parameters", zap.Error(errors.New("id is required")))
		result.FailWithMsg(response.Success, "delete collection failed")
		return
	}

	// 查询收藏夹是否存在
	collection := model.Collections{ID: cast.ToUint64(id)}
	err := h.DB.First(&collection).Error
	if err != nil {
		h.Logger.Error("collection not found", zap.Error(err))
		result.FailWithMsg(response.Success, "collection not found")
		return
	}

	//记录删除个数
	err = h.DB.Delete(&model.Collections{ID: cast.ToUint64(id)}).Error
	if err != nil {
		h.Logger.Error("delete collection failed due to database error", zap.Error(err))
		result.FailWithMsg(response.Success, "delete collection failed")
		return
	}

	result.Success(nil)
}

// 获取收藏夹列表
func (h *CollectionHandler) GetList(c *gin.Context) {
	result := response.NewResult(c)
	workspaceID := c.Query("workspace_id")

	type resp struct {
		ID           string `json:"id"`
		Name         string `json:"name"`
		Protocol     string `json:"protocol"`
		Owner        string `json:"owner"`
		Description  string `json:"description"`
		WorkspaceID  string `json:"workspace_id"`
		CollectionID string `json:"collection_id"`
		MembersCount int    `json:"members_count"`
		CreatedAt    string `json:"created_at"`
	}

	collections := []model.Collections{}
	err := h.DB.Where("workspace_id = ?", workspaceID).
		Find(&collections).Error
	if err != nil {
		h.Logger.Error("get collection list failed due to database error", zap.Error(err))
		result.FailWithMsg(response.Success, "get collection list failed")
		return
	}

	// 获取所有需要的用户ID
	userIDs := make([]uint64, 0, len(collections))
	for _, collection := range collections {
		userIDs = append(userIDs, collection.OwnerID)
	}

	// 批量查询用户信息
	users := []model.User{}
	h.DB.Where("id IN (?)", userIDs).Find(&users)

	// 创建用户ID到用户昵称的映射
	userMap := make(map[uint64]string)
	for _, user := range users {
		userMap[uint64(user.ID)] = user.Nickname
	}

	var responseList = make([]resp, 0, len(collections))
	for _, collection := range collections {
		resp := resp{
			ID:           cast.ToString(collection.ID),
			Name:         collection.Name,
			Protocol:     model.ReturnString(collection.Protocol),
			Owner:        userMap[collection.OwnerID], // 使用映射获取昵称
			Description:  collection.Description,
			WorkspaceID:  cast.ToString(collection.WorkspaceID),
			CollectionID: collection.CollectionID,
			MembersCount: collection.MembersCount,
			CreatedAt:    collection.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		responseList = append(responseList, resp)
	}

	result.Success(map[string]interface{}{
		"list": responseList,
	})
}

// 定义响应结构体
// type RequestResponse struct {
// 	ID           string `json:"id"`
// 	Name         string `json:"name"`
// 	Type         string `json:"type"`
// 	Method       string `json:"method"`
// 	CollectionID string `json:"collection_id"`
// 	RequestID    string `json:"request_id"`
// 	FolderID     string `json:"folder_id"`
// 	Kind         string `json:"kind"`
// }

// // 定义响应结构体
// type FolderResponse struct {
// 	ID           string        `json:"id"`
// 	Name         string        `json:"name"`
// 	CollectionID string        `json:"collection_id"`
// 	FolderID     string        `json:"folder_id"`
// 	Children     []interface{} `json:"children"`
// 	Kind         string        `json:"kind"`
// }

// type CollectionResponse struct {
// 	ID           string        `json:"id"`
// 	Name         string        `json:"name"`
// 	WorkspaceID  string        `json:"workspace_id"`
// 	CollectionID string        `json:"collection_id"`
// 	Kind         string        `json:"kind"`
// 	Children     []interface{} `json:"children"`
// }

// func (h *CollectionHandler) GetList(c *gin.Context) {
// 	result := response.NewResult(c)
// 	workspaceID := c.Query("workspace_id")

// 	// 获取所有集合
// 	collections := []model.Collections{}
// 	err := h.DB.Where("workspace_id = ?", workspaceID).Find(&collections).Error
// 	if err != nil {
// 		h.Logger.Error("get collection list failed due to database error", zap.Error(err))
// 		result.FailWithMsg(response.Success, "get collection list failed")
// 		return
// 	}

// 	var responseList []CollectionResponse

// 	for _, collection := range collections {

// 		// 获取所有属于当前集合的文件夹
// 		var allFolders []model.Folder
// 		err = h.DB.Where("collection_id = ?", collection.CollectionID).Find(&allFolders).Error
// 		if err != nil {
// 			h.Logger.Error("failed to get folders", zap.Error(err))
// 			result.FailWithMsg(response.Success, "failed to get folders")
// 			return
// 		}

// 		folderIDs := getFolderIDs(allFolders)

// 		// 获取所有闭包关系
// 		var folderClosures []model.FolderClosure
// 		err = h.DB.Where("ancestor IN (?)", folderIDs).Find(&folderClosures).Error
// 		if err != nil {
// 			h.Logger.Error("failed to get folder closures", zap.Error(err))
// 			result.FailWithMsg(response.Success, "failed to get folder closures")
// 			return
// 		}

// 		// 构建文件夹树
// 		folderMap := make(map[string]FolderResponse)
// 		for _, folder := range allFolders {
// 			folderMap[folder.FolderID] = FolderResponse{
// 				ID:           cast.ToString(folder.ID),
// 				Name:         folder.Name,
// 				CollectionID: folder.CollectionID,
// 				FolderID:     folder.FolderID,
// 				Kind:         "folder",
// 				Children:     []interface{}{},
// 			}
// 		}

// 		// 使用闭包表构建树
// 		for _, closure := range folderClosures {
// 			ancestor := closure.Ancestor
// 			descendant := closure.Descendant
// 			if ancestor != descendant {
// 				parentFolder, parentExists := folderMap[ancestor]
// 				childFolder, childExists := folderMap[descendant]
// 				if parentExists && childExists {
// 					parentFolder.Children = append(parentFolder.Children, childFolder)
// 					folderMap[ancestor] = parentFolder
// 				}
// 			}
// 		}

// 		// 获取所有属于当前集合的请求
// 		var allRequests []model.Request
// 		err := h.DB.Where("collection_id = ?", collection.CollectionID).Find(&allRequests).Error
// 		if err != nil {
// 			h.Logger.Error("failed to get requests", zap.Error(err))
// 			result.FailWithMsg(response.Success, "failed to get requests")
// 			return
// 		}

// 		// 构建请求映射
// 		requestMap := make(map[string][]RequestResponse)
// 		for _, request := range allRequests {
// 			reqResponse := RequestResponse{
// 				ID:           cast.ToString(request.ID),
// 				Name:         request.Name,
// 				Type:         string(request.Type),
// 				Method:       string(request.Method),
// 				CollectionID: request.CollectionID,
// 				RequestID:    request.RequestID,
// 				FolderID:     request.FolderID,
// 				Kind:         "request",
// 			}
// 			requestMap[request.FolderID] = append(requestMap[request.FolderID], reqResponse)
// 		}

// 		// 将请求添加到文件夹中
// 		for folderID, folder := range folderMap {
// 			for _, req := range requestMap[folderID] {
// 				folder.Children = append(folder.Children, req)
// 			}
// 			folderMap[folderID] = folder
// 		}

// 		// 找到根文件夹并添加请求
// 		var rootFolders []interface{}
// 		for _, folder := range folderMap {
// 			if folder.FolderID == "" || !isDescendant(folder.FolderID, folderClosures) {
// 				rootFolders = append(rootFolders, folder)
// 			}
// 		}
// 		for _, req := range requestMap[""] {
// 			rootFolders = append(rootFolders, req)
// 		}

// 		// 添加到响应列表
// 		responseList = append(responseList, CollectionResponse{
// 			ID:           cast.ToString(collection.ID),
// 			Name:         collection.Name,
// 			WorkspaceID:  cast.ToString(collection.WorkspaceID),
// 			CollectionID: collection.CollectionID,
// 			Kind:         "collection",
// 			Children:     rootFolders,
// 		})
// 	}

// 	result.Success(map[string]interface{}{
// 		"list": responseList,
// 	})
// }

// // 获取文件夹ID列表
// func getFolderIDs(folders []model.Folder) []string {
// 	ids := make([]string, len(folders))
// 	for i, folder := range folders {
// 		ids[i] = cast.ToString(folder.FolderID)
// 	}
// 	return ids
// }

// // 检查文件夹是否是其他文件夹的后代
// func isDescendant(folderID string, closures []model.FolderClosure) bool {
// 	for _, closure := range closures {
// 		if closure.Descendant == folderID && closure.Ancestor != folderID {
// 			return true
// 		}
// 	}
// 	return false
// }
