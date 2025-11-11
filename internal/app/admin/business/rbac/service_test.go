package rbac

import (
	"databaseAi/internal/app/admin/models"
	"databaseAi/internal/infra/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *database.DB {
	db, err := database.NewDB("sqlite::memory:")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// 迁移测试表
	err = db.AutoMigrate(
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestRoleService(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)
	service := NewService(repo)

	t.Run("CreateRole", func(t *testing.T) {
		role, err := service.CreateRole("admin", "管理员", "系统管理员", 1)
		assert.NoError(t, err)
		assert.NotNil(t, role)
		assert.Equal(t, "admin", role.Name)
		assert.Equal(t, "管理员", role.DisplayName)
	})

	t.Run("CreateDuplicateRole", func(t *testing.T) {
		_, err := service.CreateRole("admin", "管理员2", "重复角色", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "已存在")
	})

	t.Run("ListRoles", func(t *testing.T) {
		roles, total, err := service.ListRoles(1, 10, "", nil)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))
		assert.NotEmpty(t, roles)
	})

	t.Run("UpdateRole", func(t *testing.T) {
		roles, _, _ := service.ListRoles(1, 1, "", nil)
		if len(roles) > 0 {
			role, err := service.UpdateRole(roles[0].ID, "管理员更新", "更新后的描述", 1)
			assert.NoError(t, err)
			assert.Equal(t, "管理员更新", role.DisplayName)
		}
	})
}

func TestPermissionService(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)
	service := NewService(repo)

	t.Run("CreatePermission", func(t *testing.T) {
		perm, err := service.CreatePermission("user.create", "创建用户", "创建用户权限", "用户管理", 1)
		assert.NoError(t, err)
		assert.NotNil(t, perm)
		assert.Equal(t, "user.create", perm.Name)
	})

	t.Run("ListPermissions", func(t *testing.T) {
		perms, total, err := service.ListPermissions(1, 10, "", nil)
		assert.NoError(t, err)
		assert.Greater(t, total, int64(0))
		assert.NotEmpty(t, perms)
	})
}

func TestRolePermissionAssignment(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)
	service := NewService(repo)

	// 创建角色
	role, _ := service.CreateRole("editor", "编辑者", "内容编辑", 1)

	// 创建权限
	perm1, _ := service.CreatePermission("post.create", "创建文章", "", "内容管理", 1)
	perm2, _ := service.CreatePermission("post.edit", "编辑文章", "", "内容管理", 1)

	t.Run("AssignPermissionsToRole", func(t *testing.T) {
		err := service.AssignPermissionsToRole(role.ID, []uint{perm1.ID, perm2.ID})
		assert.NoError(t, err)
	})

	t.Run("GetRolePermissions", func(t *testing.T) {
		perms, err := service.GetRolePermissions(role.ID)
		assert.NoError(t, err)
		assert.Len(t, perms, 2)
	})
}

func TestUserPermissionCheck(t *testing.T) {
	db := setupTestDB(t)
	repo := NewRepo(db)
	service := NewService(repo)

	// 创建角色和权限
	role, _ := service.CreateRole("tester", "测试者", "测试角色", 1)
	perm, _ := service.CreatePermission("test.run", "运行测试", "", "测试管理", 1)

	// 分配权限给角色
	service.AssignPermissionsToRole(role.ID, []uint{perm.ID})

	// 为用户分配角色（假设用户ID为1）
	userID := uint(1)
	service.AssignRolesToUser(userID, []uint{role.ID})

	t.Run("CheckUserPermission", func(t *testing.T) {
		hasPermission, err := service.CheckUserPermission(userID, "test.run")
		assert.NoError(t, err)
		assert.True(t, hasPermission)
	})

	t.Run("CheckNonExistentPermission", func(t *testing.T) {
		hasPermission, err := service.CheckUserPermission(userID, "test.delete")
		assert.NoError(t, err)
		assert.False(t, hasPermission)
	})
}
