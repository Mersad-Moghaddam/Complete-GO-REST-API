package main

import (
	"go-rest/internal/database"

	"github.com/gin-gonic/gin"
)

func (app *Application) GetUserFromContext(c *gin.Context) *database.User {
	contextUser, exist := c.Get("user")
	if !exist {
		return &database.User{}
	}
	user, ok := contextUser.(*database.User)
	if !ok {
		return nil
	}
	return user
}
