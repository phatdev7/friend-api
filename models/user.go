package models

import "fmt"

type User struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func Hello() {
	fmt.Println("Hello models")
}
