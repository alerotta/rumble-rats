package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/alerotta/rumble-rats/backend/utils"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register() http.Handler {
	return http.HandlerFunc(h.handleRegister)
}

func (h *Handler) Login() http.Handler {
	return http.HandlerFunc(h.handleLogin)
}

func(h *Handler) Validate() http.Handler {
	return http.HandlerFunc(h.handleValidation)
}

func (h* Handler) handleRegister (w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req RegisterRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.svc.Register(ctx, req)
	if err != nil {
		
		var ve ValidationError
		if errors.As(err, &ve) {
			utils.WriteError(w, http.StatusBadRequest,  ve.Msg)
			return
		}
		if errors.Is(err, ErrUserAlreadyExists) { 
			utils.WriteError(w, http.StatusConflict, "user already exists" )
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, resp)
}

func (h* Handler) handleLogin (w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req LoginRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	resp, err := h.svc.Login(ctx, req)
	if err != nil {
		var ve ValidationError
		if errors.As(err, &ve) {
			utils.WriteError(w, http.StatusBadRequest,  ve.Msg)
			return
		}
		if errors.Is(err, ErrUnauthorized) { 
			utils.WriteError(w, http.StatusUnauthorized, "invalid username or password" )
			return
		}

		utils.WriteError(w, http.StatusInternalServerError, "server error")
		return
	}

	utils.WriteJSON(w, http.StatusOK, resp)

}

func (h* Handler) handleValidation (w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
	utils.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	return
	}
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req VerifyTokenRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Token == "" {
		utils.WriteError(w, http.StatusBadRequest, "token required")
		return
	}

	if _, err := h.svc.ValidateAccessToken(req.Token); err != nil {
		utils.WriteError(w, http.StatusUnauthorized, "invalid token")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{
		"valid": true,
	})

}

