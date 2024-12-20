package model

import "time"

// 定义请求类型的枚举
type RequestType string

const (
	HTTP      RequestType = "HTTP"
	WebSocket RequestType = "WebSocket"
	GRPC      RequestType = "gRPC"
)

type RequestMethod string

const (
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	PUT    RequestMethod = "PUT"
	DELETE RequestMethod = "DELETE"
)

type Request struct {
	ID           uint64        `gorm:"primaryKey;autoIncrement"`
	Name         string        `gorm:"type:varchar(128);not null"`
	CollectionID string        `gorm:"type:varchar(128);not null;index"`
	FolderID     string        `gorm:"type:varchar(128);not null;index"`
	RequestID    string        `gorm:"type:varchar(128);not null;index"`
	Method       RequestMethod `gorm:"type:varchar(64);not null"`
	Path         string        `gorm:"type:varchar(128);not null"`
	Type         RequestType   `gorm:"type:varchar(64);not null"`
	Headers      string        `gorm:"type:text"`
	Body         string        `gorm:"type:text"`
	QueryParams  string        `gorm:"type:text"`
	Status       string        `gorm:"type:varchar(64)"`
	Response     string        `gorm:"type:text"`
	Timeout      int           `gorm:"type:int"`
	RetryCount   int           `gorm:"type:int"`
	Priority     int           `gorm:"type:int"`
	Description  string        `gorm:"type:text"`
	CreatedAt    time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time     `gorm:"type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (Request) TableName() string {
	return "request"
}
