package services

import (
	"app/db"
	"app/models"
	"errors"
)

var user models.User

func Login(data *models.User) (interface{}, error) {
	if data.Email == "" {
		return "", errors.New("Email is required")
	}
	if data.Password == "" {
		return "", errors.New("Password is required")
	}

	exists, err := db.EmailExists(user.Email)
	if err != nil {
		return "Error checking user", err
	}
	if exists {
		loggedIn, err := db.LoginUser(user.Email, user.Password)
		if err != nil {
			return "Error signing in user", err
		}

		if loggedIn {
			us, _ := db.GetUserbyEmail(user.Email)

			return us, err
		} else {
			return "Invalid Password", err
		}
	} else {
		return "This email does not exist", err
	}
}

func SignUp(data *models.User) (interface{}, error) {
	if data.Email == "" {
		return "", errors.New("email is required")
	}
	if data.Password == "" {
		return "", errors.New("password is required")
	}
	if data.Fullname == "" {
		return "", errors.New("fullname is required")
	}

	exists, err := db.EmailExists(user.Email)
	if err != nil {
		return "Error checking user", err
	}
	if exists {
		return "User already exist", err
	}

	err1 := db.SetUser(&user)
	if err1 != nil {
		return "Error creating user", err
	}
	return data, nil
}
