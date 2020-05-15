package models

import (
	"fmt"
	"friend-api/db"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Users []User

func GetUsers() (Users, error) {
	db := db.GetInstance()
	rows, err := db.Query("SELECT u.id, u.email FROM users u")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(Users, 0)
	for rows.Next() {
		var id int
		var email string
		if err := rows.Scan(&id, &email); err != nil {
			return users, err
		}
		users = append(users, User{
			ID:    id,
			Email: email,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func GetOneUser(email string) (*User, error) {
	db := db.GetInstance()
	rows, err := db.Query("SELECT * FROM users u WHERE u.email=$1", email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var id int
		if err = rows.Scan(&id, &email); err != nil {
			return nil, err
		}
		users = append(users, User{
			ID:    id,
			Email: email,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("Not found user with " + email)
	}
	return &users[0], nil
}

func (user *User) InsertUser() (*User, error) {
	var lastId int
	var email string
	db := db.GetInstance()
	err := db.QueryRow("INSERT INTO users (email) VALUES ($1) RETURNING *", user.Email).Scan(&lastId, &email)
	if err != nil {
		return nil, err
	}
	user.ID = lastId
	user.Email = email
	return user, nil
}
