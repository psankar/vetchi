package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

type orgUserCtx int

const (
	OrgUserCtxKey orgUserCtx = iota
)

type Middleware struct {
	db  db.DB
	log *slog.Logger
}

func NewMiddleware(db db.DB, log *slog.Logger) *Middleware {
	return &Middleware{db: db, log: log}
}

func (m *Middleware) employerAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}

		orgUser, err := m.db.AuthOrgUser(r.Context(), authHeader)
		if err != nil {
			if errors.Is(err, db.ErrNoOrgUser) {
				http.Error(w, "", http.StatusUnauthorized)
				return
			}

			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		m.log.Debug("Authenticated org user", "orgUser", orgUser)

		ctx := context.WithValue(r.Context(), OrgUserCtxKey, orgUser)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Protect middleware only checks with the roles of the OrgUser. Whether the
// OrgUser belongs to the Org or not, should be verified via EmployerAuth. We
// should not take the OrgID/EmployerID on any of the request bodies but should
// be derived from the OrgUserToken.
func (m *Middleware) Protect(
	route string,
	handler http.Handler,
	allowedRoles []vetchi.OrgUserRole,
) {
	http.Handle(route, m.employerAuth(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			orgUser, ok := ctx.Value(OrgUserCtxKey).(db.OrgUser)
			if !ok {
				m.log.Error("Failed to get orgUser from context")
				http.Error(w, "", http.StatusInternalServerError)
				return
			}

			if !hasRoles(orgUser.OrgUserRoles, allowedRoles) {
				http.Error(w, "", http.StatusForbidden)
				return
			}

			// Call the actual handler if roles are sufficient
			handler.ServeHTTP(w, r)
		}),
	))
}

func hasRoles(
	orgUserRoles []vetchi.OrgUserRole,
	allowedRoles []vetchi.OrgUserRole,
) bool {
	// TODO: Can potentially cache this ifF there is a performance issue
	for _, orgUserRole := range orgUserRoles {
		for _, role := range allowedRoles {
			if role == orgUserRole {
				return true
			}
		}
	}

	return false
}
