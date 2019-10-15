package coreservice

import (
	"github.com/jinzhu/gorm"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"net/http"
)

const (
	errEmailInvalid = "Invalid Email"
)

type User struct {
	gorm.Model
	UUID			uuid.UUID		`gorm:"unique_index; not null" sql:"uuid"`
	Email			string			`gorm:"not null; unique_index"`
}

type UserService struct {
	db					*gorm.DB
	passwordPepper		string
}

type UserHandler struct {
	service		*UserService
}

func (s *UserService) create(user *User) error {
	s.requireUUID(user)
	err := s.db.Create(user).Error
	if err != nil {
		sentry.CaptureException(err)
		return err 
	}
	return nil
}

func (s *UserService) requireUUID(user *User) error {
	uid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	user.UUID = uid
	return nil
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.PostFormValue("email")
	if email == "" {
		errorPayload := &ErrorPayload{
			Message:	errEmailInvalid,
		}
		sendErrorResponse(w, errorPayload, http.StatusBadRequest)
	}
	newUser := &User{Email:email,}
	err := h.service.create(newUser)
	if err != nil {
		errorPayload := &ErrorPayload{
			Message:	"Internal error",
		}
		statusCode := http.StatusInternalServerError
		sendErrorResponse(w, errorPayload, statusCode)
	}
	userPayload := &map[string]interface{}{
		"email":		newUser.Email,
		"uid":			newUser.UUID,
	}
	sendSuccessResponse(w,userPayload)
}