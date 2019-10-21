package coreservice

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type Blog struct {
	gorm.Model
	UUID       uuid.UUID `gorm:"unique_index; not null" sql:"type:uuid"`
	Title      string    `gorm:"not null"`
	Content    string    `gorm:"not null"`
	CategoryID string    `gorm:"not null"`
}

type BlogService struct {
	db *gorm.DB
}

func (s BlogService) create(blog *Blog) error {
	return s.db.Table("blogs").Create(blog).Error
}

type blogPayload struct {
	Title        string `json:"name"`
	Content      string `json:"content"`
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
}

func (s BlogService) getAllBlogs() ([]*blogPayload, error) {
	result := []*blogPayload{}
	rows, err := s.db.Model(&Blog{}).Select(`blogs.title, blogs.content, blogs.category_id,blog_categories.name`).Joins("left join blog_categories on blog_categories.category_id = blogs.category_id").Order("blogs.updated_at DESC").Rows()
	if err != nil {
		fmt.Println("error is ", err)
		return nil, nil
	}

	for rows.Next() {
		var title string
		var content string
		var categoryName string
		var categoryID string
		if err := rows.Scan(&title, &content, &categoryID, &categoryName); err != nil {
			fmt.Println("scan err: ", err)
		}
		result = append(result, &blogPayload{
			Title:        title,
			Content:      content,
			CategoryID:   categoryID,
			CategoryName: categoryName,
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
	categoryID := r.PostFormValue("category_id")
	if content == "" || title == "" || categoryID == "" {
		sendErrorResponse(w, &ErrorPayload{
			Message: "Bad Request",
		}, http.StatusBadRequest)
		return
	}

	id, err := uuid.NewRandom()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	blog := &Blog{
		UUID:       id,
		Title:      title,
		Content:    content,
		CategoryID: categoryID,
	}
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
		"blog": map[string]interface{}{
			"uid":         blog.UUID,
			"category_id": blog.CategoryID,
		},
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
