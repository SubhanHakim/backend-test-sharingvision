package config

import (
	"os"
	"strings"
)

type Config struct {
	DBUser     string
	DBPass     string
	DBHost     string
	DBName     string
	DBTLS      bool
	Addr       string
	CORSOrigin string
}

func Load() Config {
	addr := env("ADDR", "")
	if addr == "" {
		if p := os.Getenv("PORT"); p != "" { // diisi otomatis oleh Render dsb.
			addr = ":" + p
		} else {
			addr = ":8082"
		}
	}
	return Config{
		DBUser:     env("DB_USER", "root"),
		DBPass:     env("DB_PASS", ""),
		DBHost:     env("DB_HOST", "127.0.0.1:3306"),
		DBName:     env("DB_NAME", "article"),
		DBTLS:      env("DB_TLS", "") == "true",
		Addr:       addr,
		CORSOrigin: env("CORS_ORIGIN", "http://localhost:5173"),
	}
}

func env(k, d string) string {
	bom := string(rune(0xFEFF))
	if v := strings.Trim(os.Getenv(k), " \t\r\n"+bom); v != "" { // spasi/BOM hasil copy-paste ikut dibuang
		return v
	}
	return d
}
