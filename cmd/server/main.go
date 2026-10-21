package main

import (
	"log"
	"net/http"

	"article/internal/app"
	"article/internal/config"
)

func main() {
	cfg := config.Load()

	h, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, h))
}
