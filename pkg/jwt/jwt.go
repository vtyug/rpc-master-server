package jwt

import (
	"FastGo/pkg/response"
	"errors"
	"time"

	jwtgo "github.com/golang-jwt/jwt/v4"
)

// 自定义错误
var (
	ErrTokenExpired     = errors.New(response.GetMessage(response.TokenExpired))
	ErrTokenNotValidYet = errors.New(response.GetMessage(response.TokenNotValidYet))
	ErrInvalidToken     = errors.New(response.GetMessage(response.InvalidToken))
)

// CustomClaims 自定义声明
type CustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     uint8  `json:"role"`
	jwtgo.StandardClaims
}

// JWT 结构体
type JWT struct {
	SigningKey []byte
}

// New 创建 JWT 实例
func New(signingKey string) *JWT {
	return &JWT{
		SigningKey: []byte(signingKey),
	}
}

// GenerateToken 生成 token
func (j *JWT) GenerateToken(userID uint, username string, role uint8) (string, error) {
	claims := CustomClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		StandardClaims: jwtgo.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwtgo.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwtgo.Token) (interface{}, error) {
		return j.SigningKey, nil
	})

	if err != nil {
		if ve, ok := err.(*jwtgo.ValidationError); ok {
			if ve.Errors&jwtgo.ValidationErrorMalformed != 0 {
				return nil, ErrInvalidToken
			} else if ve.Errors&jwtgo.ValidationErrorExpired != 0 {
				return nil, ErrTokenExpired
			} else if ve.Errors&jwtgo.ValidationErrorNotValidYet != 0 {
				return nil, ErrTokenNotValidYet
			}
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// RefreshToken 刷新 token
func (j *JWT) RefreshToken(tokenString string) (string, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	claims.StandardClaims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()

	token := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// GetUserInfo 从 token 中获取用户信息
type UserInfo struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     uint8  `json:"role"`
}

func (j *JWT) GetUserInfo(tokenString string) (*UserInfo, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}
	return &UserInfo{
		UserID:   claims.UserID,
		Username: claims.Username,
		Role:     claims.Role,
	}, nil
}

// CheckRole 检查用户角色是否满足要求
func (j *JWT) CheckRole(tokenString string, requiredRole uint8) (bool, error) {
	claims, err := j.ParseToken(tokenString)
	if err != nil {
		return false, err
	}
	return claims.Role >= requiredRole, nil
}
