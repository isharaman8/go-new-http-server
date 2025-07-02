package routes

import (
	"go-user-api/internal/handler"

	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, userHandler *handler.UserHandler) {
	const userRouteWithId string = "/users/{id}"

	r.Post("/users", userHandler.CreateUser)
	r.Get(userRouteWithId, userHandler.GetUser)
	r.Put(userRouteWithId, userHandler.UpdateUser)
	r.Delete(userRouteWithId, userHandler.DeleteUser)
}
