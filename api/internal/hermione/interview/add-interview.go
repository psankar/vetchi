package interview

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/employer"
)

func AddInterview(h wand.Wand) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered AddInterview")
		var addInterviewReq employer.AddInterviewRequest
		if err := json.NewDecoder(r.Body).Decode(&addInterviewReq); err != nil {
			h.Dbg("decoding failed", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addInterviewReq) {
			h.Dbg("validation failed", "addInterviewReq", addInterviewReq)
			return
		}
		h.Dbg("validated", "addInterviewReq", addInterviewReq)

		interviewID := util.RandomUniqueID(vetchi.InterviewIDLenBytes)

		err := h.DB().AddInterview(r.Context(), db.AddInterviewRequest{
			AddInterviewRequest: addInterviewReq,
			InterviewID:         interviewID,
		})
		if err != nil {
			if errors.Is(err, db.ErrNoCandidacy) {
				h.Dbg("no candidacy found", "error", err)
				http.Error(w, "", http.StatusNotFound)
				return
			}

			if errors.Is(err, db.ErrInvalidCandidacyState) {
				h.Dbg("candidacy not in valid state", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to add interview", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("added interview", "interviewID", interviewID)
		err = json.NewEncoder(w).Encode(employer.AddInterviewResponse{
			InterviewID: interviewID,
		})
		if err != nil {
			h.Err("failed to encode response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	})
}
