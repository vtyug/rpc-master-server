package model

import "time"

type Folder struct {
	ID           uint64    `gorm:"primarykey;autoIncrement" json:"id"`                                                     // 文件夹ID
	CollectionID uint64    `gorm:"not null;index" json:"collection_id"`                                                    // 关联到集合
	Name         string    `gorm:"type:varchar(128);not null" json:"name"`                                                 // 文件夹名称
	FolderID     string    `gorm:"type:varchar(128);not null;index" json:"folder_id"`                                      // 文件夹ID
	CreatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`                             // 创建时间
	UpdatedAt    time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"` // 更新时间
}

func (Folder) TableName() string {
	return "folders"
}

type FolderClosure struct {
	Ancestor   string `gorm:"not null;index"` // 祖先文件夹的 ID
	Descendant string `gorm:"not null;index"` // 后代文件夹的 ID
	Depth      int    `gorm:"not null"`       // 祖先与后代之间的距离
}

func (FolderClosure) TableName() string {
	return "folder_closures"
}
