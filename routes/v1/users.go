package v1

import (
	"encoding/json"
	"friend-api/models"
	"net/http"

	"github.com/go-chi/chi"
)

func userRouter(r chi.Router) {
	r.Get("/", getAllUsers)
	r.Post("/", addUser)
	r.Post("/subcribe", subcribeUser)
	r.Post("/block", blockUser)
	r.Post("/publish", publish)

	r.Group(func(r chi.Router) {
		r.Route("/friend", func(friendRoute chi.Router) {
			friendRoute.Post("/", getFriendListByEmail)
			friendRoute.Post("/mutual", getMutualFriends)
			friendRoute.Post("/make", makeFriend)
		})
	})
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetAllUsers()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	usersRes, _ := json.Marshal(users)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(usersRes))
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	newUser, err := user.AddUser()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(newUser)
}

func subcribeUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Requestor string `json:"requestor"`
		Target    string `json:"target"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if body.Requestor == "" || body.Target == "" {
		w.WriteHeader(400)
		w.Write([]byte("requestor and target should not be blank!"))
		return
	}
	if body.Requestor == body.Target {
		w.WriteHeader(400)
		w.Write([]byte("Both emails cannot be duplicated!"))
		return
	}
	err = models.SubcribeUser(body.Requestor, body.Target)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(SuccessRes{Success: true})
}

func blockUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Requestor string `json:"requestor"`
		Target    string `json:"target"`
	}
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	err = models.BlockUser(body.Requestor, body.Target)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(SuccessRes{Success: true})
}

func publish(w http.ResponseWriter, r *http.Request) {
	var body models.PublishBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	res, err := models.Publish(&body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(res)
}

type SuccessRes struct {
	Success bool `json:"success"`
}

func getFriendListByEmail(w http.ResponseWriter, r *http.Request) {
	type EmailBody struct {
		Email string `json:"email"`
	}
	var emailBody EmailBody
	err := json.NewDecoder(r.Body).Decode(&emailBody)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if emailBody.Email == "" {
		w.WriteHeader(400)
		w.Write([]byte("Email can not blank!"))
		return
	}
	user := models.GetOneUser(emailBody.Email)
	if user.Err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	friends := user.User.GetListFriend()
	if friends.Err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(friends.Data)
}

func getMutualFriends(w http.ResponseWriter, r *http.Request) {
	var emails models.Emails
	err := json.NewDecoder(r.Body).Decode(&emails)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if len(emails.Emails) < 2 || emails.Emails[0] == "" || emails.Emails[1] == "" {
		w.WriteHeader(400)
		w.Write([]byte("MutualFriends length must be 2!"))
		return
	}
	if emails.Emails[0] == emails.Emails[1] {
		w.WriteHeader(400)
		w.Write([]byte("Both emails cannot be duplicated!"))
		return
	}
	friends, err := models.GetMutualFriends(&emails)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(friends)
}

func makeFriend(w http.ResponseWriter, r *http.Request) {
	var emails models.Emails
	err := json.NewDecoder(r.Body).Decode(&emails)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	if len(emails.Emails) < 2 || emails.Emails[0] == "" || emails.Emails[1] == "" {
		w.WriteHeader(400)
		w.Write([]byte("Friends length must be 2!"))
		return
	}
	if emails.Emails[0] == emails.Emails[1] {
		w.WriteHeader(400)
		w.Write([]byte("Both emails cannot be duplicated!"))
		return
	}
	err = emails.MakeFriend()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(SuccessRes{Success: true})
}
