package services

import (
	"Resume/backend/database"
	"Resume/backend/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type InterviewService struct {
	db *gorm.DB
}

func NewInterviewService() *InterviewService {
	return &InterviewService{
		db: database.DB,
	}
}

func (s *InterviewService) CreateInterview(req *models.InterviewCreateRequest) (*models.Interview, error) {
	// 先检查简历文件是否已被其他面试记录使用
	var existingCount int64
	s.db.Model(&models.Interview{}).Where("resume_file_id = ?", req.ResumeFileID).Count(&existingCount)
	if existingCount > 0 {
		return nil, fmt.Errorf("该简历文件已被其他面试记录使用，请选择其他简历文件")
	}

	interview := models.Interview{
		CandidateName:       req.CandidateName,
		PhoneNumber:         req.PhoneNumber,
		Email:               req.Email,
		CompanyPositionID:   req.CompanyPositionID, // 使用外键ID
		ResumeFileID:        req.ResumeFileID,
		RecruitmentSourceID: req.RecruitmentSourceID, // 使用外键ID
		InterviewRound:      req.InterviewRound,
		InterviewMethod:     req.InterviewMethod,
		InterviewTime:       req.InterviewTime,
		InterviewerID:       req.InterviewerID, // 使用外键ID
		InterviewRating:     req.InterviewRating,
		Comments:            req.Comments,
		Suggestion:          req.Suggestion,
	}

	result := s.db.Create(&interview)
	if result.Error != nil {
		// 检查是否是唯一约束错误（双重保险）
		errStr := result.Error.Error()
		if strings.Contains(errStr, "Duplicate entry") &&
			(strings.Contains(errStr, "resume_file_id") || strings.Contains(errStr, "idx_interviews_resume_file_id")) {
			return nil, fmt.Errorf("该简历文件已被其他面试记录使用，请选择其他简历文件")
		}
		// 检查其他可能的约束错误
		if strings.Contains(errStr, "Error 1062") && strings.Contains(errStr, "23000") {
			return nil, fmt.Errorf("数据重复提交，请刷新页面后重试")
		}
		return nil, fmt.Errorf("创建面试记录失败: %w", result.Error)
	}

	// 使用 First 确保加载关联的 ResumeFile 以及所有最新字段
	var createdInterview models.Interview
	s.db.Preload("ResumeFile").
		Preload("CompanyPosition").
		Preload("RecruitmentSource").
		Preload("Interviewer").
		First(&createdInterview, "id = ?", interview.ID)

	return &createdInterview, nil
}

func (s *InterviewService) GetInterviewByID(id string) (*models.Interview, error) {
	var interview models.Interview
	err := s.db.Preload("ResumeFile").Preload("CompanyPosition").Preload("RecruitmentSource").Preload("Interviewer").First(&interview, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的面试记录未找到", id)
		}
		return nil, fmt.Errorf("根据ID '%s' 获取面试记录失败: %w", id, err)
	}
	return &interview, nil
}

func (s *InterviewService) ListInterviews(page, limit int, sortBy, order, keyword string, suggestion models.Suggestion) ([]models.Interview, int64, error) {
	var interviews []models.Interview
	var total int64

	query := s.db.Model(&models.Interview{})

	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.
			Joins("LEFT JOIN company_positions ON interviews.company_position_id = company_positions.id").
			Joins("LEFT JOIN recruitment_sources ON interviews.recruitment_source_id = recruitment_sources.id").
			Joins("LEFT JOIN employees ON interviews.interviewer_id = employees.id").
			Where(
				"interviews.candidate_name LIKE ? OR company_positions.name LIKE ? OR recruitment_sources.name LIKE ? OR employees.name LIKE ?",
				keyword, keyword, keyword, keyword,
			)
	}

	if suggestion != "" {
		query = query.Where("interviews.suggestion = ?", suggestion)
	}

	query.Count(&total)

	if sortBy == "" {
		sortBy = "interviews.created_at"
	} else {
		sortBy = "interviews." + sortBy
	}
	if order == "" {
		order = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, order))

	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).Preload("ResumeFile").Preload("CompanyPosition").Preload("RecruitmentSource").Preload("Interviewer").Find(&interviews).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取面试记录列表失败: %w", err)
	}

	return interviews, total, nil
}

func (s *InterviewService) UpdateInterview(id string, req *models.InterviewUpdateRequest) (*models.Interview, error) {
	var interview models.Interview
	err := s.db.First(&interview, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("ID为 '%s' 的面试记录未找到，无法更新", id)
		}
		return nil, fmt.Errorf("查找面试记录失败，无法更新: %w", err)
	}

	// 如果更新了简历文件ID，需要检查是否已被其他面试使用
	if req.ResumeFileID != nil && *req.ResumeFileID != interview.ResumeFileID {
		var existingCount int64
		s.db.Model(&models.Interview{}).
			Where("resume_file_id = ? AND id != ?", *req.ResumeFileID, id).
			Count(&existingCount)
		if existingCount > 0 {
			return nil, fmt.Errorf("该简历文件已被其他面试记录使用，请选择其他简历文件")
		}
	}

	// 使用 map 或结构体更新以避免零值问题
	updates := make(map[string]interface{})
	if req.CandidateName != nil {
		updates["candidate_name"] = *req.CandidateName
	}
	if req.PhoneNumber != nil {
		updates["phone_number"] = *req.PhoneNumber
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.CompanyPositionID != nil {
		updates["company_position_id"] = *req.CompanyPositionID
	}
	if req.ResumeFileID != nil {
		updates["resume_file_id"] = *req.ResumeFileID
	}
	if req.RecruitmentSourceID != nil {
		updates["recruitment_source_id"] = *req.RecruitmentSourceID
	}
	if req.InterviewRound != nil {
		updates["interview_round"] = *req.InterviewRound
	}
	if req.InterviewMethod != nil {
		updates["interview_method"] = *req.InterviewMethod
	}
	if req.InterviewTime != nil {
		updates["interview_time"] = *req.InterviewTime
	}
	if req.InterviewerID != nil {
		updates["interviewer_id"] = *req.InterviewerID
	}
	if req.InterviewRating != nil {
		updates["interview_rating"] = *req.InterviewRating
	}
	if req.Comments != nil {
		updates["comments"] = *req.Comments
	}
	if req.Suggestion != nil {
		updates["suggestion"] = *req.Suggestion
	}

	// 执行更新
	result := s.db.Model(&interview).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		// 检查是否是唯一约束错误（双重保险）
		errStr := result.Error.Error()
		if strings.Contains(errStr, "Duplicate entry") &&
			(strings.Contains(errStr, "resume_file_id") || strings.Contains(errStr, "idx_interviews_resume_file_id")) {
			return nil, fmt.Errorf("该简历文件已被其他面试记录使用，请选择其他简历文件")
		}
		if strings.Contains(errStr, "Error 1062") && strings.Contains(errStr, "23000") {
			return nil, fmt.Errorf("数据更新失败，可能存在重复数据")
		}
		return nil, fmt.Errorf("更新面试记录失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("ID为 '%s' 的面试记录未找到或没有更改", id)
	}

	// 重新加载更新后的面试记录，包括关联的 ResumeFile
	s.db.Preload("ResumeFile").
		Preload("CompanyPosition").
		Preload("RecruitmentSource").
		Preload("Interviewer").
		First(&interview, "id = ?", id)

	return &interview, nil
}

func (s *InterviewService) DeleteInterview(id string) error {
	result := s.db.Where("id = ?", id).Delete(&models.Interview{})
	if result.Error != nil {
		return fmt.Errorf("删除面试记录失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("ID为 '%s' 的面试记录未找到，无法删除", id)
	}
	return nil
}
