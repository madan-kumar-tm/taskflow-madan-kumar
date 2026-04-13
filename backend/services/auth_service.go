package services

import (
	"taskflow/config"
	"taskflow/models"
	"taskflow/utils"
)

func RegisterUser(name, email, password string) error {
	hashed, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
	`

	return config.DB.Exec(query, name, email, string(hashed)).Error
}

func LoginUser(email, password string) (string, string, error) {
	var user models.User

	query := `SELECT id, name, email, password FROM users WHERE email=$1`
	err := config.DB.Raw(query, email).Scan(&user).Error
	if err != nil {
		return "", "", err
	}

	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return "", "", err
	}

	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}
	return token, user.Name, nil
}
