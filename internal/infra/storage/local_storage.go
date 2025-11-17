package storage

import (
	"crypto/sha256"
	"databaseAi/internal/shared"
	"databaseAi/internal/utils"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// LocalFileStorage 本地文件存储实现
type LocalFileStorage struct {
	basePath string
	baseURL  string
	buildEnv string
}

// NewLocalFileStorage 创建本地文件存储实例
func NewLocalFileStorage(buildEnv string) *LocalFileStorage {
	var basePath string
	if buildEnv == "dev" {
		basePath = filepath.Join(utils.ProjectRoot(), "storage")
	} else {
		basePath = filepath.Join(".", "storage")
	}

	// 确保存储目录存在
	os.MkdirAll(filepath.Join(basePath, "audio"), 0755)
	os.MkdirAll(filepath.Join(basePath, "images"), 0755)
	os.MkdirAll(filepath.Join(basePath, "videos"), 0755)
	os.MkdirAll(filepath.Join(basePath, "documents"), 0755)
	os.MkdirAll(filepath.Join(basePath, "others"), 0755)
	os.MkdirAll(filepath.Join(basePath, "temp"), 0755)

	return &LocalFileStorage{
		basePath: basePath,
		baseURL:  "/storage", // 相对URL前缀
		buildEnv: buildEnv,
	}
}

// UploadFile 实现文件上传
func (s *LocalFileStorage) UploadFile(fileHeader *multipart.FileHeader, metadata *shared.FileMetadata) (string, error) {
	// 打开上传的文件
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 读取文件内容
	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// 计算文件哈希
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// 确定文件类型和存储路径
	ext := filepath.Ext(fileHeader.Filename)
	mimeType := fileHeader.Header.Get("Content-Type")
	subDir := s.getSubDirectory(mimeType, ext)

	// 生成唯一文件名
	datePrefix := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s%s", datePrefix, hashStr[:12], ext)
	relativePath := filepath.Join("storage", subDir, filename)
	fullPath := filepath.Join(s.basePath, subDir, filename)

	// 保存文件
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return relativePath, nil
}

// UploadChunk 实现分片上传
func (s *LocalFileStorage) UploadChunk(fileHeader *multipart.FileHeader, chunkInfo *shared.ChunkInfo) (bool, string, error) {
	// 分片临时存储目录
	tempDir := filepath.Join(s.basePath, "temp", chunkInfo.FileID)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return false, "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	// 保存分片
	chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", chunkInfo.ChunkIndex))
	src, err := fileHeader.Open()
	if err != nil {
		return false, "", err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return false, "", err
	}

	if err := os.WriteFile(chunkPath, data, 0644); err != nil {
		return false, "", err
	}

	// 如果是最后一个分片，合并所有分片
	if chunkInfo.ChunkIndex == chunkInfo.TotalChunks-1 {
		relativePath, err := s.mergeChunks(tempDir, chunkInfo.TotalChunks, chunkInfo.OriginalName)
		if err != nil {
			return false, "", err
		}
		return true, relativePath, nil
	}

	return false, "", nil
}

// mergeChunks 合并分片
func (s *LocalFileStorage) mergeChunks(tempDir string, totalChunks int, originalName string) (string, error) {
	ext := filepath.Ext(originalName)
	mimeType := s.getMimeTypeFromExtension(ext)
	subDir := s.getSubDirectory(mimeType, ext)

	// 临时合并文件
	mergedPath := filepath.Join(tempDir, "merged"+ext)
	merged, err := os.Create(mergedPath)
	if err != nil {
		return "", err
	}
	defer merged.Close()

	// 合并所有分片
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d", i))
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return "", err
		}
		if _, err := merged.Write(chunkData); err != nil {
			return "", err
		}
	}

	// 读取完整文件计算哈希
	merged.Seek(0, 0)
	data, err := io.ReadAll(merged)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// 移动到正式存储位置
	datePrefix := time.Now().Format("20060102")
	filename := fmt.Sprintf("%s_%s%s", datePrefix, hashStr[:12], ext)
	relativePath := filepath.Join("storage", subDir, filename)
	fullPath := filepath.Join(s.basePath, subDir, filename)

	if err := os.Rename(mergedPath, fullPath); err != nil {
		return "", err
	}

	// 清理临时目录
	os.RemoveAll(tempDir)

	return relativePath, nil
}

// DeleteFile 删除文件
func (s *LocalFileStorage) DeleteFile(path string) error {
	fullPath := s.getFullPath(path)
	return os.Remove(fullPath)
}

// GetFileURL 获取文件URL
func (s *LocalFileStorage) GetFileURL(path string) string {
	// 返回相对路径，由前端拼接基础URL
	return "/" + path
}

// FileExists 检查文件是否存在
func (s *LocalFileStorage) FileExists(path string) bool {
	fullPath := s.getFullPath(path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetFileSize 获取文件大小
func (s *LocalFileStorage) GetFileSize(path string) (int64, error) {
	fullPath := s.getFullPath(path)
	info, err := os.Stat(fullPath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// getSubDirectory 根据MIME类型获取子目录
func (s *LocalFileStorage) getSubDirectory(mimeType, ext string) string {
	if strings.HasPrefix(mimeType, "image/") {
		return "images"
	}
	if strings.HasPrefix(mimeType, "audio/") {
		return "audio"
	}
	if strings.HasPrefix(mimeType, "video/") {
		return "videos"
	}

	// 根据扩展名判断
	ext = strings.ToLower(ext)
	switch ext {
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt":
		return "documents"
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return "images"
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return "audio"
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm":
		return "videos"
	default:
		return "others"
	}
}

// getMimeTypeFromExtension 根据扩展名获取MIME类型
func (s *LocalFileStorage) getMimeTypeFromExtension(ext string) string {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".mp4":
		return "video/mp4"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}

// getFullPath 获取文件完整路径
func (s *LocalFileStorage) getFullPath(relativePath string) string {
	if s.buildEnv == "dev" {
		return filepath.Join(utils.ProjectRoot(), relativePath)
	}
	return relativePath
}
