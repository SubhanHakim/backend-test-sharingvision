package main

import (
	"log"
	"net/http"

	"article/internal/article"
	"article/internal/config"
	"article/internal/database"
	"article/internal/httpx"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mux := http.NewServeMux()
	article.NewHandler(article.NewService(article.NewRepository(db))).Register(mux)

	log.Println("listening on", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, httpx.CORS(cfg.CORSOrigin, mux)))
}
