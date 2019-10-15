package coreservice

import (
	"encoding/json"
	"github.com/getsentry/sentry-go"
	"github.com/jinzhu/gorm"
	"net/http"
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
		Blog{}, Device{}, SN{},	User{},
	}
	for _, model := range models {
		if err := db.AutoMigrate(model).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	snTableName := tx.NewScope(&SN{}).GetModelStruct().TableName(tx)
	print(snTableName)
	constrains := []struct {
		model     interface{}
		fieldName string
		refering  string
	}{
		{Device{}, "sn", snTableName + "(value)"},
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
	CreateBlog  	http.HandlerFunc
	GetAllBlogs 	http.HandlerFunc
	CreateDevice 	http.HandlerFunc
	CreateUser		http.HandlerFunc
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
	deviceService := &DeviceService{
		db:	db,
	}

	userService := &UserService {
		db : db,
	}

	userHandler := &UserHandler{
		service:	userService,
	}
	
	blogHandler := &BlogHandler{
		service:	 blogService,
	}
	deviceHandler := &DeviceHandler{
		service: 	deviceService,
	}

	return &Handler{
		CreateBlog:  	m.apply(blogHandler.CreateBlog, m.cors),
		GetAllBlogs: 	m.apply(blogHandler.GetAllBlogs, m.cors),
		CreateDevice:	m.apply(deviceHandler.CreateDevice, m.cors),
		CreateUser:		m.apply(userHandler.CreateUser, m.cors),
	}, nil

}
