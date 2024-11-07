package models

import (
	"time"
)

type Merchants struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name" binding:"required"`
	Email     string    `json:"email" binding:"required,email"`
	Address   string    `json:"address"`
	IDNumber  string    `json:"id_number" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
}
