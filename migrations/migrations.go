// Package migrations menyimpan file SQL migrasi yang di-embed ke dalam binary.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
