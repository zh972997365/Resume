package models

import "time"

type RecruitmentSource struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`
}

type RecruitmentSourceCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type RecruitmentSourceUpdateRequest struct {
	Name *string `json:"name"`
}
