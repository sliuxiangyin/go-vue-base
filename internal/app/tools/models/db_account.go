package models

import "time"

// DBAccount 数据库账号模型
type DBAccount struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name" form:"name" binding:"required"`
	Host        string    `gorm:"not null" json:"host" form:"host" binding:"required"`
	Account     string    `gorm:"not null" json:"account" form:"account" binding:"required"`
	Password    string    `gorm:"not null" json:"password" form:"password" binding:"required"`
	Port        int       `gorm:"not null" json:"port" form:"port" binding:"required"`
	Dbname      string    `gorm:"not null" json:"dbname" form:"dbname" binding:"required"`
	Ddl         string    `gorm:"column:ddl;" json:"ddl" form:"port" binding:"required"`
	CreatedTime time.Time `gorm:"autoCreateTime" json:"created_time"`
	UpdatedTime time.Time `gorm:"autoUpdateTime" json:"updated_time"`
}
