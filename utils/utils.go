package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{
		"status":  status,
		"message": message,
	}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

type Response struct {
	StatusCode uint                   `json:"status_code"`
	Error      ErrorPayload           `json:"error"`
	Data       map[string]interface{} `json:"data"`
}

type ErrorPayload struct {
	Message string `json:"message"`
}

func sendResponse(w http.ResponseWriter, payload Response) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		fmt.Println("Parse Json error", err)
	}
}

func SendSuccessResponse(w http.ResponseWriter, payload map[string]interface{}) {
	res := Response{
		StatusCode: http.StatusOK,
		Error:      ErrorPayload{},
		Data:       payload,
	}
	sendResponse(w, res)
}

func SendErrorResponse(w http.ResponseWriter, err error, code uint) {
	res := Response{
		StatusCode: code,
		Error:      ErrorPayload{Message: err.Error()},
		Data:       map[string]interface{}{},
	}
	sendResponse(w, res)
}
