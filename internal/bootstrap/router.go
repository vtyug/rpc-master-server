package bootstrap

import (
	"FastGo/internal/global"
	"FastGo/internal/handler/frontend"
	"FastGo/internal/router"
)

func SetupRoutes() {
	engine := global.Engine
	routerRegistry := router.NewRouteRegistry()

	frontend.RegisterFrontendRoutes(routerRegistry)
	routerRegistry.SetupRoutes(engine)
}
