package main

import (
	"letters/db"
	"letters/functions"
	"time"

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
	log.Println("=1f95ae=", time.Now())

	db.Connect()
	functions.AutoSend()
	server.Start(host)

}
