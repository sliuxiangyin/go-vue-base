package storage

import (
	"crypto/sha256"
	"databaseAi/internal/utils"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type FileStorage struct {
	basePath string // 存储基础路径
}

func NewFileStorage(buildEnv string) *FileStorage {
	var basePath string
	if buildEnv == "dev" {
		// 开发环境使用项目根目录
		basePath = filepath.Join(utils.ProjectRoot(), "storage", "audio")
	} else {
		// 生产环境使用当前目录
		basePath = filepath.Join(".", "storage", "audio")
	}

	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		panic(fmt.Sprintf("Failed to create storage directory: %v", err))
	}

	return &FileStorage{
		basePath: basePath,
	}
}

// DownloadFromURL 从 URL 下载文件并保存到本地，返回相对路径
func (fs *FileStorage) DownloadFromURL(url string) (string, error) {
	// 下载文件
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// 读取文件内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// 根据内容生成哈希作为文件名（避免重复）
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// 根据 Content-Type 确定文件扩展名
	ext := getExtensionFromContentType(resp.Header.Get("Content-Type"))
	if ext == "" {
		ext = ".wav" // 默认扩展名
	}

	// 生成文件名：日期_哈希前12位.扩展名
	datePrefix := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s%s", datePrefix, hashStr[:12], ext)
	relativePath := filepath.Join("storage", "audio", filename)
	fullPath := filepath.Join(fs.basePath, filename)

	// 检查文件是否已存在
	if _, err := os.Stat(fullPath); err == nil {
		// 文件已存在，直接返回路径
		return relativePath, nil
	}

	// 保存文件
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return relativePath, nil
}

// SaveBytes 直接保存字节数据到本地，返回相对路径
func (fs *FileStorage) SaveBytes(data []byte, ext string) (string, error) {
	if ext == "" {
		ext = ".wav"
	}
	if ext[0] != '.' {
		ext = "." + ext
	}

	// 根据内容生成哈希
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// 生成文件名
	datePrefix := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s%s", datePrefix, hashStr[:12], ext)
	relativePath := filepath.Join("storage", "audio", filename)
	fullPath := filepath.Join(fs.basePath, filename)

	// 检查文件是否已存在
	if _, err := os.Stat(fullPath); err == nil {
		return relativePath, nil
	}

	// 保存文件
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return relativePath, nil
}

// GetFullPath 获取文件的完整路径
func (fs *FileStorage) GetFullPath(relativePath string) string {
	return filepath.Join(utils.ProjectRoot(), relativePath)
}

// getExtensionFromContentType 根据 Content-Type 获取文件扩展名
func getExtensionFromContentType(contentType string) string {
	switch contentType {
	case "audio/mpeg", "audio/mp3":
		return ".mp3"
	case "audio/wav", "audio/wave":
		return ".wav"
	case "audio/ogg":
		return ".ogg"
	case "audio/flac":
		return ".flac"
	default:
		return ""
	}
}
