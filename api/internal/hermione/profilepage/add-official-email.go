package profilepage

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/internal/hedwig"
	"github.com/psankar/vetchi/api/internal/middleware"
	"github.com/psankar/vetchi/api/internal/util"
	"github.com/psankar/vetchi/api/internal/wand"
	"github.com/psankar/vetchi/api/pkg/vetchi"
	"github.com/psankar/vetchi/typespec/hub"
)

func AddOfficialEmail(h wand.Wand) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req hub.AddOfficialEmailRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			h.Dbg("failed to decode request", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !h.Vator().Struct(w, &req) {
			h.Dbg("invalid request", "req", req)
			return
		}
		h.Dbg("validated", "req", req)

		hubUser, ok := r.Context().Value(middleware.HubUserCtxKey).(db.HubUserTO)
		if !ok {
			h.Err("failed to get hub user from context")
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		code := util.RandomString(vetchi.AddOfficialEmailCodeLenBytes)

		email, err := h.Hedwig().GenerateEmail(hedwig.GenerateEmailReq{
			TemplateName: hedwig.AddOfficialEmail,
			Args: map[string]string{
				"Name":   hubUser.FullName,
				"Handle": hubUser.Handle,
				"Code":   code,
			},
			EmailFrom: vetchi.EmailFrom,
			EmailTo:   []string{string(req.Email)},

			// TODO: The subject should be from Hedwig, based on the template
			// This subject is used in dolores/ too. Any change
			// in either place should be synced.
			Subject: "Vetchi - Confirm Email Ownership",
		})
		if err != nil {
			h.Dbg("failed to generate email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		err = h.DB().AddOfficialEmail(db.AddOfficialEmailReq{
			Email:   email,
			Code:    code,
			HubUser: hubUser,
			Context: r.Context(),
		})
		if err != nil {
			if errors.Is(err, db.ErrTooManyOfficialEmails) {
				h.Dbg("failed to add official email", "error", err)
				http.Error(w, "", http.StatusPreconditionFailed)
				return
			} else if errors.Is(err, db.ErrNoEmployer) {
				h.Dbg("failed to add official email", "error", err)
				http.Error(w, "", http.StatusUnprocessableEntity)
				return
			}

			h.Dbg("failed to add official email", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
