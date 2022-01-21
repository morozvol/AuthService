package main

import (
	"github.com/morozvol/AuthService/internal/app/apiserver"
	"github.com/morozvol/AuthService/internal/app/config"
	"log"
)

func main() {
	config, err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
