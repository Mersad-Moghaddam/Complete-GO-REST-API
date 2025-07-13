package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *Application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		// Event Routes
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)
		v1.GET("/events/:id/attendees", app.getAttendeesForEvent)
		v1.GET("/attendees/:id/events", app.getEventsByAttendee)
		// User Routes
		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)

	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)
		authGroup.POST("/events/:event_id/attendees/:user_id", app.addAttendeeToEvent)
		authGroup.DELETE("/events/:id/attendees/:user_id", app.deleteAttendeeFromEvent)

	}
	return g
}
