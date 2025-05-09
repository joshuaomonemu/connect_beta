package handlers

import (
	"app/db"

	"github.com/gofiber/fiber/v2"
)

// import "github.com/gofiber/fiber/v2"

// func LogoutHandler(c *fiber.Ctx) error {
// 	sess, err := store.Get(c)
// 	if err != nil {
// 		return err
// 	}

// 	sess.Destroy()
// 	return c.JSON(fiber.Map{"message": "Logged out successfully"})
// }

func GetAllPhotos(c *fiber.Ctx) error {
	photos, err := db.GetAllPhotos(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(photos)
}
