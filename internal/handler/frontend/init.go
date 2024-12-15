package frontend

import "FastGo/internal/router"

func RegisterFrontendRoutes(routerRegistry *router.RouteRegistry) {
	userHandler := NewUserHandler()
	userHandler.RegisterRoutes(routerRegistry)

	collectionHandler := NewCollectionHandler()
	collectionHandler.RegisterRoutes(routerRegistry)

	folderHandler := NewFolderHandler()
	folderHandler.RegisterRoutes(routerRegistry)
}
