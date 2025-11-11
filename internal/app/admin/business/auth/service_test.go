package auth

import (
	"databaseAi/internal/app/admin/migrations"
	"databaseAi/internal/infra/config"
	"databaseAi/internal/infra/database"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func setupTestDB(t *testing.T) *database.DB {
	cfg := config.NewConfig("dev")
	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("Failed to connect database: %v", err)
	}

	// 运行迁移
	if err := migrations.MigrateAdmin(db); err != nil {
		t.Fatalf("Failed to migrate: %v", err)
	}

	return db
}

func TestService_Login(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepo(db)
	service := NewService(repo, "test-secret-key")

	// 测试登录
	token, user, err := service.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}

	if user.Username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", user.Username)
	}

	t.Logf("Login successful. Token: %s", token)
	t.Logf("User: %+v", user)
}

func TestService_LoginFailed(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepo(db)
	service := NewService(repo, "test-secret-key")

	// 测试错误密码
	_, _, err := service.Login("admin", "wrong-password")
	if err == nil {
		t.Error("Login should fail with wrong password")
	}

	// 测试不存在的用户
	_, _, err = service.Login("nonexistent", "password")
	if err == nil {
		t.Error("Login should fail with nonexistent user")
	}
}

func TestService_ValidateToken(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepo(db)
	service := NewService(repo, "test-secret-key")

	// 先登录获取 token
	token, _, err := service.Login("admin", "admin123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	// 验证 token
	claims, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("Token validation failed: %v", err)
	}

	username := (*claims)["username"].(string)
	if username != "admin" {
		t.Errorf("Expected username 'admin', got '%s'", username)
	}

	t.Logf("Token validated successfully. Claims: %+v", *claims)
}

func TestService_ChangePassword(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewRepo(db)
	service := NewService(repo, "test-secret-key")

	// 获取用户
	user, err := repo.FindByUsername("admin")
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}

	// 修改密码
	newPassword := "newpassword123"
	err = service.ChangePassword(user.ID, "admin123", newPassword)
	if err != nil {
		t.Fatalf("Failed to change password: %v", err)
	}

	// 验证新密码可以登录
	_, _, err = service.Login("admin", newPassword)
	if err != nil {
		t.Fatalf("Login with new password failed: %v", err)
	}

	// 恢复原密码供其他测试使用
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	repo.UpdateUser(user)

	t.Log("Password change successful")
}
