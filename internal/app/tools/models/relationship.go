package models

import "time"

// Relationship 关系模型
type Relationship struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	AccountID   uint      `gorm:"not null" json:"account_id"`
	Title       string    `gorm:"not null" json:"title"`
	Tables      string    `gorm:"type:text" json:"tables"` // JSON格式存储所有表信息
	Ddls        string    `gorm:"type:text" json:"ddls"`   // JSON格式存储所有表结构
	Result      string    `gorm:"type:text" json:"result"` // 结果存储
	CreatedTime time.Time `gorm:"autoCreateTime" json:"created_time"`
	UpdatedTime time.Time `gorm:"autoUpdateTime" json:"updated_time"`
}
