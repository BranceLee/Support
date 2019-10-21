package coreservice

import (
	"encoding/json"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/gorm"
)

type ErrorPayload struct {
	Message string `json:"message"`
}

type Response struct {
	StatusCode int                     `json:"status_code"`
	Error      *ErrorPayload           `json:"error"`
	Data       *map[string]interface{} `json:"data"`
}

func sendSuccessJSONResponse(w http.ResponseWriter, payload *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(payload)
	if err != nil {
		sentry.CaptureException(err)
	}
}

func sendErrorResponse(w http.ResponseWriter, err *ErrorPayload, statusCode int) {
	sendSuccessJSONResponse(w, &Response{
		StatusCode: statusCode,
		Error:      err,
		Data:       nil,
	})
}

func sendSuccessResponse(w http.ResponseWriter, payload *map[string]interface{}) {
	sendSuccessJSONResponse(w, &Response{
		StatusCode: http.StatusOK,
		Error:      nil,
		Data:       payload,
	})
}

func dbMigrate(db *gorm.DB) error {
	tx := db.Begin()

	//Close transaction.
	defer tx.Rollback()
	models := []interface{}{
		Blog{}, User{}, BlogCategory{},
	}
	for _, model := range models {
		if err := db.AutoMigrate(model).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	blCategoryTableName := tx.NewScope(&BlogCategory{}).GetModelStruct().TableName(tx)

	constrains := []struct {
		model     interface{}
		fieldName string
		refering  string
	}{
		{Blog{}, "category_id", blCategoryTableName + "(category_id)"},
	}

	// Add Foreignkey
	for _, c := range constrains {
		if err := tx.Model(c.model).AddForeignKey(c.fieldName, c.refering, "RESTRICT", "RESTRICT").Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

type Handler struct {
	CreateBlog         http.HandlerFunc
	GetAllBlogs        http.HandlerFunc
	CreateUser         http.HandlerFunc
	CreateBlogCategory http.HandlerFunc
}

func NewHandler(db *gorm.DB) (*Handler, error) {
	err := dbMigrate(db)
	if err != nil {
		return nil, err
	}
	m := middleware{}

	blogService := &BlogService{
		db: db,
	}

	blogCategoryService := &BlogCategoryService{
		db: db,
	}

	userService := &UserService{
		db: db,
	}

	userHandl := &UserHandler{
		service: userService,
	}

	blogHandl := &BlogHandler{
		service: blogService,
	}

	blogCategoryHandl := &BlogCategoryHandler{
		service: blogCategoryService,
	}

	return &Handler{
		CreateBlog:         m.apply(blogHandl.CreateBlog, m.cors),
		GetAllBlogs:        m.apply(blogHandl.GetAllBlogs, m.cors),
		CreateUser:         m.apply(userHandl.CreateUser, m.cors),
		CreateBlogCategory: m.apply(blogCategoryHandl.CreateBlogCategory, m.cors),
	}, nil

}
