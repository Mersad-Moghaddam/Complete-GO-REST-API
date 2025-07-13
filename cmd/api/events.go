package main

import (
	"go-rest/internal/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary Create event
// @Description Create a new event
// @Tags Events
// @Accept json
// @Produce json
// @Param event body database.Event true "Event info"
// @Success 201 {object} database.Event
// @Failure 400,500 {object} map[string]string
// @Router /events [post]
// @Security BearerAuth
func (app *Application) createEvent(c *gin.Context) {
	var event database.Event

	if err := c.ShouldBindBodyWithJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := app.GetUserFromContext(c)
	event.OwnerId = user.ID
	err := app.models.Events.Insert(&event)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in Inserting Event"})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// @Summary Get all events
// @Description Retrieve all events
// @Tags Events
// @Produce json
// @Success 200 {array} database.Event
// @Failure 500 {object} map[string]string
// @Router /events [get]
func (app *Application) getAllEvents(c *gin.Context) {
	events, err := app.models.Events.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

// @Summary Get event by ID
// @Description Retrieve event by ID
// @Tags Events
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} database.Event
// @Failure 400,404,500 {object} map[string]string
// @Router /events/{id} [get]
func (app *Application) getEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event, err := app.models.Events.Get(id)
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, event)
}

// @Summary Update event
// @Description Update an event by ID
// @Tags Events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body database.Event true "Event info"
// @Success 200 {object} database.Event
// @Failure 400,403,404,500 {object} map[string]string
// @Router /events/{id} [put]
// @Security BearerAuth
func (app *Application) updateEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	if existingEvent.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not Authorize"})
		return
	}
	updatedEvent := &database.Event{}
	if err := c.ShouldBindJSON(updatedEvent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	updatedEvent.ID = id // Ensure the ID is set to the existing event's ID
	if err := app.models.Events.Update(updatedEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update event"})
		return
	}
	c.JSON(http.StatusOK, updatedEvent)
}

// @Summary Delete event
// @Description Delete an event by ID
// @Tags Events
// @Param id path int true "Event ID"
// @Success 204
// @Failure 400,403,404,500 {object} map[string]string
// @Router /events/{id} [delete]
// @Security BearerAuth
func (app *Application) deleteEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	user := app.GetUserFromContext(c)
	existingEvent, err := app.models.Events.Get(id)
	if existingEvent == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event Not Found"})
		return
	}
	if existingEvent.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not Authorize to Delete"})
		return
	}
	err = app.models.Events.Delete(id)
	if err != nil {
		if err.Error() == "record not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete event"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// @Summary Add attendee to event
// @Description Add an attendee to an event
// @Tags Attendees
// @Param event_id path int true "Event ID"
// @Param user_id path int true "User ID"
// @Success 201 {object} database.Attendee
// @Failure 400,403,404,409,500 {object} map[string]string
// @Router /events/{event_id}/attendees/{user_id} [post]
// @Security BearerAuth
func (app *Application) addAttendeeToEvent(c *gin.Context) {
	eventId, err := strconv.Atoi(c.Param("event_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve the event and user from the database
	event, err := app.models.Events.Get(eventId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve event"})
		return
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}
	user := app.GetUserFromContext(c)
	if event.OwnerId != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not Authorize to Add Attendee"})
		return
	}
	// Retrieve the user from the database
	userToAdd, err := app.models.Users.GetUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user"})
		return
	}
	if userToAdd == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	existingAttendee, err := app.models.Attendees.GetByEventAndAttendee(event.ID, userToAdd.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve existing attendee"})
		return
	}
	if existingAttendee != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Attendee Already Exists"})
		return
	}
	// Add the attendee to the event
	attendee := database.Attendee{EventID: event.ID, UserID: userToAdd.ID}
	_, err = app.models.Attendees.Insert(&attendee)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add attendee"})
		return
	}
	c.JSON(http.StatusCreated, attendee)
}

// @Summary Get attendees for event
// @Description Get all attendees for an event
// @Tags Attendees
// @Param id path int true "Event ID"
// @Success 200 {array} database.User
// @Failure 400,500 {object} map[string]string
// @Router /events/{id}/attendees [get]
func (app *Application) getAttendeesForEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	users, err := app.models.Attendees.GetAttendeesByEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve attendees"})
		return
	}
	c.JSON(http.StatusOK, users)
}

// @Summary Delete attendee from event
// @Description Delete an attendee from an event
// @Tags Attendees
// @Param id path int true "Event ID"
// @Param user_id path int true "User ID"
// @Success 200 {string} string "Attendee Delete Successfully"
// @Failure 400,403,404,500 {object} map[string]string
// @Router /events/{id}/attendees/{user_id} [delete]
// @Security BearerAuth
func (app *Application) deleteAttendeeFromEvent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	event, err := app.models.Events.Get(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something Went Wrong"})
	}
	if event == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event Not Found"})
	}
	user := app.GetUserFromContext(c)
	if event.OwnerId != user.ID {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "You ate not Auth to Delete Attendee from Event"})
	}
	err = app.models.Attendees.Delete(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attendee"})
		return
	}

	c.JSON(http.StatusOK, "Attendee Delete Successfully")
}

// @Summary Get events by attendee
// @Description Get all events for an attendee
// @Tags Attendees
// @Param id path int true "Attendee ID"
// @Success 200 {array} database.Event
// @Failure 400,500 {object} map[string]string
// @Router /attendees/{id}/events [get]
func (app *Application) getEventsByAttendee(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Attendee ID"})
	}
	events, err := app.models.Attendees.GetEventsByAttendee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve events"})
		return
	}
	c.JSON(http.StatusOK, events)
}
