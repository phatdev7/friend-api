package models

import (
	"fmt"
	"friend-api/db"
	"strings"
	"sync"
	"time"
)

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type Users []User

func GetAllUsers() (Users, error) {
	q := `
	SELECT u.id, u.email
	FROM users u
	`
	rows, err := db.GetInstance().Query(q)
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

type UserResult struct {
	User *User
	Err  error
}

func GetOneUser(email string) UserResult {
	q := `
	SELECT *
	FROM users u
	WHERE u.email=$1
	`
	rows, err := db.GetInstance().Query(q, email)
	if err != nil {
		return UserResult{
			User: nil,
			Err:  err,
		}
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var id int
		if err = rows.Scan(&id, &email); err != nil {
			return UserResult{
				User: nil,
				Err:  err,
			}
		}
		users = append(users, User{
			ID:    id,
			Email: email,
		})
	}
	if err := rows.Err(); err != nil {
		return UserResult{
			User: nil,
			Err:  err,
		}
	}
	if len(users) == 0 {
		return UserResult{
			User: nil,
			Err:  fmt.Errorf("Not found user with " + email),
		}
	}

	time.Sleep(5 * time.Second)
	return UserResult{
		User: &users[0],
		Err:  nil,
	}
}

func (user *User) AddUser() (*User, error) {
	var lastId int
	var email string
	q := `
	INSERT INTO users (email)
	VALUES ($1) RETURNING *
	`
	err := db.GetInstance().QueryRow(q, user.Email).Scan(&lastId, &email)
	if err != nil {
		return nil, err
	}
	user.ID = lastId
	user.Email = email
	return user, nil
}

type FriendList struct {
	Success bool     `json:"success"`
	Friends []string `json:"friends"`
	Count   int      `json:"count"`
}

type FriendListResult struct {
	Data *FriendList
	Err  error
}

func (user *User) GetListFriend() FriendListResult {
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
	rows, err := db.GetInstance().Query(q, user.Email)
	if err != nil {
		return FriendListResult{
			Data: nil,
			Err:  err,
		}
	}
	defer rows.Close()

	friends := make([]string, 0)
	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		if err != nil {
			return FriendListResult{
				Data: nil,
				Err:  err,
			}
		}
		friends = append(friends, email)
	}
	if err = rows.Err(); err != nil {
		return FriendListResult{
			Data: nil,
			Err:  err,
		}
	}
	return FriendListResult{
		Data: &FriendList{
			Success: true,
			Friends: friends,
			Count:   len(friends),
		},
		Err: nil,
	}
}

func GetMutualFriends(emails *Emails) (*FriendList, error) {
	ch := make(chan UserResult, 2)
	go func() {
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			ch <- GetOneUser(emails.Emails[0])
			wg.Done()
		}()
		go func() {
			ch <- GetOneUser(emails.Emails[1])
			fmt.Println("22")
			wg.Done()
		}()
		go func() {
			ch <- GetOneUser(emails.Emails[1])
			fmt.Println("33")
			wg.Done()
		}()
		wg.Wait()
		close(ch)
	}()

	var userOne, userTwo UserResult
	userOne = GetOneUser(emails.Emails[0])
	userTwo = GetOneUser(emails.Emails[1])
	for v := range ch {
		fmt.Println(v.User)
		if emails.Emails[0] == v.User.Email {
			userOne = v
		} else {
			userTwo = v
		}
	}
	fmt.Println(userOne.User.Email)
	fmt.Printf(userTwo.User.Email)
	// userOne, userTwo = <-ch, <-ch

	if userOne.Err != nil {
		return nil, userOne.Err
	}
	if userTwo.Err != nil {
		return nil, userTwo.Err
	}

	resCh := make(chan FriendListResult, 2)
	resCh <- userOne.User.GetListFriend()
	resCh <- userTwo.User.GetListFriend()

	userOneListFriend, userTwoListFriend := <-resCh, <-resCh
	if userOneListFriend.Err != nil {
		return nil, userOneListFriend.Err
	}
	if userTwoListFriend.Err != nil {
		return nil, userTwoListFriend.Err
	}
	intersect := Intersect(userOneListFriend.Data.Friends, userTwoListFriend.Data.Friends)
	return &FriendList{
		Success: true,
		Friends: intersect,
		Count:   len(intersect),
	}, nil
}

func Intersect(a []string, b []string) []string {
	result := make([]string, 0)
	if len(a) <= len(b) {
		for _, v := range a {
			if Contains(b, v) {
				result = append(result, v)
			}
		}
		return result
	}
	for _, v := range b {
		if Contains(a, v) {
			result = append(result, v)
		}
	}
	return result
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
	ch := make(chan UserResult, 2)
	go func() {
		ch <- GetOneUser(Requestor)
	}()
	go func() {
		ch <- GetOneUser(Target)
	}()
	requestor, target := <-ch, <-ch
	if requestor.Err != nil {
		return requestor.Err
	}
	if target.Err != nil {
		return target.Err
	}
	subcribe_status := checkSubcribed(requestor.User, target.User)
	if subcribe_status == 1 {
		return fmt.Errorf("Subcribed!")
	}
	var requestorId, targetId, status int
	var q string
	if subcribe_status == 0 {
		q = `
		UPDATE subcribers
		SET status=$3
		WHERE requestor=$1 AND target=$2
											 AND status=0
											 RETURNING requestor, target, status
		`
	} else {
		q = `
		INSERT INTO subcribers (requestor, target, status)
		VALUES ($1, $2, $3)
		RETURNING requestor, target, status
		`
	}
	err := db.GetInstance().QueryRow(q, requestor.User.ID, target.User.ID, 1).
		Scan(&requestorId, &targetId, &status)
	if err != nil {
		return err
	}
	return nil
}

func checkSubcribed(requestor *User, target *User) int {
	q := `
	SELECT status
	FROM subcribers
	WHERE requestor=$1 AND target=$2
	`
	rows, err := db.GetInstance().Query(q, requestor.ID, target.ID)
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
		return -1
	}
	if len(subcribers) == 0 {
		return -1
	}
	return subcribers[0]
}

func BlockUser(Requestor string, Target string) error {
	ch := make(chan UserResult, 2)
	go func() {
		ch <- GetOneUser(Requestor)
	}()
	go func() {
		ch <- GetOneUser(Target)
	}()

	requestor, target := <-ch, <-ch
	if requestor.Err != nil {
		return requestor.Err
	}
	if target.Err != nil {
		return target.Err
	}
	subcribe_status := checkSubcribed(requestor.User, target.User)
	if subcribe_status == 0 {
		return fmt.Errorf("Blocked!")
	}
	var q string
	var requestorId, targetId, status int
	if subcribe_status == 1 {
		q = `
		UPDATE subcribers
		SET status=$3
		WHERE requestor=$1 AND target=$2 AND status=1
		RETURNING requestor, target, status
		`
	} else {
		q = `
		INSERT INTO subcribers (requestor, target, status)
		VALUES ($1, $2, $3) RETURNING requestor, target, status
		`
	}
	err := db.GetInstance().QueryRow(q, requestor.User.ID, target.User.ID, 0).
		Scan(&requestorId, &targetId, &status)
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
	user := GetOneUser(body.Sender)
	if user.Err != nil {
		return nil, user.Err
	}
	friends := user.User.GetListFriend()
	if friends.Err != nil {
		return nil, friends.Err
	}
	followers, err := user.User.getSubcribersStatus(1)
	if err != nil {
		return nil, err
	}
	blockers, err := user.User.getSubcribersStatus(0)
	if err != nil {
		return nil, err
	}
	combine := make([]string, 0)
	combine = append(combine, getMention(body.Text)...)
	combine = append(combine, friends.Data.Friends...)
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
	ch := make(chan string, 5)

	go func() {
		var wg sync.WaitGroup
		for _, v := range emails {
			wg.Add(1)
			go func(email string) {
				user := GetOneUser(email)
				if user.Err == nil && user.User != nil {
					ch <- user.User.Email
				}
				wg.Done()
			}(v)
		}
		wg.Wait()
		close(ch)
	}()

	for v := range ch {
		fmt.Println(v)
		users = append(users, v)
	}
	return users
}

func (user *User) getSubcribersStatus(status int) (*[]string, error) {
	q := `
	SELECT email
	FROM users
	WHERE id IN (SELECT requestor
							 FROM subcribers s JOIN users u
							 									 ON s.target=u.id AND s.status=$1 AND u.email=$2)
	`
	rows, err := db.GetInstance().Query(q, status, user.Email)
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

type Emails struct {
	Emails []string `json:"friends"`
}

func (emails *Emails) MakeFriend() error {
	var err error
	var userOne, userTwo UserResult

	ch := make(chan UserResult, 5)
	go func() {
		ch <- GetOneUser(emails.Emails[0])
	}()
	go func() {
		ch <- GetOneUser(emails.Emails[1])
	}()

	userOne, userTwo = <-ch, <-ch
	if userOne.Err != nil {
		return err
	}
	if userTwo.Err != nil {
		return err
	}

	subcribe_status := checkSubcribed(userOne.User, userTwo.User)
	if subcribe_status == 0 {
		return fmt.Errorf("%s Was blocked %s!", userOne.User.Email, userTwo.User.Email)
	}
	subcribe_status = checkSubcribed(userTwo.User, userOne.User)
	if subcribe_status == 0 {
		return fmt.Errorf("%s Was blocked %s!", userTwo.User.Email, userOne.User.Email)
	}

	var user_one_id, user_two_id, status, user_action_id int
	q := `
	INSERT INTO friends (user_one_id, user_two_id, status, user_action_id)
	VALUES ($1, $2, $3, $4)
	RETURNING *
	`
	db := db.GetInstance()
	if userOne.User.ID < userTwo.User.ID {
		err = db.QueryRow(q, userOne.User.ID, userTwo.User.ID, 1, userOne.User.ID).
			Scan(&user_one_id, &user_two_id, &status, &user_action_id)
	} else {
		err = db.QueryRow(q, userTwo.User.ID, userOne.User.ID, 1, userTwo.User.ID).
			Scan(&user_one_id, &user_two_id, &status, &user_action_id)
	}
	if err != nil {
		return err
	}
	return nil
}
