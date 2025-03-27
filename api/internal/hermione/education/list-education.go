package education

import (
	"encoding/json"
	"net/http"

	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/hub"
)

func ListEducation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var listEducationReq hub.ListEducationRequest
		err := json.NewDecoder(r.Body).Decode(&listEducationReq)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &listEducationReq) {
			h.Dbg("invalid request", "listEducationReq", listEducationReq)
			return
		}

		h.Dbg("validated", "listEducationReq", listEducationReq)

		// There may be no plural educations grammatically, but it makes it easier to understand
		educations, err := h.DB().ListEducation(r.Context(), listEducationReq)
		if err != nil {
			h.Dbg("failed to list education", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(educations)
		if err != nil {
			h.Dbg("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
