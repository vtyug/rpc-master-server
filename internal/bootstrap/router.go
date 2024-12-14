package bootstrap

import (
	"FastGo/internal/global"
	"FastGo/internal/handler/frontend"
	"FastGo/internal/router"
)

func SetupRoutes() {
	engine := global.Engine
	registry := router.NewRouteRegistry()

	frontend.RegisterFrontendRoutes(registry)
	registry.SetupRoutes(engine)
}
