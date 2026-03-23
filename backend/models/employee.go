package models

import (
	"time"
)

type Employee struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name       string `gorm:"type:varchar(50);not null" json:"name"`
	Email      string `gorm:"type:varchar(100);unique" json:"email"`
	Department string `gorm:"type:varchar(50)" json:"department"`
}

type EmployeeCreateRequest struct {
	Name       string  `json:"name" binding:"required"`
	Email      string  `json:"email" binding:"required,email"`
	Department *string `json:"department"`
}

type EmployeeUpdateRequest struct {
	Name       *string `json:"name"`
	Email      *string `json:"email"`
	Department *string `json:"department"`
}
