package user

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"go-service.codymj.io/cmd/app/util"
	"net/http"
	"strconv"
)

// getById handles requests to GET /users/:id.
func (h *handler) getById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Parse ID from request.
	idParam := p.ByName("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		util.WriteErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	// Call service to get user by ID.
	res, err := h.services.UserService.GetById(r.Context(), id)
	if err != nil {
		util.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write response.
	if res.Id == 0 {
		// No user found.
		w.WriteHeader(http.StatusNoContent)
		_ = json.NewEncoder(w).Encode(nil)
		return
	}
	bytes, _ := json.Marshal(res)
	w.Header().Set(util.ContentType, util.JsonHeader)
	_, _ = w.Write(bytes)
}
