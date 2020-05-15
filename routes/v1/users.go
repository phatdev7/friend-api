package v1

import (
	"encoding/json"
	"friend-api/models"
	"net/http"

	"github.com/go-chi/chi"
)

func userRouter(r chi.Router) {
	r.Get("/", getUsers)
	r.Post("/", insertUser)
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users, err := models.GetUsers()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		usersRes, _ := json.Marshal(users)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte(usersRes))
		// json.NewEncoder(w).Encode(user)
	}
}

func insertUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}
	result, err := user.InsertUser()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(result)
	}
}
