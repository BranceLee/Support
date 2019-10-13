package coreservice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Blog struct {
	gorm.Model
	UUID    uuid.UUID `gorm:"unique_index; not null" sql:"type:uuid"`
	Title   string    `gorm:"not null"`
	Content string    `gorm:"not null"`
}

type Service struct {
	db *gorm.DB
}

func (s Service) create(blog *Blog) error {
	return s.db.Table("blogs").Create(blog).Error
}

type blogPayload struct {
	Title   string `json:"name"`
	Content string `json:"content"`
	UUID    string `json:"uid"`
}

func (s Service) getAllBlogs() ([]*blogPayload, error) {
	result := []*blogPayload{}
	rows, err := s.db.Model(&Blog{}).Select(`title, content, uuid`).Rows()
	if err != nil {
		fmt.Println("error is ", err)
		return nil, nil
	}

	for rows.Next() {
		var title string
		var content string
		var uid string
		if err := rows.Scan(&title, &content, &uid); err != nil {
			fmt.Println("scan err: ", err)
		}
		result = append(result, &blogPayload{
			Title:   title,
			Content: content,
			UUID:    uid,
		})
	}

	return result, nil
}
