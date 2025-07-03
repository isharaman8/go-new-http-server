package handler

import (
	"encoding/json"
	"fmt"
	"go-user-api/internal/auth"
	"go-user-api/internal/middleware"
	"go-user-api/internal/model"
	"go-user-api/internal/repository"
	"net/http"

	"github.com/go-playground/validator"
)

var validate = validator.New()

type AuthRouteHandler struct {
	repo repository.UserRepository
}

func NewAuthRouteHandler(repo repository.UserRepository) *AuthRouteHandler {
	return &AuthRouteHandler{repo: repo}
}

// Signup godoc
// @Summary Signup a new user
// @Description Create a new user and return JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   user  body  model.User  true  "User Data"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid input"
// @Router /auth/signup [post]
func (h *AuthRouteHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var u model.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(u); err != nil {
		http.Error(w, "validation error"+err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := auth.HashPassword(u.Password)

	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}

	u.Password = hashedPassword

	if err := h.repo.Create(r.Context(), &u); err != nil {
		http.Error(w, "failed to signup new user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(u)

}

// Login godoc
// @Summary Login a user
// @Description Authenticate a user and return JWT token
// @Tags auth
// @Accept  json
// @Produce  json
// @Param   login  body  model.LoginInput  true  "Login Data"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/login [post]
func (h *AuthRouteHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input model.LoginInput

	// Step 1: Decode input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := validate.Struct(input); err != nil {
		http.Error(w, "validation error"+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetByEmail(r.Context(), input.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	if !auth.ComparePassword(input.Password, user.Password) {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	// generate new token
	token, jwtErr := auth.GenerateJWT(user.ID)
	if jwtErr != nil {
		http.Error(w, "could not generate jwt token", http.StatusUnauthorized)
		return
	}

	// send token
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}

// GetUserProfile godoc
// @Summary      Get the authenticated user's profile
// @Description  Requires a valid JWT token. Returns user info based on token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} model.User
// @Failure      400 {object} model.ErrorResponse
// @Failure      401 {object} model.ErrorResponse
// @Router       /auth/profile [get]
func (h *AuthRouteHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Unauthorized request: invalid token"})
	}

	fmt.Println("Got User ID:", userID)

	user, err := h.repo.Get(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "User not found"})
		return
	}

	json.NewEncoder(w).Encode(user)
}
