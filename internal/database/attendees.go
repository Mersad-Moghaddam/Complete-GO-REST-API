package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	UserID  int `json:"user_id" binding:"required"`
	EventID int `json:"event_id" binding:"required"`
}

func (m *AttendeeModel) Insert(Attendee *Attendee) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "INSERT INTO attendees (user_id, event_id) VALUES ($1, $2) RETURNING id"
	err := m.DB.QueryRowContext(ctx, query, Attendee.UserID, Attendee.EventID).Scan(&Attendee.ID)
	if err != nil {
		return nil, err
	}

	return Attendee, nil
}
func (m *AttendeeModel) GetByEventAndAttendee(eventID int, userID int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM attendees WHERE event_id = $1 AND user_id = $2"
	row := m.DB.QueryRowContext(ctx, query, eventID, userID)
	attendee := &Attendee{}
	err := row.Scan(&attendee.ID, &attendee.UserID, &attendee.EventID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return attendee, nil
}

func (m *AttendeeModel) GetAttendeesByEvent(eventId int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT u.* FROM attendees a JOIN users u ON a.user_id = u.id WHERE a.event_id = $1"
	rows, err := m.DB.QueryContext(ctx, query, eventId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.UserName, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (m *AttendeeModel) Delete(eventId int, userId int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "DELETE FROM attendees WHERE id = $1"
	_, err := m.DB.ExecContext(ctx, query, userId, eventId)
	if err != nil {
		return err
	}
	return nil
}

func (m *AttendeeModel) GetEventsByAttendee(attendeeId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := "SELECT e.* FROM attendees a JOIN events e ON a.event_id = e.id WHERE a.user_id = $1"
	rows, err := m.DB.QueryContext(ctx, query, attendeeId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	var events []*Event
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.OwnerId, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	return events, nil
}
