package openings

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/employer"
)

func GetOpening(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetOpening")
		var getOpeningReq employer.GetOpeningRequest
		err := json.NewDecoder(r.Body).Decode(&getOpeningReq)
		if err != nil {
			h.Dbg("failed to decode get opening request", "error", err)
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getOpeningReq) {
			h.Dbg("validation failed", "getOpeningReq", getOpeningReq)
			return
		}
		h.Dbg("validated", "getOpeningReq", getOpeningReq)

		opening, err := h.DB().GetOpening(r.Context(), getOpeningReq)
		if err != nil {
			if errors.Is(err, db.ErrNoOpening) {
				h.Dbg("opening not found", "id", getOpeningReq.ID)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to get opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got opening", "opening", opening)
		err = json.NewEncoder(w).Encode(opening)
		if err != nil {
			h.Err("failed to encode opening", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
