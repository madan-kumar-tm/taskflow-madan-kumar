package controllers

import (
	"net/http"
	"strings"
	"taskflow/services"
	"taskflow/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	UserName string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := make(map[string]string)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				field := fe.Field()
				tag := fe.Tag()

				switch field {
				case "UserName":
					switch tag {
					case "required":
						errors["username"] = "username is required"
					}

				case "Email":
					switch tag {
					case "required":
						errors["email"] = "email is required"
					case "email":
						errors["email"] = "email is incorrect"
					}

				case "Password":
					switch tag {
					case "required":
						errors["password"] = "Password is required"
					case "min":
						errors["password"] = "Password must be at least 6 characters"
					}
				}
			}
		} else {
			errors["body"] = "invalid request"
		}

		utils.ValidationFailed(c, errors)
		return
	}

	err := services.RegisterUser(req.UserName, req.Email, req.Password)
	if err != nil {
		println("REGISTER ERROR:", err.Error())

		if strings.Contains(err.Error(), "duplicate") {
			utils.ValidationFailed(c, map[string]string{
				"email": "email already registered",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered"})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errors := make(map[string]string)

		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				field := fe.Field()
				tag := fe.Tag()

				switch field {
				case "Email":
					switch tag {
					case "required":
						errors["email"] = "email is required"
					case "email":
						errors["email"] = "email is incorrect"
					}

				case "Password":
					switch tag {
					case "required":
						errors["password"] = "password is required"
					}
				}
			}
		} else {
			errors["body"] = "invalid request"
		}

		utils.ValidationFailed(c, errors)
		return
	}

	token, name, err := services.LoginUser(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{ 
			"name":  name,
			"email": req.Email,
		},
	})
}
