package routes

import (
	"app/handlers"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"

	"github.com/gofiber/fiber/v2"
)

// The function that runs all routes and starts the serves

func Run() {
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost", // Change to your frontend domain if needed
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

	// Auth routes
	app.Post("/auth/login", handlers.Login)
	app.Post("/auth/signup", handlers.Signup)

	//app.Post("/api/upload", handlers.JWTMiddleware, handlers.RequireRole("creator"), handlers.HandleUpload)
	// app.Post("/api/comment/id", handlers.JWTMiddleware, handCommentHandler)
	// app.Get("/api/comments/id", handlers.JWTMiddleware, CommentHandler)

	// Upload and photo APIs
	app.Post("/api/upload", handlers.HandleUpload)
	// app.Post("/api/upload", handlers.JWTMiddleware, handlers.RequireRole("creator"), handlers.HandleUpload)

	app.Get("/api/photo/all", handlers.GetAllPhotos) // Remove duplicate route

	app.Get("/api/:user_id/photos", handlers.GetPhotoByIDHandler)

	// Other API routes
	app.Post("/see", handlers.JWTMiddleware, handlers.Seerer)

	// Start the server
	log.Fatal(app.Listen(":2020"))
}
