package handlers

import (
	"app/db"
	_ "app/db"
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/gofiber/fiber/v2"
)

const (
	containerName = "videoapplication"
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func HandleUpload(c *fiber.Ctx) error {

	// Extract user ID from cookie (or session storage)
	userID := c.Cookies("user_id")
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("User not authenticated")
	}

	// Parse form fields
	title := c.FormValue("title")
	caption := c.FormValue("caption")

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return c.Status(400).SendString("File missing")
	}

	// Open file

	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	// Upload to Azure Blob
	blobURL, err := uploadToBlob(file, fileHeader.Filename)
	if err != nil {
		return err
	}

	// Save metadata to SQL
	err = db.SaveToDatabase(title, caption, blobURL, userID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"message": "Uploaded successfully",
		"url":     blobURL,
	})
}

func uploadToBlob(file multipart.File, filename string) (string, error) {

	account_pass := os.Getenv("AZURE_STORAGE_ACCOUNT")

	accountName := "videoapp"
	serviceURL := fmt.Sprintf("https://%s.blob.core.windows.net/", accountName)

	client, err := azblob.NewClientFromConnectionString(account_pass, nil)
	handleError(err)

	// Build blob name with timestamp to avoid collisions
	blobName := fmt.Sprintf("uploads/%d_%s", time.Now().Unix(), filename)

	// Upload the stream
	_, err = client.UploadStream(context.TODO(), containerName, blobName, file, nil)
	if err != nil {
		return "", err
	}

	// Return public URL (assuming blob is publicly accessible or SAS is added later)
	return serviceURL + containerName + "/" + blobName, nil
}

func GetPhotoByIDHandler(c *fiber.Ctx) error {
	userID := c.Params("user_id")
	if userID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "User ID is required"})
	}

	photos, err := db.GetPhotoByID(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(photos)
}
