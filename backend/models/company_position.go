package models

import (
	"time"
)

type CompanyPosition struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name string `gorm:"type:varchar(100);not null;unique" json:"name"`
}

type CompanyPositionCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type CompanyPositionUpdateRequest struct {
	Name *string `json:"name"`
}
