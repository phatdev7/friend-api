package v1

import (
	"encoding/json"
	"fmt"
	"friend-api/models"
	"net/http"

	"github.com/go-chi/chi"
)

func friendRouter(r chi.Router) {
	r.Get("/", getFriendListByEmail)
	r.Post("/", makeFriend)
}

type EmailBody struct {
	Email string `json:"email"`
}

func getFriendListByEmail(w http.ResponseWriter, r *http.Request) {
	var emailBody EmailBody
	err := json.NewDecoder(r.Body).Decode(&emailBody)
	if err != nil {
		panic(err)
	}
	var user &User
	if user, err = models.GetOneUser(emailBody.Email); err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	}
}

func makeFriend(w http.ResponseWriter, r *http.Request) {
	var friends models.Friends
	err := json.NewDecoder(r.Body).Decode(&friends)
	if err != nil {
		panic(err)
	}
	if friends.Friends[0] != "" && friends.Friends[1] != "" {
		err := friends.MakeFriend()
		type Response struct {
			Success bool `json:"success"`
		}
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
		} else {
			json.NewEncoder(w).Encode(Response{Success: true})
		}
	} else {
		fmt.Println("Friends length must be 2")
	}
}
