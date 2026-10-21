package article

import (
	"errors"
	"strings"
)

type Article struct {
	ID       int    `json:"id"` // hanya output; diabaikan pada create/update
	Title   string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	Status   string `json:"status"`
}

var ErrNotFound = errors.New("article not found")

// ValidationError berisi daftar pesan validasi yang gagal.
type ValidationError []string

func (e ValidationError) Error() string { return strings.Join(e, "; ") }
