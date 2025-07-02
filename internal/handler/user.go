package handler

import (
	"encoding/json"
	"go-user-api/internal/model"
	"go-user-api/internal/repository"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UserHandler struct {
	repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user and return the user object
// @Tags users
// @Accept  json
// @Produce  json
// @Param   user  body  model.User  true  "User Data"
// @Success 200 {object} model.User
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "Failed to create user"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u model.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid nput", http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(r.Context(), &u); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// GetAllUsers godoc
// @Summary Retrieve all users
// @Description Get a list of all users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} model.User
// @Failure 500 {string} string "Something went wrong"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.GetAllUsers(r.Context())
	if err != nil {
		http.Error(w, "something went wrong", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

// GetUser godoc
// @Summary Retrieve a user by ID
// @Description Get a user by their unique ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Success 200 {object} model.User
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	u, err := h.repo.Get(r.Context(), id)

	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// UpdateUser godoc
// @Summary Update a user by ID
// @Description Update a user by their unique ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Param   user  body  model.User  true  "User Data"
// @Success 200 {object} model.User
// @Failure 400 {string} string "Invalid input"
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	var u model.User

	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	u.ID = id
	if err := h.repo.Update(r.Context(), &u); err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(u)
}

// DeleteUser godoc
// @Summary Delete a user by ID
// @Description Delete a user by their unique ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   id  path  int  true  "User ID"
// @Success 204
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
