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
