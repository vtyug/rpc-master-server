package frontend

import "FastGo/internal/router"

func RegisterFrontendRoutes(routerRegistry *router.RouteRegistry) {
	userHandler := NewUserHandler()
	userHandler.RegisterRoutes(routerRegistry)

	// conllection 收藏夹
	collectionHandler := NewCollectionHandler()
	collectionHandler.RegisterRoutes(routerRegistry)

	// folder 文件夹
	folderHandler := NewFolderHandler()
	folderHandler.RegisterRoutes(routerRegistry)

	// workspace 工作区
	workspaceHandler := NewWorkspaceHandler()
	workspaceHandler.RegisterRoutes(routerRegistry)
}
