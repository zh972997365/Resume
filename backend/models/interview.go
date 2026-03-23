package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InterviewRound string

const (
	InterviewRoundInitial     InterviewRound = "初试"
	InterviewRoundReinterview InterviewRound = "复试"
	InterviewRoundFinal       InterviewRound = "终面"
)

type InterviewMethod string

const (
	InterviewMethodOnsite InterviewMethod = "现场面试"
	InterviewMethodPhone  InterviewMethod = "电话面试"
	InterviewMethodVideo  InterviewMethod = "视频面试"
)

type Suggestion string

const (
	SuggestionPassed   Suggestion = "通过"
	SuggestionPending  Suggestion = "待定"
	SuggestionRejected Suggestion = "淘汰"
)

type Interview struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	CandidateName string `gorm:"type:varchar(100);not null" json:"candidate_name"`
	PhoneNumber   string `gorm:"type:varchar(20);not null" json:"phone_number"`
	Email         string `gorm:"type:varchar(100);not null" json:"email"`

	CompanyPositionID uint             `gorm:"not null" json:"company_position_id"`
	CompanyPosition   *CompanyPosition `gorm:"foreignkey:CompanyPositionID" json:"company_position,omitempty"`

	ResumeFileID string    `gorm:"type:varchar(36);uniqueIndex" json:"resume_file_id"`
	ResumeFile   *FileInfo `gorm:"foreignkey:ResumeFileID" json:"resume_file,omitempty"`

	RecruitmentSourceID uint               `gorm:"not null" json:"recruitment_source_id"`
	RecruitmentSource   *RecruitmentSource `gorm:"foreignkey:RecruitmentSourceID" json:"recruitment_source,omitempty"`

	InterviewRound  InterviewRound  `gorm:"type:varchar(20);not null" json:"interview_round"`
	InterviewMethod InterviewMethod `gorm:"type:varchar(20);not null" json:"interview_method"`
	InterviewTime   time.Time       `gorm:"type:datetime;not null" json:"interview_time"`

	InterviewerID uint      `gorm:"not null" json:"interviewer_id"`
	Interviewer   *Employee `gorm:"foreignkey:InterviewerID" json:"interviewer,omitempty"`

	InterviewRating int        `gorm:"type:int;default:0" json:"interview_rating"`
	Comments        string     `gorm:"type:text" json:"comments,omitempty"`
	Suggestion      Suggestion `gorm:"type:varchar(20);not null" json:"suggestion"`
}

func (i *Interview) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == "" {
		i.ID = uuid.New().String()
	}
	return
}

type InterviewCreateRequest struct {
	CandidateName       string          `json:"candidate_name" binding:"required"`
	PhoneNumber         string          `json:"phone_number" binding:"required"`
	Email               string          `json:"email" binding:"required"`
	CompanyPositionID   uint            `json:"company_position_id" binding:"required"`
	ResumeFileID        string          `json:"resume_file_id" binding:"required"`
	RecruitmentSourceID uint            `json:"recruitment_source_id" binding:"required"`
	InterviewRound      InterviewRound  `json:"interview_round" binding:"required"`
	InterviewMethod     InterviewMethod `json:"interview_method" binding:"required"`
	InterviewTime       time.Time       `json:"interview_time" binding:"required"`
	InterviewerID       uint            `json:"interviewer_id" binding:"required"`
	InterviewRating     int             `json:"interview_rating"`
	Comments            string          `json:"comments"`
	Suggestion          Suggestion      `json:"suggestion" binding:"required"`
}

type InterviewUpdateRequest struct {
	CandidateName       *string          `json:"candidate_name"`
	PhoneNumber         *string          `json:"phone_number"`
	Email               *string          `json:"email"`
	CompanyPositionID   *uint            `json:"company_position_id"`
	ResumeFileID        *string          `json:"resume_file_id"`
	RecruitmentSourceID *uint            `json:"recruitment_source_id"`
	InterviewRound      *InterviewRound  `json:"interview_round"`
	InterviewMethod     *InterviewMethod `json:"interview_method"`
	InterviewTime       *time.Time       `json:"interview_time"`
	InterviewerID       *uint            `json:"interviewer_id"`
	InterviewRating     *int             `json:"interview_rating"`
	Comments            *string          `json:"comments"`
	Suggestion          *Suggestion      `json:"suggestion"`
}
