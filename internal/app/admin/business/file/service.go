package file

import (
	"crypto/sha256"
	"databaseAi/internal/shared"
	"databaseAi/internal/shared/models"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type Service struct {
	repo    shared.FileRepository
	storage shared.FileStorage
}

func NewService(repo shared.FileRepository, storage shared.FileStorage) *Service {
	return &Service{
		repo:    repo,
		storage: storage,
	}
}

// UploadFile 上传单个文件
func (s *Service) UploadFile(fileHeader *multipart.FileHeader, uploadedBy uint, referenceBy string, referenceID uint) (*models.File, error) {
	// 打开上传的文件计算哈希
	src, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// 检查是否已存在相同哈希的文件
	existingFile, err := s.repo.GetByHash(hashStr)
	if err == nil && existingFile != nil {
		return existingFile, nil
	}

	// 使用存储接口上传文件
	metadata := &shared.FileMetadata{
		UploadedBy:  uploadedBy,
		ReferenceBy: referenceBy,
		ReferenceID: referenceID,
	}

	relativePath, err := s.storage.UploadFile(fileHeader, metadata)
	if err != nil {
		return nil, err
	}

	// 创建文件记录
	ext := filepath.Ext(fileHeader.Filename)
	mimeType := fileHeader.Header.Get("Content-Type")
	fileType := getFileType(mimeType, ext)

	file := &models.File{
		Name:        fileHeader.Filename,
		Path:        relativePath,
		Size:        fileHeader.Size,
		MimeType:    mimeType,
		Extension:   ext,
		FileType:    fileType,
		Source:      models.FileSourceUpload,
		Hash:        hashStr,
		UploadedBy:  uploadedBy,
		ReferenceBy: referenceBy,
		ReferenceID: referenceID,
		StorageType: "local",
	}

	if err := s.repo.Create(file); err != nil {
		// 如果数据库创建失败，删除已上传的文件
		s.storage.DeleteFile(relativePath)
		return nil, fmt.Errorf("failed to create file record: %w", err)
	}

	return file, nil
}

// UploadChunk 上传文件分片
func (s *Service) UploadChunk(fileHeader *multipart.FileHeader, chunkIndex, totalChunks int, fileID, originalName string, uploadedBy uint) (*models.File, error) {
	chunkInfo := &shared.ChunkInfo{
		ChunkIndex:   chunkIndex,
		TotalChunks:  totalChunks,
		FileID:       fileID,
		OriginalName: originalName,
		UploadedBy:   uploadedBy,
	}

	isLast, relativePath, err := s.storage.UploadChunk(fileHeader, chunkInfo)
	if err != nil {
		return nil, err
	}

	// 如果不是最后一个分片，返回 nil
	if !isLast {
		return nil, nil
	}

	// 计算完整文件的哈希
	// 注意：这里简化处理，实际应该在合并时计算
	ext := filepath.Ext(originalName)
	mimeType := getMimeTypeFromExtension(ext)
	fileType := getFileType(mimeType, ext)

	// 获取文件大小
	fileSize, err := s.storage.GetFileSize(relativePath)
	if err != nil {
		fileSize = 0
	}

	// 创建文件记录
	file := &models.File{
		Name:        originalName,
		Path:        relativePath,
		Size:        fileSize,
		MimeType:    mimeType,
		Extension:   ext,
		FileType:    fileType,
		Source:      models.FileSourceUpload,
		Hash:        "", // 分片上传时哈希可能不准确
		UploadedBy:  uploadedBy,
		StorageType: "local",
	}

	if err := s.repo.Create(file); err != nil {
		s.storage.DeleteFile(relativePath)
		return nil, err
	}

	return file, nil
}

// GetFileList 获取文件列表
func (s *Service) GetFileList(page, pageSize int, fileType models.FileType, source models.FileSource, keyword string) ([]models.File, int64, error) {
	return s.repo.List(page, pageSize, fileType, source, keyword)
}

// GetFileByID 获取文件详情
func (s *Service) GetFileByID(id uint) (*models.File, error) {
	return s.repo.GetByID(id)
}

// UpdateFile 更新文件信息
func (s *Service) UpdateFile(id uint, description string) error {
	file, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	file.Description = description
	return s.repo.Update(file)
}

// DeleteFile 删除文件
func (s *Service) DeleteFile(id uint, deletePhysical bool) error {
	file, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 删除数据库记录
	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// 如果需要，删除物理文件
	if deletePhysical {
		if err := s.storage.DeleteFile(file.Path); err != nil {
			fmt.Printf("Warning: failed to delete physical file %s: %v\n", file.Path, err)
		}
	}

	return nil
}

// BatchDelete 批量删除文件
func (s *Service) BatchDelete(ids []uint, deletePhysical bool) error {
	if deletePhysical {
		for _, id := range ids {
			file, err := s.repo.GetByID(id)
			if err != nil {
				continue
			}
			s.storage.DeleteFile(file.Path)
		}
	}

	return s.repo.BatchDelete(ids)
}

// GetStorageStats 获取存储统计
func (s *Service) GetStorageStats() (map[string]interface{}, error) {
	return s.repo.GetStorageStats()
}

// RecordDownload 记录文件下载
func (s *Service) RecordDownload(id uint) error {
	return s.repo.IncrementDownloadCount(id)
}

// CleanOrphanFiles 清理未被引用的文件
func (s *Service) CleanOrphanFiles(dryRun bool) ([]string, error) {
	files, _, err := s.repo.GetOrphanFiles(1, 1000)
	if err != nil {
		return nil, err
	}

	var cleaned []string
	for _, file := range files {
		// 检查文件创建时间，只清理超过30天的文件
		if time.Since(file.CreatedAt) < 30*24*time.Hour {
			continue
		}

		if !dryRun {
			if err := s.DeleteFile(file.ID, true); err != nil {
				continue
			}
		}
		cleaned = append(cleaned, file.Path)
	}

	return cleaned, nil
}

// ValidateFile 验证文件是否存在且可访问
func (s *Service) ValidateFile(id uint) error {
	file, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if !s.storage.FileExists(file.Path) {
		return errors.New("physical file not found")
	}

	return nil
}

// getFileType 根据MIME类型和扩展名确定文件类型
func getFileType(mimeType, ext string) models.FileType {
	if strings.HasPrefix(mimeType, "image/") {
		return models.FileTypeImage
	}
	if strings.HasPrefix(mimeType, "audio/") {
		return models.FileTypeAudio
	}
	if strings.HasPrefix(mimeType, "video/") {
		return models.FileTypeVideo
	}

	// 根据扩展名判断
	ext = strings.ToLower(ext)
	switch ext {
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".txt":
		return models.FileTypeDocument
	case ".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp":
		return models.FileTypeImage
	case ".mp3", ".wav", ".ogg", ".flac", ".aac":
		return models.FileTypeAudio
	case ".mp4", ".avi", ".mov", ".wmv", ".flv", ".webm":
		return models.FileTypeVideo
	default:
		return models.FileTypeOther
	}
}

// getMimeTypeFromExtension 根据扩展名获取MIME类型
func getMimeTypeFromExtension(ext string) string {
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
