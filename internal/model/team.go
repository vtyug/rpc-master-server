package model

import (
    "time"
    "gorm.io/gorm"
)

// Team 团队表
type Team struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Name        string         `gorm:"type:varchar(64);not null" json:"name"`
    Description string         `gorm:"type:varchar(255)" json:"description"`
    OwnerID     uint          `gorm:"not null" json:"owner_id"`
    Logo        string         `gorm:"type:varchar(255)" json:"logo"`
    Status      uint8          `gorm:"default:1" json:"status"`                    // 1: 正常 0: 禁用
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Team) TableName() string {
    return "teams"
} 