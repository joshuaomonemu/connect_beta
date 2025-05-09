package handlers

import (
	"app/models"
	"app/services"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// The byte slice that contains user data from the request body
var bs []byte

// Secret key for signing JWT
var secretKey = []byte("supersecretkey")

// Login Handler
func Login(c *fiber.Ctx) error {

	//Parsing the user schema from the request body to our userdata struct
	user := new(models.User)

	if err := c.BodyParser(user); err != nil {
		return err
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Fullname,
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	//Sending the user data over to the auth service for it's logic
	response, err := services.Login(user)

	if err != nil {
		resp := models.AuthResp{
			Status:  "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(resp)
		c.Send(bs)
	}

	// Return data if login successful
	resp := models.AuthResp{
		Status:  "success",
		Message: "User login successful",
		Data:    response,
		Token:   tokenString,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// Signup Handler
func Signup(c *fiber.Ctx) error {
	//Parsing the user schema from the request body to our userdata struct
	user := new(models.User)

	if err := c.ReqHeaderParser(user); err != nil {
		return err
	}
	//Sending the user data over to the auth service for it's logic
	status, err := services.SignUp(user)
	//Return error on unsuccessful Registration
	if err != nil {
		resp := models.AuthResp{
			Status:  "error",
			Message: err.Error(),
		}
		bs, _ := json.Marshal(resp)
		c.Send(bs)
		return nil
	}

	//Return data if login successful
	resp := models.AuthResp{
		Status:  "success",
		Message: "User registration successful",
		Data:    status,
	}
	bs, _ := json.Marshal(resp)
	c.Send(bs)
	return err
}

// JWT Middleware (Protects Routes)
func JWTMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Extract claims and check expiration
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// Check if token has expired
	exp, ok := claims["exp"].(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid expiration time"})
	}

	if time.Now().Unix() > int64(exp) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token has expired"})
	}

	return c.Next()
}

// Role enforcement to ensure the endpoint is being requested by a role manager
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("role") != role {
			return c.Status(403).JSON(fiber.Map{"error": "Access denied"})
		}
		return c.Next()
	}
}

func Seerer(c *fiber.Ctx) error {
	//c.Set("Authorization", "Bearer "+c.Params("token"))
	return c.SendString("Hello, World!")
}
