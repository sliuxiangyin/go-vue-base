package rbac

import (
	"databaseAi/internal/app/admin/models"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *Repo
}

func NewService(repo *Repo) *Service {
	return &Service{repo: repo}
}

// ========== 角色管理 ==========

// CreateRole 创建角色
func (s *Service) CreateRole(name, displayName, description string, status int8) (*models.Role, error) {
	// 检查角色名是否已存在
	if existingRole, _ := s.repo.GetRoleByName(name); existingRole != nil {
		return nil, errors.New("角色名称已存在")
	}

	role := &models.Role{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Status:      status,
		IsSystem:    false,
	}

	if err := s.repo.CreateRole(role); err != nil {
		return nil, fmt.Errorf("创建角色失败: %w", err)
	}

	return role, nil
}

// UpdateRole 更新角色
func (s *Service) UpdateRole(id uint, displayName, description string, status int8) (*models.Role, error) {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return nil, errors.New("角色不存在")
	}

	role.DisplayName = displayName
	role.Description = description
	role.Status = status

	if err := s.repo.UpdateRole(role); err != nil {
		return nil, fmt.Errorf("更新角色失败: %w", err)
	}

	return role, nil
}

// DeleteRole 删除角色
func (s *Service) DeleteRole(id uint) error {
	role, err := s.repo.GetRoleByID(id)
	if err != nil {
		return errors.New("角色不存在")
	}

	if role.IsSystem {
		return errors.New("系统角色不允许删除")
	}

	return s.repo.DeleteRole(id)
}

// GetRole 获取角色详情
func (s *Service) GetRole(id uint) (*models.Role, error) {
	return s.repo.GetRoleByID(id)
}

// ListRoles 获取角色列表
func (s *Service) ListRoles(page, pageSize int, name string, status *int8) ([]models.Role, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.ListRoles(page, pageSize, name, status)
}

// ========== 权限管理 ==========

// CreatePermission 创建权限
func (s *Service) CreatePermission(name, displayName, description, category string, permType models.PermissionType, path, icon string, parentID *uint, sort int, status int8) (*models.Permission, error) {
	// 检查权限名是否已存在
	if existingPerm, _ := s.repo.GetPermissionByName(name); existingPerm != nil {
		return nil, errors.New("权限标识已存在")
	}

	permission := &models.Permission{
		Name:        name,
		DisplayName: displayName,
		Description: description,
		Category:    category,
		Type:        permType,
		Path:        path,
		Icon:        icon,
		ParentID:    parentID,
		Sort:        sort,
		Status:      status,
	}

	if err := s.repo.CreatePermission(permission); err != nil {
		return nil, fmt.Errorf("创建权限失败: %w", err)
	}

	return permission, nil
}

// UpdatePermission 更新权限
func (s *Service) UpdatePermission(id uint, displayName, description, category string, path, icon string, parentID *uint, sort int, status int8) (*models.Permission, error) {
	permission, err := s.repo.GetPermissionByID(id)
	if err != nil {
		return nil, errors.New("权限不存在")
	}

	permission.DisplayName = displayName
	permission.Description = description
	permission.Category = category
	// permission.Type = permType
	permission.Path = path
	permission.Icon = icon
	permission.ParentID = parentID
	permission.Sort = sort
	permission.Status = status

	if err := s.repo.UpdatePermission(permission); err != nil {
		return nil, fmt.Errorf("更新权限失败: %w", err)
	}

	return permission, nil
}

// DeletePermission 删除权限
func (s *Service) DeletePermission(id uint) error {
	if _, err := s.repo.GetPermissionByID(id); err != nil {
		return errors.New("权限不存在")
	}

	return s.repo.DeletePermission(id)
}

// GetPermission 获取权限详情
func (s *Service) GetPermission(id uint) (*models.Permission, error) {
	return s.repo.GetPermissionByID(id)
}

// ListPermissions 获取权限列表
func (s *Service) ListPermissions(page, pageSize int, category string, permType *models.PermissionType, status *int8) ([]models.Permission, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	permissions, total, err := s.repo.ListPermissions(page, pageSize, category, permType, status)
	if err != nil {
		return nil, 0, err
	}

	// 如果是前端权限，构建树形结构
	if permType != nil && *permType == models.PermissionTypeFrontend {
		return s.buildPermissionTree(permissions), total, nil
	}

	return permissions, total, nil
}

// ListAllPermissions 获取所有权限（不分页）
func (s *Service) ListAllPermissions() ([]models.Permission, error) {
	return s.repo.ListAllPermissions()
}

// ========== 角色权限关联 ==========

// AssignPermissionsToRole 为角色分配权限
func (s *Service) AssignPermissionsToRole(roleID uint, permissionIDs []uint) error {
	// 检查角色是否存在
	if _, err := s.repo.GetRoleByID(roleID); err != nil {
		return errors.New("角色不存在")
	}

	// 检查所有权限是否存在
	for _, permID := range permissionIDs {
		if _, err := s.repo.GetPermissionByID(permID); err != nil {
			return fmt.Errorf("权限ID %d 不存在", permID)
		}
	}

	return s.repo.AssignPermissionsToRole(roleID, permissionIDs)
}

// GetRolePermissions 获取角色的所有权限
func (s *Service) GetRolePermissions(roleID uint) ([]models.Permission, error) {
	return s.repo.GetRolePermissions(roleID)
}

// ========== 用户角色关联 ==========

// AssignRolesToUser 为用户分配角色
func (s *Service) AssignRolesToUser(userID uint, roleIDs []uint) error {
	// 检查所有角色是否存在且启用
	for _, roleID := range roleIDs {
		role, err := s.repo.GetRoleByID(roleID)
		if err != nil {
			return fmt.Errorf("角色ID %d 不存在", roleID)
		}
		if role.Status != 1 {
			return fmt.Errorf("角色 %s 已被禁用", role.DisplayName)
		}
	}

	return s.repo.AssignRolesToUser(userID, roleIDs)
}

// GetUserRoles 获取用户的所有角色
func (s *Service) GetUserRoles(userID uint) ([]models.Role, error) {
	return s.repo.GetUserRoles(userID)
}

// GetUserPermissions 获取用户的所有权限
func (s *Service) GetUserPermissions(userID uint) ([]models.Permission, error) {
	return s.repo.GetUserPermissions(userID)
}

// GetUserFrontendPermissions 获取用户的前端权限(用于菜单渲染)
func (s *Service) GetUserFrontendPermissions(userID uint) ([]models.Permission, error) {
	permissions, err := s.repo.GetUserPermissionsByType(userID, models.PermissionTypeFrontend)
	if err != nil {
		return nil, err
	}
	// 构建树形结构
	return s.buildPermissionTree(permissions), nil
}

// buildPermissionTree 构建权限树形结构
func (s *Service) buildPermissionTree(permissions []models.Permission) []models.Permission {
	// 创建ID到权限的映射
	permMap := make(map[uint]*models.Permission)
	for i := range permissions {
		permMap[permissions[i].ID] = &permissions[i]
		permissions[i].Children = []models.Permission{}
	}

	// 构建树形结构
	var roots []models.Permission
	for i := range permissions {
		perm := &permissions[i]
		if perm.ParentID == nil {
			// 根节点
			roots = append(roots, *perm)
		} else {
			// 子节点，添加到父节点的children中
			if parent, ok := permMap[*perm.ParentID]; ok {
				parent.Children = append(parent.Children, *perm)
			}
		}
	}

	return roots
}

// GetUserBackendPermissions 获取用户的后端API权限
func (s *Service) GetUserBackendPermissions(userID uint) ([]models.Permission, error) {
	return s.repo.GetUserPermissionsByType(userID, models.PermissionTypeBackend)
}

// CheckUserPermission 检查用户是否有某个权限
func (s *Service) CheckUserPermission(userID uint, permissionName string) (bool, error) {
	return s.repo.CheckUserPermission(userID, permissionName)
}

// ========== 管理员管理 ==========

// CreateAdminUser 创建管理员
func (s *Service) CreateAdminUser(username, password, email, nickname, avatar string, status int8) (*models.AdminUser, error) {
	// 检查用户名是否已存在
	if existingUser, _ := s.repo.GetAdminUserByUsername(username); existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if email != "" {
		if existingUser, _ := s.repo.GetAdminUserByEmail(email); existingUser != nil {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &models.AdminUser{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
		Nickname: nickname,
		Avatar:   avatar,
		Status:   status,
	}

	if err := s.repo.CreateAdminUser(user); err != nil {
		return nil, fmt.Errorf("创建管理员失败: %w", err)
	}

	return user, nil
}

// UpdateAdminUser 更新管理员
func (s *Service) UpdateAdminUser(id uint, email, nickname, avatar string, status int8) (*models.AdminUser, error) {
	user, err := s.repo.GetAdminUserByID(id)
	if err != nil {
		return nil, errors.New("管理员不存在")
	}

	// 检查邮箱是否已被其他用户使用
	if email != "" && email != user.Email {
		if existingUser, _ := s.repo.GetAdminUserByEmail(email); existingUser != nil && existingUser.ID != id {
			return nil, errors.New("邮箱已被其他用户使用")
		}
	}

	user.Email = email
	user.Nickname = nickname
	user.Avatar = avatar
	user.Status = status

	if err := s.repo.UpdateAdminUser(user); err != nil {
		return nil, fmt.Errorf("更新管理员失败: %w", err)
	}

	return user, nil
}

// UpdateAdminUserPassword 更新管理员密码
func (s *Service) UpdateAdminUserPassword(id uint, newPassword string) error {
	user, err := s.repo.GetAdminUserByID(id)
	if err != nil {
		return errors.New("管理员不存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("密码加密失败: %w", err)
	}

	user.Password = string(hashedPassword)

	if err := s.repo.UpdateAdminUser(user); err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// DeleteAdminUser 删除管理员
func (s *Service) DeleteAdminUser(id uint) error {
	if _, err := s.repo.GetAdminUserByID(id); err != nil {
		return errors.New("管理员不存在")
	}

	return s.repo.DeleteAdminUser(id)
}

// GetAdminUser 获取管理员详情
func (s *Service) GetAdminUser(id uint) (*models.AdminUser, error) {
	return s.repo.GetAdminUserByID(id)
}

// ListAdminUsers 获取管理员列表
func (s *Service) ListAdminUsers(page, pageSize int, keyword string, status *int8) ([]models.AdminUser, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.ListAdminUsers(page, pageSize, keyword, status)
}
