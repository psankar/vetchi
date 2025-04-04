package candidacy

import (
	"encoding/json"
	"net/http"

	"github.com/vetchium/vetchium/api/internal/db"
	"github.com/vetchium/vetchium/api/internal/wand"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func EmployerAddComment(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EmployerAddComment")
		var addCommentReq employer.AddEmployerCandidacyCommentRequest
		err := json.NewDecoder(r.Body).Decode(&addCommentReq)
		if err != nil {
			h.Dbg("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &addCommentReq) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("validated", "addCommentReq", addCommentReq)

		commentID, err := h.DB().
			AddEmployerCandidacyComment(r.Context(), addCommentReq)
		if err != nil {
			switch err {
			case db.ErrNoOpening:
				h.Dbg("Candidacy not found", "error", err)
				http.Error(w, "Candidacy not found", http.StatusNotFound)
			case db.ErrInvalidCandidacyState:
				h.Dbg("Invalid candidacy state", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
			case db.ErrUnauthorizedComment:
				h.Dbg("User not authorized to comment", "error", err)
				http.Error(w, "", http.StatusForbidden)
			default:
				h.Dbg("Internal error while adding comment", "error", err)
				http.Error(w, "", http.StatusInternalServerError)
			}
			return
		}

		h.Dbg("Added comment", "commentID", commentID)
		w.WriteHeader(http.StatusOK)
	}
}

func EmployerGetComments(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.Dbg("Entered EmployerGetComments")
		var getCommentsReq common.GetCandidacyCommentsRequest
		err := json.NewDecoder(r.Body).Decode(&getCommentsReq)
		if err != nil {
			h.Dbg("Error decoding request body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &getCommentsReq) {
			h.Dbg("Validation failed")
			return
		}
		h.Dbg("validated", "getCommentsReq", getCommentsReq)

		comments, err := h.DB().
			GetEmployerCandidacyComments(r.Context(), getCommentsReq)
		if err != nil {
			h.Dbg("Internal error while getting comments", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		h.Dbg("Got comments", "comments", comments)
		err = json.NewEncoder(w).Encode(comments)
		if err != nil {
			h.Err("Error encoding response", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
}
