package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"friend-api/db"
	"friend-api/models"
	"net/http"

	"github.com/go-chi/chi"
)

func friendRouter(r chi.Router) {
	r.Get("/", getAll)
}

func getAll(w http.ResponseWriter, r *http.Request) {
	var ctx context.Context
	db := db.GetInstance()
	age := "27"
	rows, err := db.QueryContext(ctx, "SELECT email FROM users;", age)
	if err != nil {
		fmt.Println("select fail")
	}
	emails := make([]string, 0)
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			fmt.Println("not found")
		}
		emails = append(emails, email)
	}
	fmt.Println(emails)

	user := models.User{ID: 1, Email: "kaken@gmail.com2"}
	resUser, _ := json.Marshal(user)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(201)
	w.Write([]byte(resUser))
	// json.NewEncoder(w).Encode(user)
}
