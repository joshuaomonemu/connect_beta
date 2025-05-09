package routes

import (
	"app/handlers"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
)

// The function that runs all routes and starts the server
func Run() {
	app := fiber.New()

	// Enable CORS for all origins
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Or restrict to specific origin like "http://localhost:5500"
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	app.Post("/auth/login", handlers.Login)
	app.Post("/auth/signup", handlers.Signup)
	app.Post("/api/upload", handlers.HandleUpload)
	//app.Post("/api/upload", handlers.JWTMiddleware, handlers.RequireRole("creator"), handlers.HandleUpload)
	// app.Post("/api/comment/id", handlers.JWTMiddleware, handCommentHandler)
	// app.Get("/api/comments/id", handlers.JWTMiddleware, CommentHandler)
	app.Get("/api/photo/all", handlers.GetAllPhotos)
	app.Get("/api/photo/all", handlers.GetAllPhotos)
	app.Get("/api/:user_id/photos", handlers.GetPhotoByIDHandler)
	app.Post("/see", handlers.JWTMiddleware, handlers.Seerer)

	app.Listen(":2020")
}
