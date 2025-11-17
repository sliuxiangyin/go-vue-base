package models

import (
	"time"

	"gorm.io/gorm"
)

// FileType 文件类型
type FileType string

const (
	FileTypeImage    FileType = "image"
	FileTypeAudio    FileType = "audio"
	FileTypeVideo    FileType = "video"
	FileTypeDocument FileType = "document"
	FileTypeOther    FileType = "other"
)

// FileSource 文件来源
type FileSource string

const (
	FileSourceUpload   FileSource = "upload"   // 用户上传
	FileSourceDownload FileSource = "download" // 从URL下载
	FileSourceGenerate FileSource = "generate" // 系统生成
)

// File 文件管理模型
type File struct {
	ID            uint           `gorm:"primarykey" json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Name          string         `gorm:"size:255;not null" json:"name"`                // 原始文件名
	Path          string         `gorm:"size:500;not null;uniqueIndex" json:"path"`    // 相对路径
	Size          int64          `gorm:"not null" json:"size"`                         // 文件大小（字节）
	MimeType      string         `gorm:"size:100" json:"mime_type"`                    // MIME类型
	Extension     string         `gorm:"size:20" json:"extension"`                     // 文件扩展名
	FileType      FileType       `gorm:"size:20;index" json:"file_type"`               // 文件类型分类
	Source        FileSource     `gorm:"size:20;default:'upload'" json:"source"`       // 文件来源
	Hash          string         `gorm:"size:64;index" json:"hash"`                    // 文件哈希（用于去重）
	Description   string         `gorm:"size:500" json:"description,omitempty"`        // 文件描述
	UploadedBy    uint           `gorm:"index" json:"uploaded_by,omitempty"`           // 上传者ID
	ReferenceBy   string         `gorm:"size:100;index" json:"reference_by,omitempty"` // 被引用的模块（如lesson、user等）
	ReferenceID   uint           `gorm:"index" json:"reference_id,omitempty"`          // 被引用的记录ID
	DownloadURL   string         `gorm:"size:1000" json:"download_url,omitempty"`      // 原始下载URL（如果是从URL下载的）
	DownloadCount int            `gorm:"default:0" json:"download_count"`              // 下载次数
	StorageType   string         `gorm:"size:20;default:'local'" json:"storage_type"`  // 存储类型：local/oss/s3等
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}
