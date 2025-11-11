package models

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string `gorm:"type:varchar(50);uniqueIndex;not null;comment:角色名称" json:"name"`
	DisplayName string `gorm:"type:varchar(100);not null;comment:角色显示名称" json:"display_name"`
	Description string `gorm:"type:varchar(255);comment:角色描述" json:"description"`
	Status      int8   `gorm:"type:tinyint;default:1;comment:1-启用 0-禁用" json:"status"`
	IsSystem    bool   `gorm:"type:tinyint;default:0;comment:是否系统角色(系统角色不可删除)" json:"is_system"`

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []AdminUser  `gorm:"many2many:user_roles;foreignKey:ID;joinForeignKey:RoleID;References:ID;joinReferences:UserID" json:"users,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// PermissionType 权限类型
type PermissionType string

const (
	PermissionTypeBackend  PermissionType = "backend"  // 后端API权限
	PermissionTypeFrontend PermissionType = "frontend" // 前端路由/菜单权限
)

// Permission 权限模型
type Permission struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Name        string         `gorm:"type:varchar(100);uniqueIndex;not null;comment:权限标识(路由名称)" json:"name"`
	DisplayName string         `gorm:"type:varchar(100);not null;comment:权限显示名称" json:"display_name"`
	Description string         `gorm:"type:varchar(255);comment:权限描述" json:"description"`
	Category    string         `gorm:"type:varchar(50);comment:权限分类" json:"category"`
	Type        PermissionType `gorm:"type:varchar(20);default:'backend';comment:权限类型(backend/frontend)" json:"type"`
	Path        string         `gorm:"type:varchar(255);comment:前端路由路径(仅前端权限)" json:"path"`
	Icon        string         `gorm:"type:varchar(50);comment:菜单图标(仅前端权限)" json:"icon"`
	ParentID    *uint          `gorm:"comment:父级权限ID(用于菜单层级)" json:"parent_id"`
	Sort        int            `gorm:"default:0;comment:排序" json:"sort"`
	Status      int8           `gorm:"type:tinyint;default:1;comment:1-启用 0-禁用" json:"status"`

	// 关联
	Roles    []Role       `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
	Children []Permission `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	RoleID       uint      `gorm:"not null;index:idx_role_permission,unique" json:"role_id"`
	PermissionID uint      `gorm:"not null;index:idx_role_permission,unique" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_role,unique" json:"user_id"`
	RoleID    uint      `gorm:"not null;index:idx_user_role,unique" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}
