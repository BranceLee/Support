package coreservice

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/gorm"
)

type errorPayload struct {
	Message string `json:"message"`
}

type response struct {
	StatusCode int                     `json:"status_code"`
	Error      *errorPayload           `json:"error"`
	Data       *map[string]interface{} `json:"data"`
}

func sendSuccessJSONResponse(w http.ResponseWriter, payload *response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(payload)
	if err != nil {
		sentry.CaptureException(err)
	}
}

func sendErrorResponse(w http.ResponseWriter, err *errorPayload, statusCode int) {
	sendSuccessJSONResponse(w, &response{
		StatusCode: statusCode,
		Error:      err,
		Data:       nil,
	})
}

func sendSuccessResponse(w http.ResponseWriter, payload *map[string]interface{}) {
	sendSuccessJSONResponse(w, &response{
		StatusCode: http.StatusOK,
		Error:      nil,
		Data:       payload,
	})
}

func dbMigrate(db *gorm.DB) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	models := []interface{}{
		Blog{}, User{}, Category{},
	}
	for _, model := range models {
		if err := db.AutoMigrate(model).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	categoryTableName := tx.NewScope(&Category{}).GetModelStruct().TableName(tx)
	constrains := []struct {
		model     interface{}
		fieldName string
		referring string
	}{
		{Blog{}, "category_id", categoryTableName + "(category_id)"},
	}

	// Add Foreignkey
	for _, c := range constrains {
		if err := tx.Model(c.model).AddForeignKey(c.fieldName, c.referring, "RESTRICT", "RESTRICT").Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// Handler contains all the http handlers for request endpoints
type Handler struct {
	CreateBlog     http.HandlerFunc
	GetAllBlogs    http.HandlerFunc
	CreateUser     http.HandlerFunc
	CreateCategory http.HandlerFunc
	GetCategory    http.HandlerFunc
}

// NewHandler returns an instance of the Handlers
func NewHandler(db *gorm.DB) (*Handler, error) {
	err := dbMigrate(db)
	if err != nil {
		return nil, err
	}
	m := middleware{}

	blogService := &blogService{
		db: db,
	}

	blogCategoryService := &blogCategoryService{
		db: db,
	}

	userService := &userService{
		db: db,
	}

	userHandl := &userHandler{
		service: userService,
	}

	blogHandl := &blogHandler{
		service: blogService,
	}

	categoryHandl := &categoryHandler{
		service: blogCategoryService,
	}

	return &Handler{
		CreateBlog:     m.apply(blogHandl.createBlog, m.cors),
		GetAllBlogs:    m.apply(blogHandl.getAllBlogs, m.configureSentry, m.cors),
		CreateUser:     m.apply(userHandl.createUser, m.cors),
		CreateCategory: m.apply(categoryHandl.createCategory, m.cors),
		GetCategory:    m.apply(categoryHandl.getCategory, m.cors),
	}, nil

}
