package locations

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/typespec/employer"
)

func GetLocation(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered GetLocation")
		var getLocationReq employer.GetLocationRequest
		err := json.NewDecoder(r.Body).Decode(&getLocationReq)
		if err != nil {
			h.Dbg("failed to decode get location request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getLocationReq) {
			h.Dbg("validation failed", "getLocationReq", getLocationReq)
			return
		}
		h.Dbg("validated", "getLocationReq", getLocationReq)

		location, err := h.DB().GetLocByName(r.Context(), getLocationReq)
		if err != nil {
			if errors.Is(err, db.ErrNoLocation) {
				h.Dbg("location not found", "title", getLocationReq.Title)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			h.Dbg("failed to get location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("got location", "location", location)
		err = json.NewEncoder(w).Encode(location)
		if err != nil {
			h.Err("failed to encode location", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
