// Command migrate menjalankan migrasi skema database sebagai step deploy terpisah.
//
//	go run ./cmd/migrate up            # terapkan semua migrasi baru
//	go run ./cmd/migrate down [N]      # batalkan N migrasi terakhir (default 1)
//	go run ./cmd/migrate version       # tampilkan versi saat ini
//	go run ./cmd/migrate force <ver>   # pulihkan state dirty ke versi tertentu
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"layanesia-server/config"
	"layanesia-server/db"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	cfg := config.MustLoadConfig()
	dsn := cfg.DSN()

	switch os.Args[1] {
	case "up":
		if err := db.MigrateUp(dsn); err != nil {
			log.Fatalf("[migrate] %v", err)
		}
		printVersion(dsn)
	case "down":
		steps := 1
		if len(os.Args) > 2 {
			var err error
			steps, err = strconv.Atoi(os.Args[2])
			if err != nil {
				log.Fatalf("[migrate] N harus angka: %v", err)
			}
		}
		if cfg.IsProduction() && steps > 1 {
			log.Fatalf("[migrate] down lebih dari 1 langkah tidak diizinkan di production")
		}
		if err := db.MigrateDown(dsn, steps); err != nil {
			log.Fatalf("[migrate] %v", err)
		}
		printVersion(dsn)
	case "version":
		printVersion(dsn)
	case "force":
		if len(os.Args) < 3 {
			usage()
			os.Exit(2)
		}
		v, err := strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatalf("[migrate] versi harus angka: %v", err)
		}
		if err := db.MigrateForce(dsn, v); err != nil {
			log.Fatalf("[migrate] %v", err)
		}
		printVersion(dsn)
	default:
		usage()
		os.Exit(2)
	}
}

func printVersion(dsn string) {
	v, dirty, err := db.Version(dsn)
	if err != nil {
		log.Fatalf("[migrate] membaca versi: %v", err)
	}
	fmt.Printf("versi migrasi saat ini: %d (dirty=%v)\n", v, dirty)
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: migrate <up|down [N]|version|force <version>>")
}
