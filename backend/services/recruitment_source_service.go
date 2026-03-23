package services

import (
	"Resume/backend/database"
	"Resume/backend/models"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type RecruitmentSourceService struct {
	db *gorm.DB
}

func NewRecruitmentSourceService() *RecruitmentSourceService {
	return &RecruitmentSourceService{
		db: database.DB,
	}
}

func (s *RecruitmentSourceService) CreateSource(req *models.RecruitmentSourceCreateRequest) (*models.RecruitmentSource, error) {
	source := models.RecruitmentSource{
		Name: req.Name,
	}
	result := s.db.Create(&source)
	if result.Error != nil {
		return nil, fmt.Errorf("创建招聘来源失败: %w", result.Error)
	}
	return &source, nil
}

func (s *RecruitmentSourceService) GetSourceByID(id string) (*models.RecruitmentSource, error) {
	var source models.RecruitmentSource
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的招聘来源ID: %s", id)
	}
	err = s.db.First(&source, "id = ?", uint(uid)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的招聘来源未找到", id)
		}
		return nil, fmt.Errorf("根据ID '%s' 获取招聘来源失败: %w", id, err)
	}
	return &source, nil
}

func (s *RecruitmentSourceService) ListSources(page, limit int, sortBy, order, keyword string) ([]models.RecruitmentSource, int64, error) {
	var sources []models.RecruitmentSource
	var total int64
	query := s.db.Model(&models.RecruitmentSource{})

	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ?", keyword)
	}

	query.Count(&total)

	if sortBy == "" {
		sortBy = "created_at"
	}
	if order == "" {
		order = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, order))

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Find(&sources).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取招聘来源列表失败: %w", err)
	}
	return sources, total, nil
}

func (s *RecruitmentSourceService) UpdateSource(id string, req *models.RecruitmentSourceUpdateRequest) (*models.RecruitmentSource, error) {
	var source models.RecruitmentSource
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的招聘来源ID: %s", id)
	}
	err = s.db.First(&source, "id = ?", uint(uid)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的招聘来源未找到，无法更新", id)
		}
		return nil, fmt.Errorf("查找招聘来源失败，无法更新: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}

	result := s.db.Model(&source).Where("id = ?", uint(uid)).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("更新招聘来源失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("ID为 '%s' 的招聘来源未找到或没有更改", id)
	}
	s.db.First(&source, "id = ?", uint(uid))
	return &source, nil
}

func (s *RecruitmentSourceService) DeleteSource(id string) error {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("无效的招聘来源ID: %s", id)
	}
	result := s.db.Where("id = ?", uint(uid)).Delete(&models.RecruitmentSource{})
	if result.Error != nil {
		return fmt.Errorf("删除招聘来源失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ID为 '%s' 的招聘来源未找到，无法删除", id)
	}
	return nil
}
