package v1

import (
	"encoding/json"
	"friend-api/models"
	"net/http"

	"github.com/go-chi/chi"
)

func friendRouter(r chi.Router) {
	r.Post("/list", getFriendListByEmail)
	r.Post("/mutual", getMutualFriends)
	r.Post("/make", makeFriend)
}

type EmailBody struct {
	Email string `json:"email"`
}

type SuccessRes struct {
	Success bool `json:"success"`
}

func getFriendListByEmail(w http.ResponseWriter, r *http.Request) {
	var emailBody EmailBody
	err := json.NewDecoder(r.Body).Decode(&emailBody)
	if err != nil {
		panic(err)
	}
	user, err := models.GetOneUser(emailBody.Email)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		friends, err := user.GetListFriend()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
		} else {
			json.NewEncoder(w).Encode(friends)
		}
	}
}

func getMutualFriends(w http.ResponseWriter, r *http.Request) {
	var emails models.Emails
	err := json.NewDecoder(r.Body).Decode(&emails)
	if err != nil {
		panic(err)
	}
	if emails.Emails[0] == "" || emails.Emails[1] == "" {
		w.WriteHeader(400)
		w.Write([]byte("MutualFriends length must be 2!"))
	} else if emails.Emails[0] == emails.Emails[1] {
		w.WriteHeader(400)
		w.Write([]byte("Both emails cannot be duplicated!"))
	} else {
		friends, err := models.GetMutualFriends(&emails)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
		} else {
			json.NewEncoder(w).Encode(friends)
		}
	}
}

func makeFriend(w http.ResponseWriter, r *http.Request) {
	var emails models.Emails
	err := json.NewDecoder(r.Body).Decode(&emails)
	if err != nil {
		panic(err)
	}
	if emails.Emails[0] == "" || emails.Emails[1] == "" {
		w.WriteHeader(400)
		w.Write([]byte("Friends length must be 2!"))
	} else if emails.Emails[0] == emails.Emails[1] {
		w.WriteHeader(400)
		w.Write([]byte("Both emails cannot be duplicated!"))
	} else {
		err := emails.MakeFriend()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
		} else {
			json.NewEncoder(w).Encode(SuccessRes{Success: true})
		}
	}
}
