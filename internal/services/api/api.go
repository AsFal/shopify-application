package api

import (
	"fmt"
	"log"
	"os"
)

func Run() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Valid Port is required to run api")
	}
	log.Println("something")
	router().Run(fmt.Sprintf(":%s", port))
}
