package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"Resume/backend/database"
	"Resume/backend/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StorageService struct {
	basePath     string
	baseURL      string
	maxFileSize  int64
	allowedTypes []string
	db           *gorm.DB
}

func NewStorageService(basePath, baseURL string, maxFileSize int64) (*StorageService, error) {
	absPath, err := filepath.Abs(basePath)
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(absPath, 0755); err != nil {
		return nil, fmt.Errorf("创建存储目录失败: %v", err)
	}

	fmt.Printf("存储服务已初始化. 基础路径: %s, 基础URL: %s\n", absPath, baseURL)

	return &StorageService{
		basePath:    absPath,
		baseURL:     baseURL,
		maxFileSize: maxFileSize,
		allowedTypes: []string{
			"application/pdf",
			"application/msword",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
		db: database.DB,
	}, nil
}

func (s *StorageService) Save(file *multipart.FileHeader) (*models.FileInfo, error) {
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	fileInfo := s.createFileInfo(file)
	destPath := filepath.Join(s.basePath, fileInfo.StoragePath)
	destDir := filepath.Dir(destPath)

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	fmt.Printf("保存文件到磁盘: %s -> %s\n", file.Filename, destPath)

	if err := s.saveFileToDisk(file, destPath); err != nil {
		return nil, err
	}

	hash, err := s.calculateFileHash(destPath)
	if err != nil {
		fmt.Printf("警告: 计算文件哈希失败 %s: %v\n", file.Filename, err)
		hash = ""
	}
	fileInfo.Hash = hash

	if err := s.db.Create(fileInfo).Error; err != nil {
		os.Remove(destPath)
		return nil, fmt.Errorf("保存文件信息到数据库失败: %w", err)
	}

	fmt.Printf("文件成功保存: %s (ID: %s) 到文件系统和数据库\n", file.Filename, fileInfo.ID)
	return fileInfo, nil
}

func (s *StorageService) SaveBatch(files []*multipart.FileHeader) ([]*models.FileInfo, []error) {
	var results []*models.FileInfo
	var errors []error

	for _, file := range files {
		fileInfo, err := s.Save(file)
		if err != nil {
			errMessage := fmt.Errorf("%s: %v", file.Filename, err)
			errors = append(errors, errMessage)
			fmt.Println(errMessage)
			continue
		}
		results = append(results, fileInfo)
	}
	return results, errors
}

func (s *StorageService) ListFiles(fileType, keyword, sortBy string, page, limit int) ([]*models.FileInfo, int64, error) {
	var files []*models.FileInfo
	var total int64

	query := s.db.Model(&models.FileInfo{})

	if fileType != "" && fileType != "all" {
		if fileType == "word" {
			query = query.Where("extension IN (?)", []string{"doc", "docx"})
		} else {
			query = query.Where("extension = ?", fileType)
		}
	}

	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("original_name LIKE ?", keyword)
	}

	query.Count(&total)

	if sortBy == "" {
		sortBy = "created_at desc"
	}
	query = query.Order(sortBy)

	offset := (page - 1) * limit
	if limit <= 0 {
		limit = 10
	}
	err := query.Offset(offset).Limit(limit).Find(&files).Error
	if err != nil {
		return nil, 0, fmt.Errorf("从数据库获取文件列表失败: %w", err)
	}
	return files, total, nil
}

func (s *StorageService) GetFile(id string) (*models.FileInfo, error) {
	var fileInfo models.FileInfo
	err := s.db.First(&fileInfo, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("数据库中未找到文件")
		}
		return nil, fmt.Errorf("从数据库获取文件失败: %w", err)
	}

	filePath := filepath.Join(s.basePath, fileInfo.StoragePath)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		s.db.Where("id = ?", id).Delete(&models.FileInfo{})
		return nil, fmt.Errorf("磁盘上未找到文件，元数据已移除")
	}

	return &fileInfo, nil
}

func (s *StorageService) DeleteFile(id string) error {
	var fileInfo models.FileInfo
	err := s.db.First(&fileInfo, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("未找到要删除的文件")
		}
		return fmt.Errorf("在数据库中查找文件失败: %w", err)
	}

	filePath := filepath.Join(s.basePath, fileInfo.StoragePath)
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除物理文件失败: %v", err)
	}

	result := s.db.Where("id = ?", id).Delete(&models.FileInfo{})
	if result.Error != nil {
		return fmt.Errorf("从数据库删除文件信息失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("尽管已初步检查，但在数据库中未找到要删除的文件")
	}

	fmt.Printf("文件成功删除: %s (ID: %s) 从文件系统和数据库\n", fileInfo.OriginalName, id)
	return nil
}

func (s *StorageService) GetFilePath(fileInfo *models.FileInfo) string {
	return filepath.Join(s.basePath, fileInfo.StoragePath)
}

func (s *StorageService) GetTotalFileSize() (int64, error) {
	var totalSize int64
	err := s.db.Model(&models.FileInfo{}).Select("sum(size)").Row().Scan(&totalSize)
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, fmt.Errorf("获取文件总大小失败: %w", err)
	}
	return totalSize, nil
}

func (s *StorageService) GetTotalFileCount() (int64, error) {
	var count int64
	err := s.db.Model(&models.FileInfo{}).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取文件总数失败: %w", err)
	}
	return count, nil
}

func (s *StorageService) RescanFiles() error {
	return fmt.Errorf("重新扫描功能尚未实现")
}

func (s *StorageService) DeleteBatch(ids []string) ([]string, []error) {
	var deletedIDs []string
	var deletionErrors []error

	for _, id := range ids {
		err := s.DeleteFile(id)
		if err != nil {
			deletionErrors = append(deletionErrors, fmt.Errorf("删除文件 %s 失败: %w", id, err))
		} else {
			deletedIDs = append(deletedIDs, id)
		}
	}
	return deletedIDs, deletionErrors
}

func (s *StorageService) validateFile(file *multipart.FileHeader) error {
	contentType := file.Header.Get("Content-Type")

	validByContentType := false
	for _, allowedType := range s.allowedTypes {
		if contentType == allowedType {
			validByContentType = true
			break
		}
	}

	if !validByContentType {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		validExts := []string{".pdf", ".doc", ".docx"}
		validByExtension := false
		for _, validExt := range validExts {
			if ext == validExt {
				validByExtension = true
				break
			}
		}

		if !validByExtension {
			return fmt.Errorf("不支持的文件类型: %s (文件扩展名: %s，仅支持 PDF, DOC, DOCX)", contentType, ext)
		}
	}

	if file.Size > s.maxFileSize {
		return fmt.Errorf("文件过大: %d 字节 (最大 %d MB)", file.Size, s.maxFileSize/(1024*1024))
	}

	return nil
}

func (s *StorageService) createFileInfo(file *multipart.FileHeader) *models.FileInfo {
	ext := strings.ToLower(filepath.Ext(file.Filename))
	originalName := file.Filename

	// 生成一个更独特的存储文件名，可以包含UUID的一部分和时间戳
	storageName := fmt.Sprintf("%s_%d%s",
		uuid.New().String()[:8], // UUID前8位
		time.Now().UnixNano(),   // 纳秒级时间戳
		ext)

	now := time.Now()
	// 按年/月组织文件目录
	yearMonthDir := filepath.Join(
		fmt.Sprintf("%d", now.Year()),
		fmt.Sprintf("%02d", now.Month()),
	)

	// 完整的相对路径（包含年月目录）
	relativePath := filepath.Join(yearMonthDir, storageName)

	// 为URL路径创建正确的格式
	urlPath := strings.ReplaceAll(relativePath, "\\", "/")

	return &models.FileInfo{
		ID:           uuid.New().String(),
		OriginalName: originalName,
		StorageName:  storageName,
		StoragePath:  relativePath,
		Size:         file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		Extension:    strings.TrimPrefix(ext, "."),
		Hash:         "",
		URL:          "/uploads/" + urlPath,
	}
}

func (s *StorageService) saveFileToDisk(file *multipart.FileHeader, destPath string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开源文件失败: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("复制文件内容失败: %v", err)
	}

	return nil
}

func (s *StorageService) calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
