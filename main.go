package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// Service log setup
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("Found error in log ", err)
	}
	log.SetOutput(file)

	// Service setup
	done := doConsume(broker, group)
	router := server()
	serverErr := http.ListenAndServe(":6020", router)
	done <- true

	if serverErr != nil {
		log.Fatal(serverErr)
	}
}
