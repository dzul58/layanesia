package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

const (
	EnvDevelopment = "development"
	EnvStaging     = "staging"
	EnvProduction  = "production"
	EnvTest        = "test"
)

type Config struct {
	Port string
	Env  string

	// Database
	DBHost            string
	DBPort            string
	DBUser            string
	DBPassword        string
	DBName            string
	DBSSLMode         string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBRunMigrations   bool
	DBLogQueries      bool

	// Auth
	JWTSecret   string
	JWTTTL      time.Duration
	AdminAPIKey string

	// HTTP
	CORSAllowOrigins string
	BodyLimitMB      int
	TrustedProxies   []string

	// Midtrans
	MidtransServerKey    string
	MidtransClientKey    string
	MidtransIsProduction bool

	// Cloudinary
	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
	CloudinaryFolder    string

	// Upload
	UploadLocalDir      string
	UploadMaxImageMB    int
	UploadMaxDocumentMB int

	// KTP verification
	KTPAutoVerify bool

	// Feature flags
	EnablePaymentSimulation bool
}

var AppConfig Config

func (c Config) IsProduction() bool { return c.Env == EnvProduction }
func (c Config) IsTest() bool       { return c.Env == EnvTest }

func (c Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Jakarta",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func (c Config) CloudinaryConfigured() bool {
	return c.CloudinaryCloudName != "" && c.CloudinaryAPIKey != "" && c.CloudinaryAPISecret != ""
}

// LoadConfig reads configuration from environment variables (and .env for local
// development). It returns an error listing every missing/invalid variable so that
// production deployments fail fast instead of silently running with insecure defaults.
func LoadConfig() (Config, error) {
	_ = godotenv.Load()

	env := strings.ToLower(getEnv("ENV", EnvDevelopment))
	isProd := env == EnvProduction

	cfg := Config{
		Port: getEnv("PORT", "8080"),
		Env:  env,

		DBHost:            getEnv("DB_HOST", "localhost"),
		DBPort:            getEnv("DB_PORT", "5432"),
		DBUser:            getEnv("DB_USER", "postgres"),
		DBPassword:        getEnv("DB_PASSWORD", ""),
		DBName:            getEnv("DB_NAME", "layanesia_db"),
		DBSSLMode:         getEnv("DB_SSLMODE", ifProd(isProd, "require", "disable")),
		DBMaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Minute),
		DBRunMigrations:   getEnvBool("DB_RUN_MIGRATIONS", !isProd),
		DBLogQueries:      getEnvBool("DB_LOG_QUERIES", !isProd),

		JWTSecret:   getEnv("JWT_SECRET", ""),
		JWTTTL:      getEnvDuration("JWT_TTL", 7*24*time.Hour),
		AdminAPIKey: getEnv("ADMIN_API_KEY", ""),

		CORSAllowOrigins: getEnv("CORS_ALLOW_ORIGINS", ifProd(isProd, "", "http://localhost:5173,http://localhost:3000")),
		BodyLimitMB:      getEnvInt("BODY_LIMIT_MB", 10),
		TrustedProxies:   splitCSV(getEnv("TRUSTED_PROXIES", "")),

		MidtransServerKey:    getEnv("MIDTRANS_SERVER_KEY", ""),
		MidtransClientKey:    getEnv("MIDTRANS_CLIENT_KEY", ""),
		MidtransIsProduction: getEnvBool("MIDTRANS_IS_PRODUCTION", false),

		CloudinaryCloudName: getEnv("CLOUDINARY_CLOUD_NAME", ""),
		CloudinaryAPIKey:    getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryAPISecret: getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryFolder:    getEnv("CLOUDINARY_FOLDER", "layanesia"),

		UploadLocalDir:      getEnv("UPLOAD_LOCAL_DIR", "./uploads"),
		UploadMaxImageMB:    getEnvInt("UPLOAD_MAX_IMAGE_MB", 3),
		UploadMaxDocumentMB: getEnvInt("UPLOAD_MAX_DOCUMENT_MB", 5),

		KTPAutoVerify: getEnvBool("KTP_AUTO_VERIFY", false),

		EnablePaymentSimulation: getEnvBool("ENABLE_PAYMENT_SIMULATION", !isProd),
	}

	var problems []string

	if cfg.JWTSecret == "" {
		if isProd {
			problems = append(problems, "JWT_SECRET wajib diisi")
		} else {
			cfg.JWTSecret = "dev-only-insecure-jwt-secret-change-me"
			log.Println("[config] WARNING: JWT_SECRET kosong, memakai secret development. JANGAN dipakai di production.")
		}
	} else if len(cfg.JWTSecret) < 32 {
		if isProd {
			problems = append(problems, "JWT_SECRET minimal 32 karakter")
		} else {
			log.Println("[config] WARNING: JWT_SECRET pendek (<32 karakter).")
		}
	}

	if cfg.DBPassword == "" {
		if isProd {
			problems = append(problems, "DB_PASSWORD wajib diisi")
		} else {
			cfg.DBPassword = "postgres"
		}
	}

	if isProd {
		if cfg.DBSSLMode == "disable" {
			problems = append(problems, "DB_SSLMODE tidak boleh 'disable' di production")
		}
		if cfg.CORSAllowOrigins == "" || cfg.CORSAllowOrigins == "*" {
			problems = append(problems, "CORS_ALLOW_ORIGINS wajib berisi daftar origin eksplisit (bukan '*')")
		}
		if cfg.MidtransServerKey == "" {
			problems = append(problems, "MIDTRANS_SERVER_KEY wajib diisi")
		}
		if !cfg.MidtransIsProduction {
			log.Println("[config] WARNING: MIDTRANS_IS_PRODUCTION=false di production; transaksi akan masuk sandbox.")
		}
		if !cfg.CloudinaryConfigured() {
			problems = append(problems, "CLOUDINARY_CLOUD_NAME / CLOUDINARY_API_KEY / CLOUDINARY_API_SECRET wajib diisi")
		}
		if cfg.AdminAPIKey == "" {
			problems = append(problems, "ADMIN_API_KEY wajib diisi (dipakai untuk approve verifikasi KTP)")
		}
		if cfg.EnablePaymentSimulation {
			problems = append(problems, "ENABLE_PAYMENT_SIMULATION harus false di production")
		}
		if cfg.KTPAutoVerify {
			problems = append(problems, "KTP_AUTO_VERIFY harus false di production")
		}
		if cfg.DBRunMigrations {
			log.Println("[config] WARNING: DB_RUN_MIGRATIONS=true di production; disarankan menjalankan migrasi sebagai step deploy terpisah (cmd/migrate).")
		}
	}

	if len(problems) > 0 {
		return cfg, fmt.Errorf("konfigurasi tidak valid:\n  - %s", strings.Join(problems, "\n  - "))
	}

	AppConfig = cfg
	return cfg, nil
}

// MustLoadConfig is LoadConfig that terminates the process on error.
func MustLoadConfig() Config {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("[config] %v", err)
	}
	return cfg
}

func ifProd(isProd bool, prodVal, devVal string) string {
	if isProd {
		return prodVal
	}
	return devVal
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if v := strings.TrimSpace(p); v != "" {
			out = append(out, v)
		}
	}
	return out
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok && strings.TrimSpace(value) != "" {
		return strings.TrimSpace(value)
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := getEnv(key, "")
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("[config] %s=%q bukan angka, memakai default %d", key, v, fallback)
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	v := strings.ToLower(getEnv(key, ""))
	switch v {
	case "":
		return fallback
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	}
	log.Printf("[config] %s=%q bukan boolean, memakai default %v", key, v, fallback)
	return fallback
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	v := getEnv(key, "")
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		log.Printf("[config] %s=%q bukan durasi (contoh: 30m, 168h), memakai default %s", key, v, fallback)
		return fallback
	}
	return d
}
