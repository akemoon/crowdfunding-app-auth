package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/akemoon/crowdfunding-app-auth/metrics"
	"github.com/akemoon/crowdfunding-app-auth/service/auth"
	"github.com/akemoon/crowdfunding-app-auth/service/token"
)

// @Summary Sign up
// @Description Create credentials and user profile
// @Accept json
// @Produce json
// @Param payload body domain.SignUpRequest true "Sign up payload"
// @Success 201 "User created"
// @Failure 400 "Invalid request"
// @Failure 405 "Method not allowed"
// @Failure 409 "Conflict"
// @Failure 500 "Internal server error"
// @Router /signup [post]
func SignUp(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req domain.SignUpRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.SignUp(r.Context(), req)
		if err != nil {
			log.Printf("service: %s", err.Error())

			status, resp := mapErrToHTTP(err)
			writeJSON(w, status, resp)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// @Summary Sign in
// @Description Authenticate and return access/refresh tokens
// @Accept json
// @Produce json
// @Param payload body domain.SignInRequest true "Sign in payload"
// @Success 200 {object} domain.SignInResponse "Tokens issued"
// @Failure 400 "Invalid request"
// @Failure 405 "Method not allowed"
// @Failure 500 "Internal server error"
// @Router /signin [post]
func SignIn(svc *auth.Service, m *metrics.AuthMetrics) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req domain.SignInRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		resp, err := svc.SignIn(r.Context(), req)
		if err != nil {
			log.Printf("auth service: %s", err)

			m.AuthSignInTotal.WithLabelValues("failure").Inc()

			status, errResp := mapErrToHTTP(err)
			writeJSON(w, status, errResp)
			return
		}

		m.AuthSignInTotal.WithLabelValues("success").Inc()

		writeJSON(w, http.StatusOK, resp)
	}
}

// @Summary Sign out
// @Description Revoke refresh token
// @Accept json
// @Produce json
// @Param payload body domain.SignOutRequest true "Sign out payload"
// @Success 200 "Signed out"
// @Failure 400 "Invalid request"
// @Failure 405 "Method not allowed"
// @Failure 500 "Internal server error"
// @Router /signout [post]
func SignOut(svc *auth.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req domain.SignOutRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		err = svc.SignOut(r.Context(), req.RefreshToken)
		if err != nil {
			log.Printf("auth service: %s", err)

			status, resp := mapErrToHTTP(err)
			writeJSON(w, status, resp)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

// @Summary Check access token
// @Description Validate access token from Authorization header
// @Accept json
// @Produce json
// @Param Authorization header string true "Authorization header with access token"
// @Success 200 "Access token is valid"
// @Header  200 {string} X-User-Id "Authenticated user UUID"
// @Failure 401 "Unauthorized"
// @Failure 405 "Method not allowed"
// @Failure 500 "Internal server error"
// @Router /check [get]
func CheckAccessToken(svc *token.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if strings.TrimSpace(authHeader) == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		userID, err := svc.ValidateAccessToken(authHeader)
		if err != nil {
			log.Printf("token service: %s", err)

			status, resp := mapErrToHTTP(err)
			writeJSON(w, status, resp)
			return
		}

		w.Header().Set("X-User-Id", userID.String())
		w.WriteHeader(http.StatusOK)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
