package shared

import (
	"databaseAi/internal/shared/models"
	"mime/multipart"
)

// FileStorage 文件存储接口，用于抽象不同的存储实现（本地、OSS、S3等）
type FileStorage interface {
	// UploadFile 上传单个文件
	// 返回：文件相对路径、错误
	UploadFile(file *multipart.FileHeader, metadata *FileMetadata) (string, error)

	// UploadChunk 上传文件分片
	// 返回：是否为最后一片、文件路径（仅最后一片）、错误
	UploadChunk(chunk *multipart.FileHeader, chunkInfo *ChunkInfo) (bool, string, error)

	// DeleteFile 删除文件
	DeleteFile(path string) error

	// GetFileURL 获取文件访问URL
	GetFileURL(path string) string

	// FileExists 检查文件是否存在
	FileExists(path string) bool

	// GetFileSize 获取文件大小
	GetFileSize(path string) (int64, error)
}

// FileMetadata 文件元数据
type FileMetadata struct {
	UploadedBy  uint   // 上传者ID
	ReferenceBy string // 引用模块
	ReferenceID uint   // 引用ID
	Description string // 描述
}

// ChunkInfo 分片信息
type ChunkInfo struct {
	ChunkIndex   int    // 当前分片索引
	TotalChunks  int    // 总分片数
	FileID       string // 文件唯一标识
	OriginalName string // 原始文件名
	UploadedBy   uint   // 上传者ID
}

// FileRepository 文件仓库接口
type FileRepository interface {
	Create(file *models.File) error
	GetByID(id uint) (*models.File, error)
	GetByPath(path string) (*models.File, error)
	GetByHash(hash string) (*models.File, error)
	List(page, pageSize int, fileType models.FileType, source models.FileSource, keyword string) ([]models.File, int64, error)
	Update(file *models.File) error
	Delete(id uint) error
	BatchDelete(ids []uint) error
	IncrementDownloadCount(id uint) error
	GetStorageStats() (map[string]interface{}, error)
	GetOrphanFiles(page, pageSize int) ([]models.File, int64, error)
}
