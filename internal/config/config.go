package config

import "os"

type Config struct {
	DBUser     string
	DBPass     string
	DBHost     string
	Addr       string
	CORSOrigin string
}

func Load() Config {
	return Config{
		DBUser:     env("DB_USER", "root"),
		DBPass:     env("DB_PASS", ""),
		DBHost:     env("DB_HOST", "127.0.0.1:3306"),
		Addr:       env("ADDR", ":8081"),
		CORSOrigin: env("CORS_ORIGIN", "http://localhost:5173"),
	}
}

func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
