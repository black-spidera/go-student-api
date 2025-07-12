package students

import (
	"encoding/json"
	"errors"
	"fmt"
	"github/black-spidera/student-api/internal/types"
	"github/black-spidera/student-api/internal/utils"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		person := &types.Person{}

		err := json.NewDecoder(r.Body).Decode(person)
		if errors.Is(err, io.EOF) {
			utils.WriteJSONResponse(w, http.StatusBadRequest, "Request body cannot be empty")
			return
		} else if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		err = validator.New().Struct(person)
		if err != nil {
			validatorError := err.(validator.ValidationErrors)
			utils.WriteJSONResponse(w, http.StatusBadRequest, utils.ValidationErrorFormat(validatorError))
			return
		}

		utils.WriteJSONResponse(w, http.StatusCreated, fmt.Sprintf("Student %s created successfully", person.Name))

	}
}
