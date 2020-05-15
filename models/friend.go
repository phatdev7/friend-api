package models

import (
	"friend-api/db"
)

type Friends struct {
	Friends []string `json:"friends"`
}

func (friends *Friends) MakeFriend() error {
	var user1 *User
	var err error

	user1, err = GetOneUser(friends.Friends[0])
	if err != nil {
		return err
	}
	var user2 *User
	user2, err = GetOneUser(friends.Friends[1])
	if err != nil {
		return err
	}

	var user_id1, user_id2 int
	db := db.GetInstance()
	err = db.QueryRow("INSERT INTO friends (user_id1, user_id2) VALUES ($1, $2) RETURNING *",
		user1.ID, user2.ID).Scan(&user_id1, &user_id2)
	if err != nil {
		return err
	}
	return nil
}
