package user

import (
	"encoding/json"
	"go-service.codymj.io/cmd/app/util"
	"net/http"
	"strings"
)

// getByParams handles requests to GET /users.
func (h *handler) getByParams(w http.ResponseWriter, r *http.Request) {
	// Parse params from request.
	u, _ := r.URL.Parse(r.URL.String())
	params := make(map[string]string)
	if !strings.EqualFold("", u.RawQuery) {
		params = util.ParseQueryString(u.RawQuery)
	}

	// Call service to get users.
	res, err := h.services.UserService.GetByParams(r.Context(), params)
	if err != nil {
		util.WriteErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	// Write response.
	if len(res) == 0 {
		// No users found.
		w.WriteHeader(http.StatusNoContent)
		_ = json.NewEncoder(w).Encode(nil)
		return
	}
	bytes, _ := json.Marshal(res)
	w.Header().Set(util.ContentType, util.JsonHeader)
	_, _ = w.Write(bytes)
}
