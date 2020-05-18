package models

import (
	"fmt"
	"friend-api/db"
)

type Emails struct {
	Emails []string `json:"friends"`
}

func (emails *Emails) MakeFriend() error {
	var userOne *User
	var err error

	userOne, err = GetOneUser(emails.Emails[0])
	if err != nil {
		return err
	}
	var userTwo *User
	userTwo, err = GetOneUser(emails.Emails[1])
	if err != nil {
		return err
	}

	subcribe_status := checkSubcribed(userOne, userTwo)
	if subcribe_status == 0 {
		return fmt.Errorf("%s Was blocked %s!", userOne.Email, userTwo.Email)
	}
	subcribe_status = checkSubcribed(userTwo, userOne)
	if subcribe_status == 0 {
		return fmt.Errorf("%s Was blocked %s!", userTwo.Email, userOne.Email)
	}

	var user_one_id, user_two_id, status, user_action_id int
	db := db.GetInstance()
	if userOne.ID < userTwo.ID {
		err = db.QueryRow("INSERT INTO friends (user_one_id, user_two_id, status, user_action_id) VALUES ($1, $2, $3, $4) RETURNING *",
			userOne.ID, userTwo.ID, 1, userOne.ID).Scan(&user_one_id, &user_two_id, &status, &user_action_id)
	} else {
		err = db.QueryRow("INSERT INTO friends (user_one_id, user_two_id, status, user_action_id) VALUES ($1, $2, $3, $4) RETURNING *",
			userTwo.ID, userOne.ID, 1, userTwo.ID).Scan(&user_one_id, &user_two_id, &status, &user_action_id)
	}
	if err != nil {
		return err
	}
	return nil
}
