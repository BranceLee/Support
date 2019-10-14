package coreservice

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/getsentry/sentry-go"
	"net/http"
	"os"
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

func Cors(w http.ResponseWriter) {
	env, _ := os.LookupEnv("ENV")
	if env == "DEV" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	}
}

type Handler struct {
	service *BlogService
}

func (h Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	Cors(w)
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
		"blog":		blog.UUID,
	})
}

func (h Handler) GetAllBlogs(w http.ResponseWriter, r *http.Request) {
	Cors(w)
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

func dbMigrate(db *gorm.DB) error{
	tx := db.Begin()

	//Close transaction.
	defer tx.Rollback()

	userTableName := tx.NewScope(&Blog{}).GetModelStruct().TableName(tx)
	print(userTableName)	
	models := []interface{}{
		Blog{}, Device{}, SN{},
	}
	for _, model := range models {
		if err := db.AutoMigrate(model).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func NewHandler(db *gorm.DB) (*Handler, error) {
	err := dbMigrate(db)
	if err != nil {
		return nil, err
	}
	service := &BlogService{
		db: db,
	}
	return &Handler{
		service: service,
	}, nil
}
