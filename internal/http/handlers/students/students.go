package students

import (
	"encoding/json"
	"github/black-spidera/student-api/internal/types"
	"github/black-spidera/student-api/internal/utils"
	"net/http"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		person := &types.Person{}

		if err := json.NewDecoder(r.Body).Decode(person); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

	}
}
