package services

import (
	"Resume/backend/database"
	"Resume/backend/models"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type CompanyPositionService struct {
	db *gorm.DB
}

func NewCompanyPositionService() *CompanyPositionService {
	return &CompanyPositionService{
		db: database.DB,
	}
}

func (s *CompanyPositionService) CreatePosition(req *models.CompanyPositionCreateRequest) (*models.CompanyPosition, error) {
	position := models.CompanyPosition{
		Name: req.Name,
	}
	result := s.db.Create(&position)
	if result.Error != nil {
		return nil, fmt.Errorf("创建岗位失败: %w", result.Error)
	}
	return &position, nil
}

func (s *CompanyPositionService) GetPositionByID(id string) (*models.CompanyPosition, error) {
	var position models.CompanyPosition
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的岗位ID: %s", id)
	}
	err = s.db.First(&position, "id = ?", uint(uid)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的岗位未找到", id)
		}
		return nil, fmt.Errorf("根据ID '%s' 获取岗位失败: %w", id, err)
	}
	return &position, nil
}

func (s *CompanyPositionService) ListPositions(page, limit int, sortBy, order, keyword string) ([]models.CompanyPosition, int64, error) {
	var positions []models.CompanyPosition
	var total int64
	query := s.db.Model(&models.CompanyPosition{})

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
	err := query.Offset(offset).Limit(limit).Find(&positions).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取岗位列表失败: %w", err)
	}
	return positions, total, nil
}

func (s *CompanyPositionService) UpdatePosition(id string, req *models.CompanyPositionUpdateRequest) (*models.CompanyPosition, error) {
	var position models.CompanyPosition
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的岗位ID: %s", id)
	}
	err = s.db.First(&position, "id = ?", uint(uid)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的岗位未找到，无法更新", id)
		}
		return nil, fmt.Errorf("查找岗位失败，无法更新: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}

	result := s.db.Model(&position).Where("id = ?", uint(uid)).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("更新岗位失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("ID为 '%s' 的岗位未找到或没有更改", id)
	}
	s.db.First(&position, "id = ?", uint(uid))
	return &position, nil
}

func (s *CompanyPositionService) DeletePosition(id string) error {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("无效的岗位ID: %s", id)
	}
	result := s.db.Where("id = ?", uint(uid)).Delete(&models.CompanyPosition{})
	if result.Error != nil {
		return fmt.Errorf("删除岗位失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ID为 '%s' 的岗位未找到，无法删除", id)
	}
	return nil
}
