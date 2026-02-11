package server

import (
	"net/http"
	"strings"

	"github.com/alerotta/rumble-rats/backend/internal/auth"
	"github.com/alerotta/rumble-rats/backend/utils"
)

func RequireAuth(svc *auth.Service) func (http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
			h := r.Header.Get("Authorization")
			if h == "" {
				utils.WriteError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			const prefix = "Bearer "
			if !strings.HasPrefix(h, prefix) {
				utils.WriteError(w, http.StatusUnauthorized, "invalid authorization header")
				return
			}

			token := strings.TrimSpace(strings.TrimPrefix(h, prefix))
			if token == "" {
				utils.WriteError(w, http.StatusUnauthorized, "token required")
				return
			}

			if _, err := svc.ValidateAccessToken(token); err != nil {
				utils.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}