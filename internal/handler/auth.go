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
	}

	u.Password = hashedPassword

	if err := h.repo.Create(r.Context(), &u); err != nil {
		http.Error(w, "failed to signup new user", http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(u)
}

func (h *AuthRouteHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required,min=3"`
	}

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

func (h *AuthRouteHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)

	fmt.Println("Got User ID:", userID)

	user, err := h.repo.Get(r.Context(), userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}
