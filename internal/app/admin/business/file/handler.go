package file

import (
	"databaseAi/internal/app/admin/business/auth"
	"databaseAi/internal/shared/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service     *Service
	authService *auth.Service
}

func NewHandler(service *Service, authService *auth.Service) *Handler {
	return &Handler{
		service:     service,
		authService: authService,
	}
}

// UploadFile 上传单个文件
func (h *Handler) UploadFile(c *fiber.Ctx) error {
	// 获取当前用户ID
	userID := h.getUserID(c)

	// 获取文件
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "No file uploaded",
		})
	}

	// 获取可选参数
	referenceBy := c.FormValue("reference_by", "")
	referenceIDStr := c.FormValue("reference_id", "0")
	referenceID, _ := strconv.ParseUint(referenceIDStr, 10, 32)

	// 上传文件
	uploadedFile, err := h.service.UploadFile(file, userID, referenceBy, uint(referenceID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"file_id":   uploadedFile.ID,
			"file_path": uploadedFile.Path,
			"file_name": uploadedFile.Name,
			"file_size": uploadedFile.Size,
			"file_type": uploadedFile.FileType,
		},
	})
}

// UploadChunk 上传文件分片
func (h *Handler) UploadChunk(c *fiber.Ctx) error {
	userID := h.getUserID(c)

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "No file chunk uploaded",
		})
	}

	// 获取分片信息
	chunkIndexStr := c.FormValue("chunk_index")
	totalChunksStr := c.FormValue("total_chunks")
	fileID := c.FormValue("file_id")
	originalName := c.FormValue("original_name")

	chunkIndex, _ := strconv.Atoi(chunkIndexStr)
	totalChunks, _ := strconv.Atoi(totalChunksStr)

	uploadedFile, err := h.service.UploadChunk(file, chunkIndex, totalChunks, fileID, originalName, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	// 如果是最后一个分片，返回完整文件信息
	if uploadedFile != nil {
		return c.JSON(fiber.Map{
			"code":    0,
			"message": "File uploaded successfully",
			"data": fiber.Map{
				"file_id":   uploadedFile.ID,
				"file_path": uploadedFile.Path,
				"file_name": uploadedFile.Name,
				"file_size": uploadedFile.Size,
				"file_type": uploadedFile.FileType,
			},
		})
	}

	// 中间分片
	return c.JSON(fiber.Map{
		"code":    0,
		"message": "Chunk uploaded successfully",
		"data": fiber.Map{
			"chunk_uploaded": true,
		},
	})
}

// GetFileList 获取文件列表
func (h *Handler) GetFileList(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	fileType := models.FileType(c.Query("file_type", ""))
	source := models.FileSource(c.Query("source", ""))
	keyword := c.Query("keyword", "")

	files, total, err := h.service.GetFileList(page, pageSize, fileType, source, keyword)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"list":      files,
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

// GetFileDetail 获取文件详情
func (h *Handler) GetFileDetail(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid file ID",
		})
	}

	file, err := h.service.GetFileByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "File not found",
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data":    file,
	})
}

// UpdateFile 更新文件信息
func (h *Handler) UpdateFile(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid file ID",
		})
	}

	var req struct {
		Description string `json:"description"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid request body",
		})
	}

	if err := h.service.UpdateFile(uint(id), req.Description); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "File updated successfully",
	})
}

// DeleteFile 删除文件
func (h *Handler) DeleteFile(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid file ID",
		})
	}

	deletePhysical := c.QueryBool("delete_physical", true)

	if err := h.service.DeleteFile(uint(id), deletePhysical); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "File deleted successfully",
	})
}

// BatchDelete 批量删除文件
func (h *Handler) BatchDelete(c *fiber.Ctx) error {
	var req struct {
		IDs            []uint `json:"ids"`
		DeletePhysical bool   `json:"delete_physical"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid request body",
		})
	}

	if len(req.IDs) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "No file IDs provided",
		})
	}

	if err := h.service.BatchDelete(req.IDs, req.DeletePhysical); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "Files deleted successfully",
	})
}

// GetStorageStats 获取存储统计
func (h *Handler) GetStorageStats(c *fiber.Ctx) error {
	stats, err := h.service.GetStorageStats()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data":    stats,
	})
}

// DownloadFile 下载文件（记录下载次数）
func (h *Handler) DownloadFile(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Invalid file ID",
		})
	}

	file, err := h.service.GetFileByID(uint(id))
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "File not found",
		})
	}

	// 记录下载次数
	go h.service.RecordDownload(uint(id))

	// 返回文件路径，由前端或静态文件服务器处理实际下载
	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"file_path": file.Path,
			"file_name": file.Name,
		},
	})
}

// CleanOrphanFiles 清理孤立文件
func (h *Handler) CleanOrphanFiles(c *fiber.Ctx) error {
	dryRun := c.QueryBool("dry_run", true)

	cleaned, err := h.service.CleanOrphanFiles(dryRun)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"code":    500,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"code":    0,
		"message": "success",
		"data": fiber.Map{
			"cleaned_files": cleaned,
			"count":         len(cleaned),
			"dry_run":       dryRun,
		},
	})
}

// getUserID 从上下文获取用户ID
func (h *Handler) getUserID(c *fiber.Ctx) uint {
	userID := c.Locals("user_id") // 注意：与中间件保持一致，使用 user_id
	if userID == nil {
		return 0
	}
	if id, ok := userID.(uint); ok {
		return id
	}
	if id, ok := userID.(float64); ok {
		return uint(id)
	}
	// 尝试从字符串解析
	if idStr, ok := userID.(string); ok {
		id, _ := strconv.ParseUint(idStr, 10, 32)
		return uint(id)
	}
	return 0
}

// RegisterRoutes 注册路由
func (h *Handler) RegisterRoutes(router fiber.Router) {
	files := router.Group("/files")

	// 文件上传
	files.Post("/upload", h.UploadFile)
	files.Post("/upload/chunk", h.UploadChunk)

	// 文件管理
	files.Get("", h.GetFileList)
	files.Get("/:id", h.GetFileDetail)
	files.Put("/:id", h.UpdateFile)
	files.Delete("/:id", h.DeleteFile)
	files.Post("/batch-delete", h.BatchDelete)

	// 文件下载
	files.Get("/:id/download", h.DownloadFile)

	// 统计和清理
	files.Get("/stats", h.GetStorageStats)
	files.Post("/clean-orphan", h.CleanOrphanFiles)
}
