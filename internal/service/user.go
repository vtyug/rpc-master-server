package service

import (
	"FastGo/internal/global"
	"FastGo/internal/model"
	"FastGo/utils"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LoginService struct {
	DB *gorm.DB
}

func NewLoginService() *LoginService {
	return &LoginService{DB: global.GetDB()}
}

// Login 使用多种方式登录
func (s *LoginService) Login(c *gin.Context) (*model.User, error) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password" binding:"required"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		GitHubID    string `json:"github_id"`
		GiteeID     string `json:"gitee_id"`
		WeChatID    string `json:"wechat_id"`
		QQID        string `json:"qq_id"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		return nil, err
	}

	var user model.User
	var err error

	switch {
	case loginData.Username != "":
		err = s.DB.Where("username = ?", loginData.Username).First(&user).Error
	case loginData.Email != "":
		err = s.DB.Where("email = ?", loginData.Email).First(&user).Error
	case loginData.PhoneNumber != "":
		err = s.DB.Where("phone_number = ?", loginData.PhoneNumber).First(&user).Error
	case loginData.GitHubID != "":
		err = s.DB.Where("github_id = ?", loginData.GitHubID).First(&user).Error
	case loginData.GiteeID != "":
		err = s.DB.Where("gitee_id = ?", loginData.GiteeID).First(&user).Error
	case loginData.WeChatID != "":
		err = s.DB.Where("wechat_id = ?", loginData.WeChatID).First(&user).Error
	case loginData.QQID != "":
		err = s.DB.Where("qq_id = ?", loginData.QQID).First(&user).Error
	default:
		return nil, errors.New("必须提供用户名、邮箱或手机号之一")
	}

	if err != nil {
		return nil, errors.New("用户不存在或密码错误")
	}

	if !utils.CheckPasswordHash(loginData.Password, user.Password) {
		return nil, errors.New("用户不存在或密码错误")
	}

	return &user, nil
}
