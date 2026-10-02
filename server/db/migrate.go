// Package db menyediakan runner migrasi SQL yang di-embed ke dalam binary.
// Skema database HANYA dikelola melalui file di db/migrations (bukan GORM AutoMigrate),
// sehingga ERD, file migrasi, dan skema fisik selalu satu sumber kebenaran.
package db

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // driver database/sql "pgx"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// newMigrator membuka koneksi *terpisah* untuk migrasi. Driver pgx milik
// golang-migrate menutup *sql.DB yang diberikan saat Close(), sehingga koneksi
// aplikasi (GORM) tidak boleh dipakai bersama.
func newMigrator(dsn string) (*migrate.Migrate, error) {
	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("membuka koneksi migrasi: %w", err)
	}
	sqlDB.SetMaxOpenConns(2)

	src, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("membaca file migrasi embedded: %w", err)
	}
	driver, err := pgxmigrate.WithInstance(sqlDB, &pgxmigrate.Config{})
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("inisialisasi driver migrasi: %w", err)
	}
	m, err := migrate.NewWithInstance("iofs", src, "pgx5", driver)
	if err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("inisialisasi migrator: %w", err)
	}
	return m, nil
}

// MigrateUp menjalankan semua migrasi yang belum diterapkan. Aman dipanggil berulang.
func MigrateUp(dsn string) error {
	m, err := newMigrator(dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("menjalankan migrasi: %w", err)
	}
	return nil
}

// MigrateDown membatalkan N migrasi terakhir (steps <= 0 berarti semua).
func MigrateDown(dsn string, steps int) error {
	m, err := newMigrator(dsn)
	if err != nil {
		return err
	}
	defer m.Close()

	if steps <= 0 {
		err = m.Down()
	} else {
		err = m.Steps(-steps)
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("membatalkan migrasi: %w", err)
	}
	return nil
}

// MigrateForce menandai versi tertentu tanpa menjalankan SQL (untuk memulihkan state dirty).
func MigrateForce(dsn string, version int) error {
	m, err := newMigrator(dsn)
	if err != nil {
		return err
	}
	defer m.Close()
	return m.Force(version)
}

// Version mengembalikan versi migrasi saat ini dan status dirty.
func Version(dsn string) (uint, bool, error) {
	m, err := newMigrator(dsn)
	if err != nil {
		return 0, false, err
	}
	defer m.Close()

	v, dirty, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil
	}
	return v, dirty, err
}
