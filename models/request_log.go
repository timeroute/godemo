package models

import (
	"time"

	"gorm.io/gorm"
)

// RequestLog 请求日志模型
type RequestLog struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	Method       string         `gorm:"size:10;not null" json:"method"`           // HTTP 方法
	Path         string         `gorm:"size:500;not null" json:"path"`            // 请求路径
	Query        string         `gorm:"size:1000" json:"query"`                   // 查询参数
	Status       int            `gorm:"not null" json:"status"`                   // 响应状态码
	Latency      int64          `json:"latency"`                                  // 响应时间（毫秒）
	ClientIP     string         `gorm:"size:45;index" json:"client_ip"`           // 客户端 IP
	UserAgent    string         `gorm:"size:500" json:"user_agent"`               // User-Agent
	ErrorMessage string         `gorm:"size:1000" json:"error_message,omitempty"` // 错误信息
	RequestBody  string         `gorm:"type:text" json:"request_body,omitempty"`  // 请求体（可选）
	ResponseBody string         `gorm:"type:text" json:"response_body,omitempty"` // 响应体（可选）
	UserID       *uint          `gorm:"index" json:"user_id,omitempty"`           // 用户 ID（如果已认证）
	Username     string         `gorm:"size:100" json:"username,omitempty"`       // 用户名（如果已认证）
	Country      string         `gorm:"size:100;index" json:"country,omitempty"`  // 国家
	City         string         `gorm:"size:100;index" json:"city,omitempty"`     // 城市
	Latitude     float64        `json:"latitude,omitempty"`                       // 纬度
	Longitude    float64        `json:"longitude,omitempty"`                      // 经度
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (RequestLog) TableName() string {
	return "request_logs"
}
