package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"layanesia-server/config"
	"layanesia-server/db"
	"layanesia-server/internal/app"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	cfg := config.MustLoadConfig()

	gormDB, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("[db] %v", err)
	}
	sqlDB, _ := gormDB.DB()
	defer sqlDB.Close()

	if cfg.DBRunMigrations {
		if err := db.MigrateUp(cfg.DSN()); err != nil {
			log.Fatalf("[migrate] %v", err)
		}
		log.Println("[migrate] skema database up-to-date")
	} else {
		v, dirty, err := db.Version(cfg.DSN())
		if err != nil {
			log.Fatalf("[migrate] tidak dapat membaca versi skema: %v", err)
		}
		if dirty {
			log.Fatalf("[migrate] skema dalam keadaan dirty (versi %d); perbaiki dengan `cmd/migrate force`", v)
		}
		log.Printf("[migrate] versi skema: %d (migrasi otomatis nonaktif)", v)
	}

	server := app.New(cfg, gormDB)

	// Graceful shutdown: berhenti menerima koneksi baru, tunggu request berjalan.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("[server] sinyal berhenti diterima, menutup server...")
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if err := server.ShutdownWithContext(ctx); err != nil {
			log.Printf("[server] shutdown tidak bersih: %v", err)
		}
	}()

	log.Printf("[server] Layanesia API v%s (%s) mendengarkan di :%s", app.Version, cfg.Env, cfg.Port)
	if err := server.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("[server] %v", err)
	}
	log.Println("[server] berhenti.")
}
