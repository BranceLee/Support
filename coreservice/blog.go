package coreservice

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Blog struct {
	gorm.Model
	UUID    uuid.UUID `gorm:"unique_index; not null" sql:"type:uuid"`
	Title   string    `gorm:"not null"`
	Content string    `gorm:"not null"`
}

type BlogService struct {
	db *gorm.DB
}

func (s BlogService) create(blog *Blog) error {
	return s.db.Table("blogs").Create(blog).Error
}

type blogPayload struct {
	Title   string `json:"name"`
	Content string `json:"content"`
	UUID    string `json:"uid"`
}

func (s BlogService) getAllBlogs() ([]*blogPayload, error) {
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

type BlogHandler struct {
	service *BlogService
}

func (h BlogHandler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	content := r.PostFormValue("content")
	title := r.PostFormValue("title")

	if content == "" {
		sendErrorResponse(w, &ErrorPayload{
			Message: "Content can not be null",
		}, http.StatusBadRequest)
	}
	id, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	blog := &Blog{
		UUID:    id,
		Title:   title,
		Content: content,
	}
	print(blog)
	dbErr := h.service.create(blog)
	if dbErr != nil {
		statusCode := http.StatusInternalServerError
		errorPayload := &ErrorPayload{
			Message: "Internal Error",
		}
		sendErrorResponse(w, errorPayload, statusCode)
		return
	}
	sendSuccessResponse(w, &map[string]interface{}{
		"message": "Save success",
		"blog":    blog.UUID,
	})
}

func (h BlogHandler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.service.getAllBlogs()
	if err != nil {
		sendErrorResponse(w, &ErrorPayload{
			Message: "Internal Error",
		}, 404)
	}
	sendSuccessResponse(w, &map[string]interface{}{
		"blogs": blogs,
	})
}
