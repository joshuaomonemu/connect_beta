// Go connection Sample Code:
package db

import (
	"app/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/microsoft/go-mssqldb"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var server = "vidapp.database.windows.net"
var port = 1433
var user = "videoapproot"
var password = "Mylovefordogs1$"
var database = "videoapp"

func Conn() (*sql.DB, error) {
	// Build connection string
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)
	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	return db, nil
}

func EmailExists(email string) (bool, error) {
	db, _ := Conn()

	const query = "SELECT 1 FROM users WHERE email = ?"
	row := db.QueryRow(query, email)
	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func SetUser(user *models.User) error {
	db, err1 := Conn()
	if err1 != nil {
		return err1
	}
	query := `INSERT INTO users (email, fullname, displayname, password)
              VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, user.Email, user.Fullname, user.DisplayName, user.Password)

	return err
}

func LoginUser(email, password string) (bool, error) {
	db, _ := Conn()
	// Retrieve the OTP from the database

	var exists bool

	// Prepare the SQL query
	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ? AND password = ?)"

	// Execute the query
	err := db.QueryRow(query, email, password).Scan(&exists)
	if err != nil {
		return false, err
	}

	// Return the existence of the user with the given email and password
	return exists, nil
}

func DeleteUser(email string) (string, error) {
	db, _ := Conn()

	// Start a new transaction
	tx, err := db.Begin()
	if err != nil {
		return "", fmt.Errorf("could not begin transaction: %v", err)
	}

	// Defer a rollback in case of error
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete transactions linked to the email
	deleteTransactionsQuery := "DELETE FROM transactions WHERE user = ?"
	_, err = tx.Exec(deleteTransactionsQuery, email)
	if err != nil {
		return "", fmt.Errorf("could not delete transactions: %v", err)
	}

	// Delete user from the users table
	deleteUserQuery := "DELETE FROM users WHERE email = ?"
	_, err = tx.Exec(deleteUserQuery, email)
	if err != nil {
		return "", fmt.Errorf("could not delete user: %v", err)
	}

	// Commit the transaction if both deletions succeed
	if err = tx.Commit(); err != nil {
		return "", fmt.Errorf("could not commit transaction: %v", err)
	}

	return "ok", nil
}

func GetUserbyEmail(email string) (models.User, error) {
	db, _ := Conn()
	var user models.User

	// Prepare the SQL query
	query := "SELECT fullname, phone, wallet, email, verified_email FROM users WHERE email = ?"

	// Execute the query
	err := db.QueryRow(query, email).Scan(&user.Fullname, &user.Role, &user.Email, &user.DisplayName)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, err
		}
		return models.User{}, err
	}

	return user, nil
}
func SaveToDatabase(title, caption, file_url, user_id string) error {
	db, _ := Conn()
	query := `
		INSERT INTO photo (title, caption, file_url, user_id)
		VALUES (@p1, @p2, @p3, @p4)
	`

	_, err := db.Exec(query, title, caption, file_url, user_id)
	if err != nil {
		return fmt.Errorf("failed to insert photo metadata: %w", err)
	}

	return nil
}

func GetAllPhotos(c *fiber.Ctx) ([]models.Photo, error) {
	db, _ := Conn()
	rows, err := db.Query(`SELECT id, title, caption, file_url, user_id FROM photo ORDER BY id DESC`)
	if err != nil {
		log.Println("Query error:", err)
		return nil, fmt.Errorf("failed to query photos: %w", err)
	}
	defer rows.Close()

	var photos []models.Photo

	for rows.Next() {
		var p models.Photo
		if err := rows.Scan(&p.ID, &p.Title, &p.Caption, &p.FileURL, &p.UserID); err != nil {
			log.Println("Row scan error:", err)
			return nil, fmt.Errorf("failed to read photo row: %w", err)
		}
		photos = append(photos, p)
	}

	return photos, nil
}

func GetPhotoByID(userID string) ([]models.Photo, error) {
	db, _ := Conn()
	query := `SELECT id, title, caption, file_url, user_id, uploaded_at FROM photo WHERE user_id = @p1 ORDER BY id DESC`

	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query photos: %w", err)
	}
	defer rows.Close()

	var photos []models.Photo

	for rows.Next() {
		var p models.Photo
		if err := rows.Scan(&p.ID, &p.Title, &p.Caption, &p.FileURL, &p.UserID, &p.Uploaded); err != nil {
			return nil, fmt.Errorf("failed to scan photo row: %w", err)
		}
		photos = append(photos, p)
	}

	return photos, nil
}
