package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ravenocx/cat-socialx/cmd"
	"github.com/ravenocx/cat-socialx/internal/db"
)

func main() {
	godotenv.Load()

	dbConn, err := db.CreateConnection()
	if err != nil {
		log.Printf("Error connect to database : %+v", err)
	}

	defer func(){
		if err := dbConn.Close(); err != nil {
			log.Printf("Error close database connection : %+v", err)
		}
	}()

	h := cmd.New(&cmd.Http{
		DB : dbConn,
	})

	h.StartApp()
}
