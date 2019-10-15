package coreservice

import (
	"github.com/jinzhu/gorm"
	"github.com/support/token"
	"github.com/getsentry/sentry-go"
	"net/http"
)

type SN struct {
	gorm.Model
	Value 		string `gorm:"unique_index"`
}

type snService struct {
	db			*gorm.DB
}

func (sn *snService) generateSN() *string{
	deviceSN, err := token.RandomToken(16)
	if err != nil {
		sentry.CaptureException(err)
		return nil
	}
	return deviceSN
}

type Device struct {
	gorm.Model
	SN 			string
}

type DeviceService struct {
	db 			*gorm.DB
}

type DeviceHandler struct {
	service		*DeviceService
}

func (dvs *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	email := r.PostFormValue("email")
	print(email)
}