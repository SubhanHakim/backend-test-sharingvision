// Package app merakit seluruh dependensi menjadi satu http.Handler,
// dipakai bersama oleh server lokal (cmd/server) dan Vercel (api/).
package app

import (
	"net/http"

	"article/internal/article"
	"article/internal/config"
	"article/internal/database"
	"article/internal/httpx"
)

func New(cfg config.Config) (http.Handler, error) {
	db, err := database.Open(cfg)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(3) // serverless: batasi koneksi per instance

	mux := http.NewServeMux()
	article.NewHandler(article.NewService(article.NewRepository(db))).Register(mux)
	return httpx.CORS(cfg.CORSOrigin, mux), nil
}
