package services

import (
	"taskflow/config"
	"taskflow/models"
)

func ListUsers() ([]models.User, error) {
	users := []models.User{}
	query := `
		SELECT id, name, email, CAST(created_at AS TEXT) AS created_at
		FROM users
		ORDER BY name ASC
	`
	err := config.DB.Raw(query).Scan(&users).Error
	return users, err
}
