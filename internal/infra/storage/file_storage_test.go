package storage

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileStorage_SaveBytes(t *testing.T) {
	// 创建临时测试目录
	tempDir := t.TempDir()

	// 创建测试用的 FileStorage
	fs := &FileStorage{
		basePath: tempDir,
	}

	// 测试数据
	testData := []byte("test audio content")

	// 保存文件
	relativePath, err := fs.SaveBytes(testData, ".mp3")
	if err != nil {
		t.Fatalf("SaveBytes failed: %v", err)
	}

	// 验证相对路径格式
	if relativePath == "" {
		t.Error("relative path is empty")
	}

	// 验证文件是否存在
	fullPath := filepath.Join(tempDir, filepath.Base(relativePath))
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		t.Errorf("file does not exist at %s", fullPath)
	}

	// 验证文件内容
	savedData, err := os.ReadFile(fullPath)
	if err != nil {
		t.Fatalf("failed to read saved file: %v", err)
	}

	if string(savedData) != string(testData) {
		t.Errorf("saved data does not match: got %s, want %s", savedData, testData)
	}

	t.Logf("Saved file to: %s", relativePath)
}

func TestFileStorage_SaveBytes_Duplicate(t *testing.T) {
	tempDir := t.TempDir()
	fs := &FileStorage{
		basePath: tempDir,
	}

	testData := []byte("duplicate test")

	// 第一次保存
	path1, err := fs.SaveBytes(testData, ".wav")
	if err != nil {
		t.Fatalf("first save failed: %v", err)
	}

	// 第二次保存相同内容
	path2, err := fs.SaveBytes(testData, ".wav")
	if err != nil {
		t.Fatalf("second save failed: %v", err)
	}

	// 应该返回相同的路径（因为内容相同）
	if path1 != path2 {
		t.Errorf("paths should be the same for duplicate content: got %s and %s", path1, path2)
	}

	t.Logf("Duplicate file returned same path: %s", path1)
}
