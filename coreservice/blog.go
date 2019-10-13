package coreservice

import (
	"github.com/jinzhu/gorm"
	"fmt"
)

type Blog struct {
	gorm.Model
	Title   string
	Content string
}

type Service struct {
	db *gorm.DB
}

func (s Service) create(blog *Blog) error {
	return s.db.Table("blogs").Create(blog).Error
}

func (s Service) getAllBlogs() ([]*Blog, error){
	result := []*Blog{}
	rows, err := s.db.Model(&Blog{}).Select("title, content").Rows()
	if err != nil {
		fmt.Println("error is ",err)
		return nil, nil
	}

	for rows.Next(){
		var title		string
		var content 	string
		if err := rows.Scan(&title, &content); err != nil {
			fmt.Println("scan err: ",err)
		}
		result = append(result, &Blog{
			Title:		title,
			Content:	content,
		})
	}

	return result, nil
}