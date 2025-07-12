package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/black-spidera/student-api/internal/storage"
	"github/black-spidera/student-api/internal/types"
	"github/black-spidera/student-api/internal/utils"
	"io"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		student := &types.Student{}

		err := json.NewDecoder(r.Body).Decode(student)
		if errors.Is(err, io.EOF) {
			utils.WriteJSONResponse(w, http.StatusBadRequest, "Request body cannot be empty")
			return
		} else if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		err = validator.New().Struct(student)
		if err != nil {
			validatorError := err.(validator.ValidationErrors)
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ValidationErrorFormat(validatorError))
			return
		}

		id, err := storage.CreateStudent(student.Name, student.Email, student.Age)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create student: %v", err))
			return
		}

		utils.WriteJSONResponse(w, http.StatusCreated, map[string]string{fmt.Sprintf("student_%d", id): "created"})

	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")

		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, "Invalid student ID")
			return
		}

		student, err := storage.GetStudentById(id)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve student: %v", err))
			return
		}

		if student.Id == 0 {
			utils.WriteJSONResponse(w, http.StatusNotFound, "Student not found")
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, student)
	}
}
