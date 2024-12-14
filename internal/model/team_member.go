package model

import (
	"time"

	"gorm.io/gorm"
)

// TeamMember 团队成员表
type TeamMember struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	TeamID    uint           `gorm:"not null;index" json:"team_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Role      uint8          `gorm:"default:1" json:"role"`   // 1: 普通成员 2: 管理员
	Status    uint8          `gorm:"default:1" json:"status"` // 1: 正常 0: 禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (TeamMember) TableName() string {
	return "team_members"
}
