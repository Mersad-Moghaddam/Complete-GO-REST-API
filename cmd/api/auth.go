package main

import (
	"go-rest/internal/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserName string `json:"username" binding:"required,min=2"`
}

func (app *Application) registerUser(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went Wrong in Generating Password"})
		return
	}

	req.Password = string(hashedPassword)
	user := database.User{
		Email:    req.Email,
		Password: req.Password,
		UserName: req.UserName,
	}
	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went Wrong in Inserting User"})
		return
	}
	c.JSON(http.StatusCreated, user)
}
