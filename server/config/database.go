package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB adalah koneksi GORM global. Diisi oleh ConnectDB.
var DB *gorm.DB

// ConnectDB membuka koneksi PostgreSQL, mengatur connection pool, dan memverifikasi
// koneksi dengan ping. Mengembalikan error (bukan warning) jika gagal - server tidak
// boleh berjalan tanpa database.
//
// Skema database tidak lagi dibuat lewat GORM AutoMigrate; gunakan db.MigrateUp
// (dijalankan otomatis saat DB_RUN_MIGRATIONS=true) atau `go run ./cmd/migrate up`.
func ConnectDB(cfg Config) (*gorm.DB, error) {
	logLevel := logger.Error
	if cfg.DBLogQueries {
		logLevel = logger.Info
	}

	gormDB, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                 logger.Default.LogMode(logLevel),
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		NowFunc:                func() time.Time { return time.Now().In(jakarta()) },
	})
	if err != nil {
		return nil, fmt.Errorf("membuka koneksi database: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("mengambil *sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database %s@%s:%s/%s gagal: %w", cfg.DBUser, cfg.DBHost, cfg.DBPort, cfg.DBName, err)
	}

	log.Printf("[db] terhubung ke %s:%s/%s (sslmode=%s, pool max=%d)", cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode, cfg.DBMaxOpenConns)
	DB = gormDB
	return gormDB, nil
}

// PingDB dipakai health check untuk memastikan database benar-benar bisa dijangkau.
func PingDB(ctx context.Context, gormDB *gorm.DB) error {
	if gormDB == nil {
		return fmt.Errorf("database belum terinisialisasi")
	}
	sqlDB, err := gormDB.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func jakarta() *time.Location {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return time.FixedZone("WIB", 7*3600)
	}
	return loc
}
