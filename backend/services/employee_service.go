package services

import (
	"Resume/backend/database"
	"Resume/backend/models"
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type EmployeeService struct {
	db *gorm.DB
}

func NewEmployeeService() *EmployeeService {
	return &EmployeeService{
		db: database.DB,
	}
}

func (s *EmployeeService) CreateEmployee(req *models.EmployeeCreateRequest) (*models.Employee, error) {
	employee := models.Employee{
		Name:  req.Name,
		Email: req.Email,
	}
	if req.Department != nil {
		employee.Department = *req.Department
	}
	result := s.db.Create(&employee)
	if result.Error != nil {
		return nil, fmt.Errorf("创建面试官失败: %w", result.Error)
	}
	return &employee, nil
}

func (s *EmployeeService) GetEmployeeByID(id string) (*models.Employee, error) {
	var employee models.Employee
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的面试官ID: %s", id)
	}
	err = s.db.First(&employee, "id = ?", uint(uid)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的面试官未找到", id)
		}
		return nil, fmt.Errorf("根据ID '%s' 获取面试官失败: %w", id, err)
	}
	return &employee, nil
}

func (s *EmployeeService) ListEmployees(page, limit int, sortBy, order, keyword string) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64
	query := s.db.Model(&models.Employee{})

	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where("name LIKE ? OR email LIKE ? OR department LIKE ?", keyword, keyword, keyword)
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
	err := query.Offset(offset).Limit(limit).Find(&employees).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取面试官列表失败: %w", err)
	}
	return employees, total, nil
}

func (s *EmployeeService) UpdateEmployee(id string, req *models.EmployeeUpdateRequest) (*models.Employee, error) {
	var employee models.Employee
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("无效的面试官ID: %s", id)
	}
	err = s.db.First(&employee, "id = ?", uint(uid)).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的面试官未找到，无法更新", id)
		}
		return nil, fmt.Errorf("查找面试官失败，无法更新: %w", err)
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Department != nil {
		updates["department"] = *req.Department
	}

	result := s.db.Model(&employee).Where("id = ?", uint(uid)).Updates(updates)
	if result.Error != nil {
		return nil, fmt.Errorf("更新面试官失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("ID为 '%s' 的面试官未找到或没有更改", id)
	}
	s.db.First(&employee, "id = ?", uint(uid))
	return &employee, nil
}

func (s *EmployeeService) DeleteEmployee(id string) error {
	uid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		return fmt.Errorf("无效的面试官ID: %s", id)
	}
	result := s.db.Where("id = ?", uint(uid)).Delete(&models.Employee{})
	if result.Error != nil {
		return fmt.Errorf("删除面试官失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ID为 '%s' 的面试官未找到，无法删除", id)
	}
	return nil
}
