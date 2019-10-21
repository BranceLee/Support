package coreservice

import (
	"errors"
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

const (
	errEmailInvalid = "Invalid Email"
	errEmailTaken   = "Email has been Taken"
	errInternal     = "Internal error"
	errGenericError = "An error occurred. Please try again later."
)

type User struct {
	gorm.Model
	UUID  uuid.UUID `gorm:"unique_index; not null" sql:"uuid"`
	Email string    `gorm:"not null; unique_index"`
}

type UserService struct {
	db             *gorm.DB
	passwordPepper string
}

type UserHandler struct {
	service *UserService
}

func runValidator(user *User, fns ...func(*User) error) error {
	for _, fn := range fns {
		if err := fn(user); err != nil {
			return err
		}
	}
	return nil
}

func (s *UserService) create(user *User) error {
	err := runValidator(user, s.requireUUID, s.requireUniqueEmail)
	if err != nil {
		return err
	}
	err = s.db.Create(user).Error
	if err != nil {
		sentry.CaptureException(err)
		return err
	}
	return nil
}

func (s *UserService) byEmail(email string) (*User, *ModelError) {
	identity := &User{}
	err := s.db.Where(&User{Email: email}).First(identity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			e := &ModelError{
				Kind: ErrTypeNotFound,
				Err:  err,
			}
			return nil, e
		}
		e := &ModelError{
			Kind: ErrTypeDBError,
			Err:  err,
		}
		sentry.CaptureException(err)
		return nil, e
	}
	return identity, nil
}

func (s *UserService) requireUniqueEmail(user *User) error {
	var count int
	err := s.db.Model(&User{}).Where("email = ?", user.Email).Count(&count).Error
	if err != nil {
		sentry.CaptureException(err)
		return errors.New(errGenericError)
	}
	if count != 0 {
		return errors.New(errEmailTaken)
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
			Message: errEmailInvalid,
		}
		sendErrorResponse(w, errorPayload, http.StatusBadRequest)
	}
	newUser := &User{Email: email}
	dbErr := h.service.create(newUser)
	if dbErr != nil {
		if dbErr.Error() == errEmailTaken {
			errorPayload := &ErrorPayload{
				Message: errEmailTaken,
			}
			statusCode := http.StatusBadRequest
			sendErrorResponse(w, errorPayload, statusCode)
			return
		}
		errorPayload := &ErrorPayload{
			Message: errInternal,
		}
		statusCode := http.StatusInternalServerError
		sendErrorResponse(w, errorPayload, statusCode)
		return
	}
	userPayload := &map[string]interface{}{
		"email": newUser.Email,
		"uid":   newUser.UUID,
	}
	sendSuccessResponse(w, userPayload)
}
