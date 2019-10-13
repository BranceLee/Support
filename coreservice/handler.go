package coreservice

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	u "github.com/support/utils"
	"net/http"
	"fmt"
	"os"
)

type ErrorPayload struct {
	Message			string		`json:"message"`
}

type Response struct {
	StatusCode		int							`json:"status_code"`
	Error			*ErrorPayload				`json:"error"`
	Data			*map[string]interface{}		`json:"data"`
}

func sendSuccessJSONResponse(w http.ResponseWriter, payload *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(payload)
	if err != nil {
		fmt.Println("error")
	}
}

func sendErrorResponse(w http.ResponseWriter, err *ErrorPayload, statusCode int) {
	sendSuccessJSONResponse(w, &Response{
		StatusCode: statusCode,
		Error:		err,
		Data:		nil,
	})
}

func sendSuccessResponse(w http.ResponseWriter, payload *map[string]interface{}){
	sendSuccessJSONResponse(w, &Response{
		StatusCode: http.StatusOK,
		Error:		nil,
		Data:		payload,
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
	service *Service
}

func (h Handler) CreateBlog(w http.ResponseWriter, r *http.Request) {
	Cors(w)
	content := r.PostFormValue("content")
	title := r.PostFormValue("title")

	if content == "" {
		sendErrorResponse(w, &ErrorPayload{
			Message:	"Content can not be null",
		}, 404)
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

func (h Handler) GetAllBlogs(w http.ResponseWriter, r *http.Request){
	Cors(w)
	blogs, err := h.service.getAllBlogs()
	if err != nil {
		sendErrorResponse(w, &ErrorPayload{
			Message:	"Internal Error",
		}, 404)
	}
	fmt.Println(blogs)
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
