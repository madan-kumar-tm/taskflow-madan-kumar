package controllers

import (
	"net/http"
	"taskflow/services"

	"github.com/gin-gonic/gin"
)

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func ListUsers(c *gin.Context) {
	users, err := services.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch users"})
		return
	}

	result := make([]userResponse, 0, len(users))
	for _, user := range users {
		result = append(result, userResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}
	c.JSON(http.StatusOK, result)
}
