package utils

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func WriteJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Status:  status,
		Message: http.StatusText(status),
		Data:    data,
	}
	return json.NewEncoder(w).Encode(response)
}

func ValidationErrorFormat(err validator.ValidationErrors) string {
	var errorsSlice []string
	for _, fieldError := range err {
		switch fieldError.ActualTag() {
		case "required":
			errorsSlice = append(errorsSlice, fieldError.Field()+" is required")
		default:
			errorsSlice = append(errorsSlice, fieldError.Field()+" is invalid")
		}
	}
	return strings.Join(errorsSlice, ", ")
}
