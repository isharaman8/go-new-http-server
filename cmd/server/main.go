package main

import (
	"go-user-api/internal/db"
	"go-user-api/internal/handler"
	"go-user-api/internal/middleware"
	"go-user-api/internal/repository"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	// load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found proceeding with system env vars")
	}

	conn, err := db.Connect()

	if err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	defer conn.Close()

	UserRepo := repository.NewUserRepo(conn)
	userHandler := handler.NewUserHandler(UserRepo)
	authHandler := handler.NewAuthRouteHandler(UserRepo)

	r := chi.NewRouter()

	r.Post("/users", userHandler.CreateUser)
	r.Get("/users/{id}", userHandler.GetUser)
	r.Put("/users/{id}", userHandler.UpdateUser)
	r.Delete("/users/{id}", userHandler.DeleteUser)
	r.Post("/auth/signup", authHandler.Signup)
	r.Post("/auth/login", authHandler.Login)
	r.With(middleware.JWTAuthMiddleware).Get("/auth/profile", authHandler.GetUserProfile)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
