package main

import (
	"goDev/web-socket-chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// router creating
	mux := routes()

	// listening to client messages
	log.Println("starting channel listening")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")
	// starting webserver
	_ = http.ListenAndServe(":8080", mux)

}
