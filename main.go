package main

import (
	"letters/db"
	"letters/functions"

	"letters/server"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	host := "127.0.0.1:80"
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Файл .env не найден")
		host = ":80"

	}
	db.Connect()
	slice := functions.CreateBirthdaysSlice()
	log.Println("=22b5b8=", slice)
	functions.AutoSend()
	server.Start(host)

}

