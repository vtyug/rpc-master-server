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

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
func (h *UserHandler) RegisterRoutes(registry *router.RouteRegistry) {
	registry.Register("POST", "user", "/login", h.Login, 1, "用户登录")
	registry.Register("POST", "user", "/register", h.Register, 1, "用户注册")
	registry.Register("GET", "user", "/profile", h.Profile, 1, "获取用户信息")
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
		result.FailWithError(response.ServerError, "Token 生成失败")
		return
	}

	result.Success(map[string]interface{}{
		"user_id": user.ID,
		"token":   token,
	})
}

func (h *UserHandler) Register(c *gin.Context) {
	var registerData struct {
		Username    string  `json:"username" binding:"required,username"`
		Password    string  `json:"password" binding:"required,password"`
		Email       *string `json:"email,omitempty"`
		PhoneNumber *string `json:"phone_number,omitempty"`
		GitHubID    *string `json:"github_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&registerData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validator.TranslateError(err)})
		return
	}

	// 检查非空字段的唯一性
	if registerData.Email != nil {
		var existingUser model.User
		if err := h.DB.Where("email = ?", *registerData.Email).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被注册"})
			return
		}
	}

	if registerData.GitHubID != nil {
		var existingUser model.User
		if err := h.DB.Where("git_hub_id = ?", *registerData.GitHubID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "GitHub ID 已被注册"})
			return
		}
	}

	hashedPassword, err := utils.HashPassword(registerData.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 生成唯一的 API Key
	apiKey := uuid.New().String()

	user := model.User{
		Username:    registerData.Username,
		Password:    hashedPassword,
		Email:       registerData.Email,
		PhoneNumber: registerData.PhoneNumber,
		GitHubID:    registerData.GitHubID,
		ApiKey:      apiKey,
	}

	if err := h.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}

	result := response.NewResult(c)
	result.Success("注册成功")
}

func (h *UserHandler) Profile(c *gin.Context) {
	// 获取用户信息逻辑
}
