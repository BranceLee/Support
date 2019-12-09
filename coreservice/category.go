package coreservice

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/gorm"
	"github.com/BranceLee/Support/token"
)

// Category contains all the category of the blog
type Category struct {
	gorm.Model
	CategoryID string `gorm:"not null;unique_index"`
	Name       string `gorm:"not null"`
}

type blogCategoryService struct {
	db *gorm.DB
}

func (s *blogCategoryService) create(blogCategory *Category) error {
	return s.db.Create(blogCategory).Error
}

func (s *blogCategoryService) getCategoryByName(name string) (*Category, *ModelError) {
	if name == "" {
		return nil, &ModelError{
			Kind: ErrTypeValidation,
			Err:  errors.New("Empty uuid"),
		}
	}

	category := &Category{
		Name: name,
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
	Name       string `json:"name"`
	CategoryID string `json:"category_id"`
}

func (s *blogCategoryService) getCategory() ([]*categoryPayload, error) {
	result := []*categoryPayload{}
	rows, err := s.db.Model(&Category{}).Select(`name, category_id`).Rows()
	if err != nil {
		sentry.CaptureException(err)
		return nil, err
	}

	for rows.Next() {
		var name string
		var uid string
		if err := rows.Scan(&name, &uid); err != nil {
			sentry.CaptureException(err)
			return nil, err
		}
		result = append(result, &categoryPayload{
			Name:       name,
			CategoryID: uid,
		})
	}

	return result, nil
}

type categoryHandler struct {
	service *blogCategoryService
}

func (ch *categoryHandler) createCategory(w http.ResponseWriter, r *http.Request) {
	categoryName := r.PostFormValue("category")

	if categoryName == "" {
		sendErrorResponse(w, &errorPayload{
			Message: "Category can not be null",
		}, http.StatusBadRequest)
		return
	}

	_, err := ch.service.getCategoryByName(categoryName)

	if err != nil && err.Kind == ErrTypeDBError {
		sendErrorResponse(w, &errorPayload{
			Message: "Internal Error",
		}, http.StatusBadRequest)
		return
	}

	if err != nil && err.Kind == ErrTypeNotFound {
		id, uidErr := token.RandomToken(16)
		if uidErr != nil {
			sendErrorResponse(w, &errorPayload{
				Message: "Internal Error",
			}, http.StatusInternalServerError)
			return
		}

		category := &Category{
			CategoryID: *id,
			Name:       categoryName,
		}
		dbErr := ch.service.create(category)
		if dbErr != nil {
			sentry.CaptureException(dbErr)
			sendErrorResponse(w, &errorPayload{
				Message: "Internal Error",
			}, http.StatusInternalServerError)
			return
		}

		sendSuccessResponse(w, &map[string]interface{}{
			"uid":  category.CategoryID,
			"name": category.Name,
		})
		return
	}

	sendErrorResponse(w, &errorPayload{
		Message: "Category has been taken",
	}, http.StatusBadRequest)
}

func (ch *categoryHandler) getCategory(w http.ResponseWriter, r *http.Request) {
	category, err := ch.service.getCategory()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sendSuccessResponse(w, &map[string]interface{}{
		"category": category,
	})
}
