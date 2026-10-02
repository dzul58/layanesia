// Package app merakit seluruh dependency (repo -> service -> controller) dan
// mendaftarkan route. Dipisahkan dari main agar bisa dipakai integration test.
package app

import (
	"log"
	"strings"
	"time"

	"layanesia-server/config"
	"layanesia-server/internal/apperrors"
	"layanesia-server/internal/controllers"
	"layanesia-server/internal/middlewares"
	"layanesia-server/internal/repositories"
	"layanesia-server/internal/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"gorm.io/gorm"
)

const Version = "1.1.0"

// Options memungkinkan test menyuntikkan implementasi eksternal (Midtrans, upload).
type Options struct {
	Midtrans services.MidtransService
	Upload   services.UploadService
}

func New(cfg config.Config, db *gorm.DB, opts ...Options) *fiber.App {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}

	app := fiber.New(fiber.Config{
		AppName:                 "Layanesia API",
		ServerHeader:            "",
		DisableStartupMessage:   cfg.IsProduction() || cfg.IsTest(),
		BodyLimit:               cfg.BodyLimitMB * 1024 * 1024,
		ReadTimeout:             15 * time.Second,
		WriteTimeout:            30 * time.Second,
		IdleTimeout:             60 * time.Second,
		ErrorHandler:            middlewares.ErrorHandler,
		EnableTrustedProxyCheck: len(cfg.TrustedProxies) > 0,
		TrustedProxies:          cfg.TrustedProxies,
		ProxyHeader:             proxyHeader(cfg),
	})

	// --- Global middlewares ---
	app.Use(recover.New(recover.Config{EnableStackTrace: !cfg.IsProduction()}))
	app.Use(requestid.New())
	app.Use(helmet.New())
	if !cfg.IsTest() {
		app.Use(logger.New(logger.Config{
			Format:     "${time} ${locals:requestid} ${ip} ${status} ${latency} ${method} ${path} ${error}\n",
			TimeFormat: time.RFC3339,
		}))
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins(cfg),
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-Admin-Key",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
		AllowCredentials: false,
		MaxAge:           600,
	}))

	// --- Dependencies ---
	locationRepo := repositories.NewLocationRepository(db)
	userRepo := repositories.NewUserRepository(db)
	subRepo := repositories.NewSubscriptionRepository(db)
	skillRepo := repositories.NewSkillPostingRepository(db)
	offerRepo := repositories.NewJobOfferRepository(db)
	jobRepo := repositories.NewJobPostingRepository(db)
	appRepo := repositories.NewJobApplicationRepository(db)
	connRepo := repositories.NewConnectionRepository(db)
	ratingRepo := repositories.NewRatingRepository(db)

	midtransSvc := opt.Midtrans
	if midtransSvc == nil {
		midtransSvc = services.NewMidtransService(services.MidtransConfig{
			ServerKey:    cfg.MidtransServerKey,
			IsProduction: cfg.MidtransIsProduction,
		})
	}
	uploadSvc := opt.Upload
	if uploadSvc == nil {
		uploadSvc = services.NewUploadService(services.UploadConfig{
			CloudName:     cfg.CloudinaryCloudName,
			APIKey:        cfg.CloudinaryAPIKey,
			APISecret:     cfg.CloudinaryAPISecret,
			Folder:        cfg.CloudinaryFolder,
			LocalDir:      cfg.UploadLocalDir,
			MaxImageBytes: int64(cfg.UploadMaxImageMB) * 1024 * 1024,
			MaxDocBytes:   int64(cfg.UploadMaxDocumentMB) * 1024 * 1024,
			AllowLocal:    !cfg.IsProduction(),
		})
	}

	locationSvc := services.NewLocationService(locationRepo)
	authSvc := services.NewAuthService(userRepo)
	userSvc := services.NewUserService(userRepo, cfg.KTPAutoVerify)
	subSvc := services.NewSubscriptionService(db, subRepo, userRepo, midtransSvc, cfg.EnablePaymentSimulation)
	skillSvc := services.NewSkillPostingService(skillRepo, userRepo, subRepo)
	offerSvc := services.NewJobOfferService(db, offerRepo, skillRepo, userRepo, connRepo)
	jobSvc := services.NewJobPostingService(db, jobRepo, appRepo, userRepo)
	appSvc := services.NewJobApplicationService(db, appRepo, jobRepo, userRepo, subRepo, connRepo)
	connSvc := services.NewConnectionService(db, connRepo, jobRepo, appRepo)
	ratingSvc := services.NewRatingService(ratingRepo, connRepo)

	locationCtrl := controllers.NewLocationController(locationSvc)
	authCtrl := controllers.NewAuthController(authSvc)
	userCtrl := controllers.NewUserController(userSvc, uploadSvc)
	subCtrl := controllers.NewSubscriptionController(subSvc)
	skillCtrl := controllers.NewSkillPostingController(skillSvc)
	offerCtrl := controllers.NewJobOfferController(offerSvc)
	jobCtrl := controllers.NewJobPostingController(jobSvc)
	appCtrl := controllers.NewJobApplicationController(appSvc)
	connCtrl := controllers.NewConnectionController(connSvc)
	ratingCtrl := controllers.NewRatingController(ratingSvc)
	healthCtrl := controllers.NewHealthController(db, Version)

	// File lokal hanya di dev (Cloudinary tidak dikonfigurasi). Di production
	// Cloudinary wajib (divalidasi config) sehingga route ini tidak pernah aktif.
	if uploadSvc.UsesLocalStorage() {
		app.Static("/uploads", uploadSvc.LocalDir(), fiber.Static{Browse: false, MaxAge: 3600})
		log.Printf("[app] penyimpanan file lokal aktif di %s (mode pengembangan)", uploadSvc.LocalDir())
	}

	// --- Rate limiters ---
	authLimiter := limiter.New(limiter.Config{
		Max:        20,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return apperrors.TooManyRequests("Terlalu banyak percobaan. Coba lagi dalam satu menit.")
		},
	})
	writeLimiter := limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if id, ok := middlewares.GetUserID(c); ok {
				return id.String()
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return apperrors.TooManyRequests("Terlalu banyak permintaan. Coba lagi sebentar.")
		},
	})
	externalLimiter := limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return apperrors.TooManyRequests("Terlalu banyak pencarian lokasi. Coba lagi sebentar.")
		},
	})

	protected := middlewares.Protected()
	admin := middlewares.AdminOnly(cfg.AdminAPIKey)

	// --- Routes ---
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"app": "Layanesia Backend API", "version": Version, "health": "/api/v1/health"})
	})

	api := app.Group("/api/v1")
	api.Get("/health", healthCtrl.CheckHealth)
	api.Get("/health/live", healthCtrl.Live)

	auth := api.Group("/auth", authLimiter)
	auth.Post("/register", authCtrl.Register)
	auth.Post("/login", authCtrl.Login)

	users := api.Group("/users", protected)
	users.Get("/me", userCtrl.GetMe)
	users.Put("/profile", userCtrl.UpdateProfile)
	users.Patch("/switch-mode", userCtrl.SwitchMode)
	users.Post("/verify-ktp", writeLimiter, userCtrl.VerifyKTP)
	users.Post("/upload-resume", userCtrl.UploadResume)
	users.Post("/upload-file", writeLimiter, userCtrl.UploadFile)

	adminGroup := api.Group("/admin", admin)
	adminGroup.Get("/ktp/pending", userCtrl.AdminListPendingKTP)
	adminGroup.Post("/users/:id/ktp/review", userCtrl.AdminReviewKTP)

	subs := api.Group("/subscriptions")
	subs.Post("/webhook", externalLimiter, subCtrl.HandleWebhook) // dipanggil Midtrans, tanpa JWT
	subs.Post("/checkout", protected, writeLimiter, subCtrl.Checkout)
	subs.Get("/my", protected, subCtrl.GetMySubscriptions)
	subs.Get("/active", protected, subCtrl.CheckActiveStatus)
	subs.Post("/:id/sync", protected, writeLimiter, subCtrl.SyncStatus)
	if cfg.EnablePaymentSimulation && !cfg.IsProduction() {
		subs.Post("/simulate-activate", protected, subCtrl.SimulateActivate)
		log.Println("[app] WARNING: endpoint /subscriptions/simulate-activate aktif (mode pengembangan)")
	}

	skills := api.Group("/skills")
	skills.Get("/", skillCtrl.Search)
	skills.Get("/my-postings", protected, skillCtrl.GetMyPostings)
	skills.Get("/:id", skillCtrl.GetByID)
	skills.Post("/", protected, writeLimiter, skillCtrl.Create)
	skills.Patch("/:id/availability", protected, skillCtrl.UpdateAvailability)
	skills.Delete("/:id", protected, skillCtrl.Delete)

	offers := api.Group("/offers", protected)
	offers.Post("/", writeLimiter, offerCtrl.Create)
	offers.Get("/received", offerCtrl.GetReceived)
	offers.Get("/sent", offerCtrl.GetSent)
	offers.Patch("/:id/respond", offerCtrl.Respond)

	jobs := api.Group("/jobs")
	jobs.Get("/", jobCtrl.Search)
	jobs.Get("/my-postings", protected, jobCtrl.GetMyPostings)
	jobs.Get("/my-applications", protected, appCtrl.GetMyApplications)
	jobs.Patch("/applications/:id/accept", protected, appCtrl.Accept)
	jobs.Patch("/applications/:id/reject", protected, appCtrl.Reject)
	jobs.Get("/:id", jobCtrl.GetByID)
	jobs.Post("/", protected, writeLimiter, jobCtrl.Create)
	jobs.Post("/:id/apply", protected, writeLimiter, appCtrl.Apply)
	jobs.Get("/:id/applicants", protected, appCtrl.GetApplicants)
	jobs.Patch("/:id/close", protected, jobCtrl.Close)

	conns := api.Group("/connections", protected)
	conns.Get("/", connCtrl.List)
	conns.Get("/active", connCtrl.GetActive)
	conns.Get("/:id", connCtrl.GetByID)
	conns.Patch("/:id/complete", connCtrl.Complete)

	ratings := api.Group("/ratings")
	ratings.Post("/", protected, writeLimiter, ratingCtrl.Create)
	ratings.Get("/user/:userID", ratingCtrl.GetUserRatings)

	locs := api.Group("/locations")
	locs.Get("/provinces", locationCtrl.GetProvinces)
	locs.Get("/cities", locationCtrl.GetCities)
	locs.Get("/districts", locationCtrl.GetDistricts)
	locs.Get("/search", locationCtrl.SearchLocations)
	locs.Get("/nominatim/search", externalLimiter, locationCtrl.GeocodeNominatim)

	// 404 JSON konsisten.
	app.Use(func(c *fiber.Ctx) error {
		return fiber.NewError(fiber.StatusNotFound, "Endpoint tidak ditemukan.")
	})

	return app
}

func corsOrigins(cfg config.Config) string {
	if strings.TrimSpace(cfg.CORSAllowOrigins) == "" {
		// Hanya terjadi di non-prod (prod divalidasi config). Default aman: tidak ada origin.
		return "http://localhost:5173"
	}
	return cfg.CORSAllowOrigins
}

func proxyHeader(cfg config.Config) string {
	if len(cfg.TrustedProxies) > 0 {
		return fiber.HeaderXForwardedFor
	}
	return ""
}
