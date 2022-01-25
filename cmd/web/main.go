package main

import (
	"goDev/web-socket-chat/internal/handlers"
	"log"
	"net/http"
)

func main() {
	// создаем маршрутизатор
	mux := routes()

	//асинхронно раскидываем клиентам полученные сообщения
	log.Println("starting channel listening")
	go handlers.ListenToWsChannel()

	log.Println("Starting web server on port 8080")
	// запускаем tcp сервер с указанием маршрутизатора
	_ = http.ListenAndServe(":8080", mux)

}
