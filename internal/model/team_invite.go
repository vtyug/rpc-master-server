package model

import (
    "time"
)

// TeamInvite 团队邀请表
type TeamInvite struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    TeamID    uint      `gorm:"not null;index" json:"team_id"`
    Email     string    `gorm:"type:varchar(128);not null" json:"email"`
    Code      string    `gorm:"type:varchar(32);not null" json:"code"`
    ExpiredAt time.Time `json:"expired_at"`
    Status    uint8     `gorm:"default:1" json:"status"`                          // 1: 未使用 2: 已使用 3: 已过期
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

func (TeamInvite) TableName() string {
    return "team_invites"
} 