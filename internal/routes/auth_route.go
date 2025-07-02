package routes

import (
	"go-user-api/internal/handler"
	"go-user-api/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, authHandler *handler.AuthRouteHandler) {
	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)
	r.With(middleware.JWTAuthMiddleware).Get("/auth/profile", authHandler.GetUserProfile)
}
