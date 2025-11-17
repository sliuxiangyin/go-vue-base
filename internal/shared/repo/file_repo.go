package repo

import (
	"databaseAi/internal/shared/models"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func NewFileRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

// Create 创建文件记录
func (r *Repo) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// GetByID 根据ID获取文件
func (r *Repo) GetByID(id uint) (*models.File, error) {
	var file models.File
	err := r.db.First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetByPath 根据路径获取文件
func (r *Repo) GetByPath(path string) (*models.File, error) {
	var file models.File
	err := r.db.Where("path = ?", path).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// GetByHash 根据哈希获取文件
func (r *Repo) GetByHash(hash string) (*models.File, error) {
	var file models.File
	err := r.db.Where("hash = ?", hash).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// List 获取文件列表（分页）
func (r *Repo) List(page, pageSize int, fileType models.FileType, source models.FileSource, keyword string) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := r.db.Model(&models.File{})

	// 条件筛选
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}
	if source != "" {
		query = query.Where("source = ?", source)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}

// Update 更新文件信息
func (r *Repo) Update(file *models.File) error {
	return r.db.Save(file).Error
}

// Delete 删除文件记录（软删除）
func (r *Repo) Delete(id uint) error {
	return r.db.Delete(&models.File{}, id).Error
}

// IncrementDownloadCount 增加下载次数
func (r *Repo) IncrementDownloadCount(id uint) error {
	return r.db.Model(&models.File{}).Where("id = ?", id).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1)).Error
}

// GetStorageStats 获取存储统计信息
func (r *Repo) GetStorageStats() (map[string]interface{}, error) {
	var result struct {
		TotalFiles int64
		TotalSize  int64
	}

	// 总文件数
	if err := r.db.Model(&models.File{}).Count(&result.TotalFiles).Error; err != nil {
		return nil, err
	}

	// 总大小
	if err := r.db.Model(&models.File{}).Select("COALESCE(SUM(size), 0)").Scan(&result.TotalSize).Error; err != nil {
		return nil, err
	}

	// 按类型统计
	var typeStats []struct {
		FileType models.FileType
		Count    int64
		Size     int64
	}
	if err := r.db.Model(&models.File{}).
		Select("file_type, COUNT(*) as count, COALESCE(SUM(size), 0) as size").
		Group("file_type").
		Scan(&typeStats).Error; err != nil {
		return nil, err
	}

	stats := map[string]interface{}{
		"total_files": result.TotalFiles,
		"total_size":  result.TotalSize,
		"by_type":     typeStats,
	}

	return stats, nil
}

// BatchDelete 批量删除
func (r *Repo) BatchDelete(ids []uint) error {
	return r.db.Delete(&models.File{}, ids).Error
}

// GetOrphanFiles 获取未被引用的文件
func (r *Repo) GetOrphanFiles(page, pageSize int) ([]models.File, int64, error) {
	var files []models.File
	var total int64

	query := r.db.Model(&models.File{}).Where("reference_by = ? OR reference_by IS NULL", "")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		return nil, 0, err
	}

	return files, total, nil
}
