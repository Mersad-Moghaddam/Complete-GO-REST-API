package database

import "database/sql"

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id" binding:"required"`
	EventID int `json:"event_id" binding:"required"`
}
