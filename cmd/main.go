package main

import (
	"log"
	"GoGoGo/pkg/handler"
	serv "GoGoGo"
)

func main() {

	handlers := new(handler.Handler)

	srv := new(serv.Server)

	err := srv.Run("8080", handlers.NewRouter())
	if err != nil {
		log.Fatal(err)
	}
}
