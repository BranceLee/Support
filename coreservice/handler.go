package coreservice

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	u "github.com/support/utils"
	"net/http"
	"os"
)

func response(statusCode int, payload string) map[string]interface{} {
	return map[string]interface{}{
		"statusCode": statusCode,
		"data":       payload,
	}
}

func sendErrorResponse(w http.ResponseWriter, err string, statusCode int) {
	payload := response(statusCode, err)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}

func Core(w http.ResponseWriter) {
	env, _ := os.LookupEnv("ENV")
	if env == "DEV" {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	}
}

type Handler struct {
	service *Service
}

func (h Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	Core(w)
	content := r.PostFormValue("content")
	title := r.PostFormValue("title")

	if content == "" {
		sendErrorResponse(w, "Content can not be null", 404)
		return
	}

	blog := &Blog{
		Title:   title,
		Content: content,
	}
	print(blog)
	err := h.service.create(blog)
	if err != nil {
		res := u.Message(false, "Internal Error")
		u.Respond(w, res)
		return
	}
	res := u.Message(true, "Success")
	u.Respond(w, res)
}

func NewHandler(db *gorm.DB) (*Handler, error) {
	db.AutoMigrate(Blog{})
	service := &Service{
		db: db,
	}
	return &Handler{
		service: service,
	}, nil
}
