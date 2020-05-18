package v1

import (
	"encoding/json"
	"friend-api/models"
	"net/http"

	"github.com/go-chi/chi"
)

func userRouter(r chi.Router) {
	r.Get("/", getUsers)
	r.Post("/", addUser)
	r.Post("/subcribe", subcribeUser)
	r.Post("/block", blockUser)
	r.Post("/publish", publish)
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
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		panic(err)
	}
	result, err := user.AddUser()
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

type SubcribeBody struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

func subcribeUser(w http.ResponseWriter, r *http.Request) {
	var body SubcribeBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
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

type BlockBody struct {
	Requestor string `json:"requestor"`
	Target    string `json:"target"`
}

func blockUser(w http.ResponseWriter, r *http.Request) {
	var body BlockBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(""))
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
