package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Resume/backend/models"
	"Resume/backend/services"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	storageService *services.StorageService
}

func NewUploadHandler(storageService *services.StorageService) *UploadHandler {
	return &UploadHandler{
		storageService: storageService,
	}
}

func (h *UploadHandler) UploadSingle(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "文件上传失败: " + err.Error(),
			},
		})
		return
	}

	fileInfo, err := h.storageService.Save(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.UploadResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "文件保存失败: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "文件上传成功",
		},
		File: fileInfo,
	})
}

func (h *UploadHandler) UploadBatch(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]

	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, models.UploadResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "没有选择文件上传",
			},
		})
		return
	}

	savedFiles, errors := h.storageService.SaveBatch(files)

	if len(errors) > 0 {
		errorMessages := []string{}
		for _, err := range errors {
			errorMessages = append(errorMessages, err.Error())
		}
		c.JSON(http.StatusOK, models.UploadResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "部分文件上传失败",
			},
			Files:  savedFiles,
			Errors: errorMessages,
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: fmt.Sprintf("成功上传 %d 个文件", len(savedFiles)),
		},
		Files: savedFiles,
	})
}

func (h *UploadHandler) DeleteBatch(c *gin.Context) {
	var req struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": fmt.Sprintf("无效的请求载荷: %s", err.Error())})
		return
	}

	deletedIDs, errors := h.storageService.DeleteBatch(req.IDs)

	if len(errors) > 0 {
		errorMessages := []string{}
		for _, err := range errors {
			errorMessages = append(errorMessages, err.Error())
		}
		c.JSON(http.StatusMultiStatus, gin.H{
			"success":     false,
			"message":     "部分文件删除失败",
			"deleted_ids": deletedIDs,
			"errors":      errorMessages,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     fmt.Sprintf("成功删除 %d 个文件", len(deletedIDs)),
		"deleted_ids": deletedIDs,
	})
}

func (h *UploadHandler) ListFiles(c *gin.Context) {
	fileType := c.DefaultQuery("type", "all")
	keyword := c.DefaultQuery("keyword", "")
	sortBy := c.DefaultQuery("sortBy", "created_at")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	files, total, err := h.storageService.ListFiles(fileType, keyword, sortBy, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FileListResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "获取文件列表失败: " + err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, models.FileListResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "文件列表获取成功",
		},
		Files: files,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}

func (h *UploadHandler) GetFile(c *gin.Context) {
	id := c.Param("id")
	fileInfo, err := h.storageService.GetFile(id)
	if err != nil {
		status := http.StatusNotFound
		if strings.Contains(err.Error(), "从数据库获取文件失败") {
			status = http.StatusInternalServerError
		}
		c.JSON(status, models.UploadResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "文件获取失败: " + err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, models.UploadResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "文件获取成功",
		},
		File: fileInfo,
	})
}

func (h *UploadHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	err := h.storageService.DeleteFile(id)
	if err != nil {
		status := http.StatusInternalServerError
		if strings.Contains(err.Error(), "文件未找到") {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"success": false, "message": "文件删除失败: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "文件删除成功"})
}

func (h *UploadHandler) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	fileInfo, err := h.storageService.GetFile(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "文件未找到或已删除: " + err.Error()})
		return
	}

	filePath := h.storageService.GetFilePath(fileInfo)
	c.FileAttachment(filePath, fileInfo.OriginalName)
}

func (h *UploadHandler) SearchFiles(c *gin.Context) {
	keyword := c.DefaultQuery("q", "")
	fileType := c.DefaultQuery("type", "all")

	files, _, err := h.storageService.ListFiles(fileType, keyword, "", 1, 99999)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.FileSearchResponse{
			BasicResponse: models.BasicResponse{
				Success: false,
				Message: "搜索文件失败: " + err.Error(),
			},
		})
		return
	}
	c.JSON(http.StatusOK, models.FileSearchResponse{
		BasicResponse: models.BasicResponse{
			Success: true,
			Message: "文件搜索成功",
		},
		Files: files,
	})
}

func (h *UploadHandler) GetStats(c *gin.Context) {
	totalCount, err := h.storageService.GetTotalFileCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取文件总数失败: " + err.Error()})
		return
	}

	totalSize, err := h.storageService.GetTotalFileSize()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "获取文件总大小失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"total_files": totalCount,
		"total_size":  totalSize,
	})
}

func (h *UploadHandler) RescanFiles(c *gin.Context) {
	err := h.storageService.RescanFiles()
	if err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"success": false, "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "文件重新扫描成功"})
}
