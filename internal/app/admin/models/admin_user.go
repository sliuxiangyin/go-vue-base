package models

import (
	"time"

	"gorm.io/gorm"
)

// AdminUser 管理员用户模型
type AdminUser struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Username  string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Nickname  string     `gorm:"type:varchar(50)" json:"nickname"`
	Avatar    string     `gorm:"type:varchar(255)" json:"avatar"`
	Status    int8       `gorm:"type:tinyint;default:1;comment:1-启用 0-禁用" json:"status"`
	LastLogin *time.Time `json:"last_login"`

	// 关联
	Roles []Role `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:RoleID" json:"roles,omitempty"`
}

// TableName 指定表名
func (AdminUser) TableName() string {
	return "admin_users"
}
