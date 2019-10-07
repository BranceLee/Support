package coreservice

import (
	"github.com/jinzhu/gorm"
	// "github.com/google/uuid"
)

type Blog struct {
	gorm.Model
	// UUID			uuid.UUID 	`gorm:"unique_index;not null" sql:"type:uuid"`
	Title   string
	Content string
}

type Service struct {
	db *gorm.DB
}

func (s Service) create(blog *Blog) error {
	return s.db.Table("blogs").Create(blog).Error
}
