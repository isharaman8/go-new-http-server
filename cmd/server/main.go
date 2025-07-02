// @title       Go User API
// @version     1.0
// @description This is a sample user API using JWT authentication
// @host        localhost:8080
// @BasePath    /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token.
package main

import (
	_ "go-user-api/docs"
	"go-user-api/internal/db"
	"go-user-api/internal/handler"
	"go-user-api/internal/repository"
	"go-user-api/internal/routes"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
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

	// register routes
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	routes.RegisterUserRoutes(r, userHandler)
	routes.RegisterAuthRoutes(r, authHandler)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}
