package coreservice

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/gorm"
	"github.com/support/token"
)

type BlogCategory struct {
	gorm.Model
	CategoryID string `gorm:"not null;unique_index"`
	Name       string `gorm:"not null"`
}

type BlogCategoryService struct {
	db *gorm.DB
}

func (s *BlogCategoryService) create(blogCategory *BlogCategory) error {
	return s.db.Table("blog_categories").Create(blogCategory).Error
}

func (s *BlogCategoryService) getCategoryByName(ctName string) (*BlogCategory, *ModelError) {
	if ctName == "" {
		return nil, &ModelError{
			Kind: ErrTypeValidation,
			Err:  errors.New("Empty uuid"),
		}
	}

	category := &BlogCategory{
		Name: ctName,
	}
	err := s.db.Where("name = ?", category.Name).Take(category).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &ModelError{
				Err:  err,
				Kind: ErrTypeNotFound,
			}
		}
		return nil, &ModelError{
			Err:  err,
			Kind: ErrTypeDBError,
		}
	}

	return category, nil
}

type categoryPayload struct {
	CategoryID string `json:"uid"`
	Name       string `json:"name"`
}

func (s *BlogCategoryService) getBlogCategory() ([]*categoryPayload, error) {
	result := []*categoryPayload{}
	rows, err := s.db.Model(&BlogCategory{}).Select(`name, category_id`).Rows()
	if err != nil {
		fmt.Println("error is ", err)
		return nil, nil
	}

	for rows.Next() {
		var name string
		var uid string
		if err := rows.Scan(&name, &uid); err != nil {
			fmt.Println("scan err: ", err)
		}
		result = append(result, &categoryPayload{
			CategoryID: uid,
			Name:       name,
		})
	}

	return result, nil
}

type BlogCategoryHandler struct {
	service *BlogCategoryService
}

func (bl *BlogCategoryHandler) CreateBlogCategory(w http.ResponseWriter, r *http.Request) {
	categoryName := r.PostFormValue("category")

	if categoryName == "" {
		sendErrorResponse(w, &ErrorPayload{
			Message: "Category can not be null",
		}, http.StatusBadRequest)
		return
	}

	_, err := bl.service.getCategoryByName(categoryName)

	if err != nil && err.Kind == ErrTypeDBError {
		sendErrorResponse(w, &ErrorPayload{
			Message: "Internal Error",
		}, http.StatusBadRequest)
		return
	}

	if err != nil && err.Kind == ErrTypeNotFound {
		id, uidErr := token.RandomToken(16)
		if uidErr != nil {
			sendErrorResponse(w, &ErrorPayload{
				Message: "Internal Error",
			}, http.StatusInternalServerError)
			return
		}

		category := &BlogCategory{
			CategoryID: *id,
			Name:       categoryName,
		}

		dbErr := bl.service.create(category)
		if dbErr != nil {
			statusCode := http.StatusInternalServerError
			errorPayload := &ErrorPayload{
				Message: "Internal Error",
			}
			sentry.CaptureException(dbErr)
			sendErrorResponse(w, errorPayload, statusCode)
			return
		}
		sendSuccessResponse(w, &map[string]interface{}{
			"message": "Save success",
			"uid":     category.CategoryID,
			"name":    category.Name,
		})
		return
	}

	sendErrorResponse(w, &ErrorPayload{
		Message: "Category has been taken",
	}, http.StatusBadRequest)
}