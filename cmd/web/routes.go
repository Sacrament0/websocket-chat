package main

import (
	"goDev/web-socket-chat/internal/handlers"
	"net/http"

	"github.com/bmizerany/pat"
)

// serves client requests
func routes() http.Handler {
	// creating multiplexer
	mux := pat.New()
	// serves Home page
	mux.Get("/", http.HandlerFunc(handlers.Home))
	// serves websocket connection
	mux.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	// serves static content (nice alerts)
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
