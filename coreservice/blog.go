package coreservice

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Blog contains all the blog information
type Blog struct {
	gorm.Model
	UUID       uuid.UUID `gorm:"unique_index; not null" sql:"type:uuid"`
	Title      string    `gorm:"not null"`
	Content    string    `gorm:"not null"`
	CategoryID string    `gorm:"not null"`
}

type blogService struct {
	db *gorm.DB
}

func (s *blogService) create(blog *Blog) error {
	return s.db.Table("blogs").Create(blog).Error
}

type blogPayload struct {
	Title        string `json:"name"`
	Content      string `json:"content"`
	CategoryID   string `json:"category_id"`
	CategoryName string `json:"category_name"`
}

func (s *blogService) getAllBlogs() ([]*blogPayload, error) {
	result := []*blogPayload{}
	rows, err := s.db.Model(&Blog{}).Select(`blogs.title, blogs.content, blogs.category_id,blog_categories.name`).Joins("left join categories on categories.category_id = blogs.category_id").Order("blogs.updated_at DESC").Rows()
	if err != nil {
		sentry.CaptureException(err)
		return nil, nil
	}

	for rows.Next() {
		var title string
		var content string
		var categoryName string
		var categoryID string
		if err := rows.Scan(&title, &content, &categoryID, &categoryName); err != nil {
			sentry.CaptureException(err)
			return nil, err
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

type blogHandler struct {
	service *blogService
}

func (h blogHandler) createBlog(w http.ResponseWriter, r *http.Request) {
	content := r.PostFormValue("content")
	title := r.PostFormValue("title")
	categoryID := r.PostFormValue("category_id")
	if content == "" || title == "" || categoryID == "" {
		sendErrorResponse(w, &errorPayload{
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
		errorPayload := &errorPayload{
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

func (h *blogHandler) getAllBlogs(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.service.getAllBlogs()
	if err != nil {
		sendErrorResponse(w, &errorPayload{
			Message: "Internal Error",
		}, 404)
	}
	sendSuccessResponse(w, &map[string]interface{}{
		"blogs": blogs,
	})
}
