package frontend

import "FastGo/internal/router"

func RegisterFrontendRoutes(registry *router.RouteRegistry) {
	userHandler := NewUserHandler()
	userHandler.RegisterRoutes(registry)
}
