package database

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"article/internal/config"
	"article/migrations"
)

const dbName = "article"

// Open membuat database jika belum ada, menjalankan migrate, lalu mengembalikan koneksi.
func Open(c config.Config) (*sql.DB, error) {
	base := c.DBUser + ":" + c.DBPass + "@tcp(" + c.DBHost + ")/"

	boot, err := sql.Open("mysql", base)
	if err != nil {
		return nil, err
	}
	_, err = boot.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	boot.Close()
	if err != nil {
		return nil, err
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, "mysql://"+base+dbName+"?multiStatements=true")
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return sql.Open("mysql", base+dbName)
}
