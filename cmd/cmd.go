package cmd

import (
	"friend-api/db"
	"friend-api/server"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err.Error())
	}
	godotenv.Load(wd + "/.env")
}

func Execute() {
	db.Init()
	dbInstance := db.GetInstance()
	defer dbInstance.Close()

	s := server.NewServer()
	s.Start()
}
