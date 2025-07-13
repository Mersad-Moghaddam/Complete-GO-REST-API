package database

import (
	"context"
	"database/sql"
	"time"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (m *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	return m.DB.QueryRowContext(ctx, query, user.UserName, user.Email, user.Password).Scan(&user.ID)
}

func (m *UserModel) GetUser(id int) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User
	query := `SELECT id, username, email, password FROM users WHERE id = $1`
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *UserModel) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User
	// Only select the columns needed for scanning
	query := `SELECT id, username, email, password FROM users WHERE email = $1`
	err := m.DB.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.UserName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
