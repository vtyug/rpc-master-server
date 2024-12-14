package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Username    string         `gorm:"type:varchar(32);not null;unique" json:"username"`
	Password    string         `gorm:"type:varchar(128);not null" json:"-"`
	Email       *string        `gorm:"type:varchar(128);unique" json:"email,omitempty"`
	PhoneNumber *string        `gorm:"type:varchar(15);unique" json:"phone_number,omitempty"`
	GitHubID    *string        `gorm:"type:varchar(64);unique" json:"github_id,omitempty"`
	GiteeID     *string        `gorm:"type:varchar(64);unique" json:"gitee_id,omitempty"`
	WeChatID    *string        `gorm:"type:varchar(64);unique" json:"wechat_id,omitempty"`
	QQID        *string        `gorm:"type:varchar(64);unique" json:"qq_id,omitempty"`
	Nickname    string         `gorm:"type:varchar(32)" json:"nickname"`
	Avatar      string         `gorm:"type:varchar(255)" json:"avatar"`
	Role        uint8          `gorm:"default:1" json:"role"`
	LastLoginAt *time.Time     `gorm:"type:datetime" json:"last_login_at"`
	LastLoginIP string         `gorm:"type:varchar(64)" json:"last_login_ip"`
	Status      uint8          `gorm:"default:1" json:"status"`
	ApiKey      string         `gorm:"type:varchar(64);unique" json:"api_key,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}
