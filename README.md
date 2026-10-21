# Article Service (Sharing Vision Backend Test)

Microservice Go + MySQL untuk use case Post Article.

## Menjalankan
Prasyarat: Go 1.22+ dan MySQL berjalan di `127.0.0.1:3306`.

```
go mod tidy
go run ./cmd/server
```

Saat start, program membuat database `article` lalu menjalankan migrate
(`migrations/`) untuk membuat tabel `posts`. Server berjalan di `:8082`.

Konfigurasi lewat environment variable (semua opsional):

| Variabel  | Default          |
|-----------|------------------|
| `DB_USER` | `root`           |
| `DB_PASS` | *(kosong)*       |
| `DB_HOST` | `127.0.0.1:3306` |
| `DB_NAME` | `article`        |
| `DB_TLS`  | *(kosong)*, isi `true` untuk MySQL cloud |
| `PORT`    | dipakai jika `ADDR` kosong (diisi otomatis oleh Render) |
| `ADDR`    | `:8082`          |
| `CORS_ORIGIN` | `http://localhost:5173` |

## Struktur
```
cmd/server/main.go          entry point, wiring dependensi
internal/config/            baca konfigurasi dari env
internal/database/          buat DB + jalankan migrate + koneksi
internal/article/           model, repository (SQL), service (validasi), handler (HTTP)
internal/httpx/             helper respons JSON dan middleware CORS
migrations/                 file SQL migrasi (di-embed ke binary)
```
Alur request: handler → service → repository → MySQL.

Pembuatan manual database/tabel (soal Database no. 1): jalankan `schema.sql`.

## Endpoint

| Method | URL | Keterangan |
|---|---|---|
| POST | `/article/` | Buat article |
| GET | `/article/{limit}/{offset}` | List dengan paging |
| GET | `/article/{id}` | Detail |
| PUT | `/article/{id}` | Ubah |
| DELETE | `/article/{id}` | Hapus |

Body create/update:
```json
{ "title": "", "content": "", "category": "", "status": "publish" }
```

Validasi: `title` min 20 karakter, `content` min 200, `category` min 3,
`status` salah satu dari `publish`, `draft`, `thrash`. Gagal validasi → 400.
Id tidak ditemukan → 404.

## Postman
Import `postman_collection.json` (variabel `base` = `http://localhost:8082`).

## Deploy gratis (Render + MySQL cloud)
1. Buat database MySQL gratis (mis. TiDB Cloud Serverless atau Aiven), catat host:port, user, password, nama database.
   Hindari karakter spesial (`@ / : ?`) pada password karena ikut masuk ke URL migrate.
2. Di Render: **New → Web Service**, hubungkan repo GitHub ini.
   - Runtime: `Go`
   - Build Command: `go build -o app ./cmd/server`
   - Start Command: `./app`
3. Environment variables:
   `DB_HOST` (host:port), `DB_USER`, `DB_PASS`, `DB_NAME`, `DB_TLS=true`, `CORS_ORIGIN` (URL FE).
4. Deploy. Tabel dibuat otomatis oleh migrate saat start. Base URL = URL service Render.
