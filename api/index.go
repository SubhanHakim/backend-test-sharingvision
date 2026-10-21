// Package handler adalah entry point serverless untuk Vercel.
package handler

import (
	"net/http"
	"sync"

	"article/internal/app"
	"article/internal/config"
)

var (
	mu sync.Mutex
	h  http.Handler
)

// Handler dipanggil Vercel per request. App dirakit sekali per instance (cold start)
// dan dicoba lagi pada request berikutnya jika sebelumnya gagal.
func Handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	if h == nil {
		built, err := app.New(config.Load())
		if err != nil {
			mu.Unlock()
			http.Error(w, `{"error":"database unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		h = built
	}
	cur := h
	mu.Unlock()
	cur.ServeHTTP(w, r)
}
