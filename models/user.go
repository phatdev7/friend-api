package models

import (
	"fmt"
	"friend-api/db"
	"strings"
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

func (user *User) AddUser() (*User, error) {
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

type FriendListStruct struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

func (user *User) GetListFriend() (*FriendListStruct, error) {
	db := db.GetInstance()
	q := `
	SELECT u.email
	FROM users u
	WHERE u.id IN (SELECT user_two_id
								 FROM friends f, users u
								 WHERE f.user_one_id=u.id AND u.email=$1)
				OR u.id IN (SELECT user_one_id
										FROM friends f, users u
										WHERE f.user_two_id=u.id AND u.email=$1);
	`
	rows, err := db.Query(q, user.Email)
	if err != nil {
		return nil, err
	}

	friends := make([]string, 0)
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		friends = append(friends, email)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &FriendListStruct{
		Success: true,
		Friends: friends,
		Count:   len(friends),
	}, nil
}

func GetMutualFriends(emails *Emails) (*FriendListStruct, error) {
	userOne, err := GetOneUser(emails.Emails[0])
	if err != nil {
		return nil, err
	}
	userTwo, err := GetOneUser(emails.Emails[1])
	if err != nil {
		return nil, err
	}
	userOneListFriend, err := userOne.GetListFriend()
	if err != nil {
		return nil, err
	}
	userTwoListFriend, err := userTwo.GetListFriend()
	if err != nil {
		return nil, err
	}
	intersect := Intersect(userOneListFriend.Friends, userTwoListFriend.Friends)
	return &FriendListStruct{
		Success: true,
		Friends: intersect,
		Count:   len(intersect),
	}, nil
}

func Intersect(a []string, b []string) (result []string) {
	if len(a) <= len(b) {
		for _, v := range a {
			if Contains(b, v) {
				result = append(result, v)
			}
		}
	} else {
		for _, v := range b {
			if Contains(a, v) {
				result = append(result, v)
			}
		}
	}
	return
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func SubcribeUser(Requestor string, Target string) error {
	requestor, err := GetOneUser(Requestor)
	if err != nil {
		return err
	}
	target, err := GetOneUser(Target)
	if err != nil {
		return err
	}
	subcribe_status := checkSubcribed(requestor, target)
	if subcribe_status == 1 {
		return fmt.Errorf("Subcribed!")
	}
	if subcribe_status == 0 {
		q := `
		UPDATE subcribers SET status=$3 WHERE requestor=$1 AND target=$2 AND status=0 RETURNING requestor, target, status
		`
		var requestorId, targetId, status int
		db := db.GetInstance()
		err := db.QueryRow(q, requestor.ID, target.ID, 1).Scan(&requestorId, &targetId, &status)
		if err != nil {
			return err
		}
		return nil
	}
	q := `
	INSERT INTO subcribers (requestor, target, status) VALUES ($1, $2, $3) RETURNING requestor, target, status
	`
	var requestorId, targetId, status int
	db := db.GetInstance()
	err = db.QueryRow(q, requestor.ID, target.ID, 1).Scan(&requestorId, &targetId, &status)
	if err != nil {
		return err
	}
	return nil
}

func checkSubcribed(requestor *User, target *User) int {
	q := `
	SELECT status FROM subcribers WHERE requestor=$1 AND target=$2
	`
	db := db.GetInstance()
	rows, err := db.Query(q, requestor.ID, target.ID)
	if err != nil {
		return -1
	}
	subcribers := make([]int, 0)
	for rows.Next() {
		var status int
		err = rows.Scan(&status)
		if err != nil {
			return -1
		}
		subcribers = append(subcribers, status)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return -1
	}
	if len(subcribers) == 0 {
		return -1
	}
	return subcribers[0]
}

func BlockUser(Requestor string, Target string) error {
	requestor, err := GetOneUser(Requestor)
	if err != nil {
		return err
	}
	target, err := GetOneUser(Target)
	if err != nil {
		return err
	}
	subcribe_status := checkSubcribed(requestor, target)
	if subcribe_status == 0 {
		return fmt.Errorf("Blocked!")
	}
	if subcribe_status == 1 {
		q := `
		UPDATE subcribers
		SET status=$3
		WHERE requestor=$1 AND target=$2 AND status=1
		RETURNING requestor, target, status
		`
		var requestorId, targetId, status int
		db := db.GetInstance()
		err := db.QueryRow(q, requestor.ID, target.ID, 0).Scan(&requestorId, &targetId, &status)
		if err != nil {
			return err
		}
		return nil
	}
	q := `
	INSERT INTO subcribers (requestor, target, status)
	VALUES ($1, $2, $3) RETURNING requestor, target, status
	`
	var requestorId, targetId, status int
	db := db.GetInstance()
	err = db.QueryRow(q, requestor.ID, target.ID, 0).Scan(&requestorId, &targetId, &status)
	if err != nil {
		return err
	}
	return nil
}

type PublishBody struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

type PublishRes struct {
	Success    bool     `json:"success"`
	Recipients []string `json:"recipients"`
}

func Publish(body *PublishBody) (*PublishRes, error) {
	user, err := GetOneUser(body.Sender)
	if err != nil {
		return nil, err
	}
	friends, err := user.GetListFriend()
	if err != nil {
		return nil, err
	}
	followers, err := user.getSubcribersStatus(1)
	if err != nil {
		return nil, err
	}
	blockers, err := user.getSubcribersStatus(0)
	if err != nil {
		return nil, err
	}
	combine := make([]string, 0)
	combine = append(combine, getMention(body.Text)...)
	combine = append(combine, friends.Friends...)
	combine = append(combine, *followers...)
	combine = unique(combine)
	combine = removeEmail(combine, *blockers)

	return &PublishRes{
		Success:    true,
		Recipients: combine,
	}, nil
}

func removeEmail(emails []string, removeItems []string) []string {
	for _, v := range removeItems {
		for i := 0; i < len(emails); i++ {
			if emails[i] == v {
				emails = append(emails[:i], emails[i+1:]...)
				// i-- // form the remove item index to start iterate next item
			}
		}
	}

	return emails
}

func getMention(s string) []string {
	arr := strings.Split(s, " ")
	emails := make([]string, 0)
	for _, v := range arr {
		if strings.Contains(v, "@") {
			emails = append(emails, v)
		}
	}

	users := make([]string, 0)
	for _, v := range emails {
		u, err := GetOneUser(v)
		if err == nil && u != nil {
			users = append(users, u.Email)
		}
	}
	return users
}

func (user *User) getSubcribersStatus(status int) (*[]string, error) {
	q := `
	SELECT email
	FROM users
	WHERE id IN (SELECT requestor
							 FROM subcribers s JOIN users u ON s.target=u.id AND s.status=$1 AND u.email=$2)
	`
	db := db.GetInstance()
	rows, err := db.Query(q, status, user.Email)
	if err != nil {
		return nil, err
	}
	emails := make([]string, 0)
	for rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return &emails, nil
}
