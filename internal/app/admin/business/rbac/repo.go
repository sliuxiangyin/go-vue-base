package rbac

import (
	"databaseAi/internal/app/admin/models"
	"databaseAi/internal/infra/database"
	"errors"

	"gorm.io/gorm"
)

type Repo struct {
	db *database.DB
}

func NewRepo(db *database.DB) *Repo {
	return &Repo{db: db}
}

// ========== 角色相关 ==========

// CreateRole 创建角色
func (r *Repo) CreateRole(role *models.Role) error {
	return r.db.Create(role).Error
}

// UpdateRole 更新角色
func (r *Repo) UpdateRole(role *models.Role) error {
	if role.ID == 0 {
		return errors.New("role id is required")
	}
	return r.db.Save(role).Error
}

// DeleteRole 删除角色
func (r *Repo) DeleteRole(id uint) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// GetRoleByID 根据ID查找角色
func (r *Repo) GetRoleByID(id uint) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName 根据名称查找角色
func (r *Repo) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles 获取角色列表
func (r *Repo) ListRoles(page, pageSize int, name string, status *int8) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	query := r.db.Model(&models.Role{})

	// 条件过滤
	if name != "" {
		query = query.Where("name LIKE ? OR display_name LIKE ?", "%"+name+"%", "%"+name+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Permissions").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&roles).Error

	return roles, total, err
}

// ========== 权限相关 ==========

// CreatePermission 创建权限
func (r *Repo) CreatePermission(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

// UpdatePermission 更新权限
func (r *Repo) UpdatePermission(permission *models.Permission) error {
	if permission.ID == 0 {
		return errors.New("permission id is required")
	}
	return r.db.Save(permission).Error
}

// DeletePermission 删除权限
func (r *Repo) DeletePermission(id uint) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

// GetPermissionByID 根据ID查找权限
func (r *Repo) GetPermissionByID(id uint) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// GetPermissionByName 根据名称查找权限
func (r *Repo) GetPermissionByName(name string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		return nil, err
	}
	return &permission, nil
}

// ListPermissions 获取权限列表
func (r *Repo) ListPermissions(page, pageSize int, category string, permType *models.PermissionType, status *int8) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	query := r.db.Model(&models.Permission{})

	// 条件过滤
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if permType != nil {
		query = query.Where("type = ?", *permType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Offset(offset).
		Limit(pageSize).
		Order("category ASC, sort ASC, created_at DESC").
		Find(&permissions).Error

	return permissions, total, err
}

// ListAllPermissions 获取所有权限（不分页）
func (r *Repo) ListAllPermissions() ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("status = 1").Order("category ASC, sort ASC, created_at DESC").Find(&permissions).Error
	return permissions, err
}

// GetUserPermissionsByType 根据类型获取用户权限
func (r *Repo) GetUserPermissionsByType(userID uint, permType models.PermissionType) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.status = 1 AND permissions.type = ?", userID, permType).
		Order("permissions.sort ASC, permissions.created_at DESC").
		Find(&permissions).Error
	return permissions, err
}

// ========== 角色权限关联 ==========

// AssignPermissionsToRole 为角色分配权限
func (r *Repo) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	// 使用事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除原有权限
		if err := tx.Where("role_id = ?", roleID).Delete(&models.RolePermission{}).Error; err != nil {
			return err
		}

		// 批量插入新权限
		if len(permissionIDs) > 0 {
			rolePermissions := make([]models.RolePermission, 0, len(permissionIDs))
			for _, permissionID := range permissionIDs {
				rolePermissions = append(rolePermissions, models.RolePermission{
					RoleID:       roleID,
					PermissionID: permissionID,
				})
			}
			if err := tx.Create(&rolePermissions).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetRolePermissions 获取角色的所有权限
func (r *Repo) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ? AND permissions.status = 1", roleID).
		Find(&permissions).Error
	return permissions, err
}

// ========== 用户角色关联 ==========

// AssignRolesToUser 为用户分配角色
func (r *Repo) AssignRolesToUser(userID uint, roleIDs []uint) error {
	// 使用事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除原有角色
		if err := tx.Where("user_id = ?", userID).Delete(&models.UserRole{}).Error; err != nil {
			return err
		}

		// 批量插入新角色
		if len(roleIDs) > 0 {
			userRoles := make([]models.UserRole, 0, len(roleIDs))
			for _, roleID := range roleIDs {
				userRoles = append(userRoles, models.UserRole{
					UserID: userID,
					RoleID: roleID,
				})
			}
			if err := tx.Create(&userRoles).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetUserRoles 获取用户的所有角色
func (r *Repo) GetUserRoles(userID uint) ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ? AND roles.status = 1", userID).
		Find(&roles).Error
	return roles, err
}

// GetUserPermissions 获取用户的所有权限（通过角色）
func (r *Repo) GetUserPermissions(userID uint) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Distinct("permissions.*").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.status = 1", userID).
		Find(&permissions).Error
	return permissions, err
}

// CheckUserPermission 检查用户是否有某个权限
func (r *Repo) CheckUserPermission(userID uint, permissionName string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Permission{}).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ? AND permissions.name = ? AND permissions.status = 1", userID, permissionName).
		Count(&count).Error
	return count > 0, err
}

// ========== 管理员用户相关 ==========

// CreateAdminUser 创建管理员
func (r *Repo) CreateAdminUser(user *models.AdminUser) error {
	return r.db.Create(user).Error
}

// UpdateAdminUser 更新管理员
func (r *Repo) UpdateAdminUser(user *models.AdminUser) error {
	if user.ID == 0 {
		return errors.New("user id is required")
	}
	return r.db.Save(user).Error
}

// DeleteAdminUser 删除管理员
func (r *Repo) DeleteAdminUser(id uint) error {
	return r.db.Delete(&models.AdminUser{}, id).Error
}

// GetAdminUserByID 根据ID查找管理员
func (r *Repo) GetAdminUserByID(id uint) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAdminUserByUsername 根据用户名查找管理员
func (r *Repo) GetAdminUserByUsername(username string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAdminUserByEmail 根据邮箱查找管理员
func (r *Repo) GetAdminUserByEmail(email string) (*models.AdminUser, error) {
	var user models.AdminUser
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListAdminUsers 获取管理员列表
func (r *Repo) ListAdminUsers(page, pageSize int, keyword string, status *int8) ([]models.AdminUser, int64, error) {
	var users []models.AdminUser
	var total int64

	query := r.db.Model(&models.AdminUser{})

	// 条件过滤
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Roles").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}
