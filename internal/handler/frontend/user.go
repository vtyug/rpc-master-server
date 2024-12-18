package frontend

import (
	"FastGo/internal/handler"
	"FastGo/internal/model"
	"FastGo/internal/router"
	"FastGo/internal/service"
	"FastGo/pkg/jwt"
	"FastGo/pkg/response"
	"FastGo/pkg/validator"
	"FastGo/utils"

	"github.com/gin-gonic/gin"
)

// UserHandler 处理用户相关的请求
type UserHandler struct {
	*handler.CommonHandler
	LoginService *service.LoginService
}

// NewUserHandler 创建并返回一个 UserHandler 实例
func NewUserHandler() *UserHandler {
	return &UserHandler{
		CommonHandler: handler.NewCommonHandler(),
		LoginService:  service.NewLoginService(),
	}
}

// RegisterRoutes 注册用户模块的路由
func (h *UserHandler) RegisterRoutes(routerRegistry *router.RouteRegistry) {
	routerRegistry.Register("POST", "user", "/login", h.Login, 1, "用户登录")
	routerRegistry.Register("POST", "user", "/register", h.Register, 1, "用户注册")
	routerRegistry.Register("GET", "user", "/profile", h.Profile, 1, "获取用户信息")
}

func (h *UserHandler) Login(c *gin.Context) {
	result := response.NewResult(c)

	user, err := h.LoginService.Login(c)
	if err != nil {
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	// 生成 JWT
	jwtService := jwt.New("yug-fastgo")
	token, err := jwtService.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		// result.FailWithError(response.ServerError, "Token 生成失败")
		return
	}

	result.Success(map[string]interface{}{
		"user_id": user.ID,
		"token":   token,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var registerData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	result := response.NewResult(c)

	if err := c.ShouldBindJSON(&registerData); err != nil {
		result.FailWithError(response.InvalidParams, validator.TranslateError(err))
		return
	}

	

	// 检查用户名是否已存在
	var existingUser model.User
	if err := h.DB.Where("username = ?", registerData.Username).First(&existingUser).Error; err == nil {
		result.FailWithMsg(response.InvalidParams, "username already exists")
		return
	}

	hashedPassword, err := utils.HashPassword(registerData.Password)
	if err != nil {
		result.FailWithMsg(response.ServerError, "password encryption failed")
		return
	}

	user := model.User{
		Username: registerData.Username,
		Password: hashedPassword,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		result.FailWithMsg(response.ServerError, "register failed")
		return
	}

	result.Success(nil)
}

func (h *UserHandler) Profile(c *gin.Context) {
	// 获取用户信息逻辑
}
