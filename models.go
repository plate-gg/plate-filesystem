package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Path     string    `gorm:"not null"`
	CID      string    `gorm:"not null"`
	FileName string    `gorm:"not null"`
	Size    int64     `gorm:"not null"`
}
