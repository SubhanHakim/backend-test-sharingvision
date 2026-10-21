package database

import (
	"database/sql"
	"errors"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"article/internal/config"
	"article/migrations"
)

// Open membuat database jika belum ada, menjalankan migrate, lalu mengembalikan koneksi.
func Open(c config.Config) (*sql.DB, error) {
	base := c.DBUser + ":" + c.DBPass + "@tcp(" + c.DBHost + ")/"
	tls := ""
	if c.DBTLS {
		tls = "&tls=true"
	}

	// MySQL cloud biasanya melarang CREATE DATABASE; abaikan jika gagal karena DB sudah disediakan.
	if boot, err := sql.Open("mysql", base+"?"+tls); err == nil {
		if _, err := boot.Exec("CREATE DATABASE IF NOT EXISTS `" + c.DBName + "`"); err != nil {
			log.Println("skip create database:", err)
		}
		// migrate memakai isolation SERIALIZABLE yang ditolak TiDB kecuali pengecekannya dimatikan.
		var version string
		if boot.QueryRow("SELECT VERSION()").Scan(&version) == nil && strings.Contains(version, "TiDB") {
			tls += "&tidb_skip_isolation_level_check=1"
		}
		boot.Close()
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, err
	}
	m, err := migrate.NewWithSourceInstance("iofs", src, "mysql://"+base+c.DBName+"?multiStatements=true"+tls)
	if err != nil {
		return nil, err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return sql.Open("mysql", base+c.DBName+"?"+tls)
}
