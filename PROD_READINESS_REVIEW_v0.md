# Production Readiness Review - Layanesia

**Proyek:** Layanesia (On-Demand Gig Marketplace)
**Scope Review:** `server/` (Go Fiber + GORM + PostgreSQL), `client/` (React + Vite), skema database aktual `layanesia_db`, file migrasi & seed, dokumen PRD/ERD/Sprint Planning.
**Tanggal Review:** 28 September 2026
**Referensi:** [PRD_v0.md](./PRD_v0.md) | [ERD_v0.md](./ERD_v0.md) | [SPRINT_PLANNING_LAYANESIA_v0.md](./SPRINT_PLANNING_LAYANESIA_v0.md)

---

## 1. Kesimpulan Eksekutif

**Verdict: 🔴 NO-GO. Belum layak naik production.**

Kode berhasil di-compile (`go build`, `go vet` bersih) dan frontend berhasil di-build, tetapi hasil pengujian end-to-end terhadap server yang benar-benar dijalankan menunjukkan **fitur inti Model 1 dan Model 2 tidak berfungsi sama sekali** (14 endpoint terproteksi selalu mengembalikan `401`), **monetisasi Midtrans efektif tidak ada** (langganan aktif tanpa pembayaran), dan **data pribadi pengguna (email, HP, nomor KTP, alamat) bocor ke publik tanpa login**.

Ringkasan jumlah temuan:

| Severity | Jumlah | Keterangan |
| :--- | :---: | :--- |
| 🔴 P0 - Blocker | 6 | Wajib diperbaiki sebelum rilis; aplikasi tidak berfungsi / kerugian finansial / pelanggaran privasi |
| 🟠 P1 - High | 9 | Wajib diperbaiki sebelum rilis; risiko keamanan & integritas data |
| 🟡 P2 - Medium | 10 | Sangat disarankan sebelum rilis; perbaiki paling lambat sprint pertama pasca-rilis |
| ⚪ P3 - Low / Hygiene | 8 | Kualitas kode & operasional |

Estimasi kasar effort untuk mencapai status GO: **2-3 sprint** (P0 + P1 + infrastruktur deploy minimum).

---

## 2. Metodologi & Bukti Verifikasi

Review dilakukan dengan kombinasi pembacaan kode menyeluruh dan pengujian empiris:

1. `go build ./...` dan `go vet ./...` pada `server/` → **lolos** (namun tidak ada satu pun file `_test.go`).
2. `npx oxlint` dan `npx vite build` pada `client/` → build lolos, 10 warning lint.
3. Inspeksi langsung skema PostgreSQL lokal `layanesia_db` (`\d`, `pg_constraint`, `pg_indexes`) dan dibandingkan dengan model GORM, file `db/migrations/*.sql`, dan ERD.
4. Menjalankan server aktual (`PORT=8099 go run ./cmd/main.go`) lalu memanggil 16 skenario API dengan `curl` menggunakan token JWT valid dari user seed. Semua data uji yang tercipta sudah dibersihkan dan data seed dikembalikan ke kondisi semula.

Hasil pengujian yang menjadi dasar temuan P0 ditampilkan pada bagian masing-masing.

---

## 3. 🔴 P0 - Blocker

### P0-1. Key `Locals` tidak konsisten → 14 endpoint inti selalu `401 Unauthorized`

**Lokasi:** `server/internal/middlewares/auth_middleware.go` vs 6 controller.

Middleware menyimpan identitas user sebagai:

```go
c.Locals("user_id", claims.UserID)   // key "user_id", tipe uuid.UUID
```

Namun `skill_posting_controller.go`, `job_posting_controller.go`, `job_application_controller.go`, `job_offer_controller.go`, `connection_controller.go`, dan `rating_controller.go` membacanya sebagai:

```go
userIDStr, ok := c.Locals("userID").(string)   // key "userID", tipe string
```

Key berbeda **dan** tipe berbeda, sehingga type assertion selalu gagal dan handler langsung mengembalikan `401 "Tidak terautentikasi"` sebelum menyentuh service. Compiler tidak mendeteksi ini karena `Locals` bertipe `interface{}`.

**Bukti (token valid, user seed `budi.worker@gmail.com`):**

```text
GET  /api/v1/users/me            → 200 OK        (memakai "user_id" - benar)
GET  /api/v1/skills/my-postings  → 401 {"error":"Tidak terautentikasi"}
GET  /api/v1/offers/received     → 401 {"error":"Tidak terautentikasi"}
GET  /api/v1/jobs/my-applications→ 401 {"error":"Tidak terautentikasi"}
GET  /api/v1/connections/active  → 401 {"error":"Tidak terautentikasi"}
POST /api/v1/skills              → 401 {"error":"Tidak terautentikasi"}
```

**Endpoint terdampak (14):** `POST /skills`, `GET /skills/my-postings`, `POST /offers`, `GET /offers/received`, `GET /offers/sent`, `PATCH /offers/:id/respond`, `POST /jobs`, `GET /jobs/my-postings`, `GET /jobs/my-applications`, `POST /jobs/:id/apply`, `GET /jobs/:id/applicants`, `PATCH /jobs/applications/:id/accept`, `GET /connections/active`, `POST /ratings`.

Artinya seluruh flow Sprint 4, 5, dan 6 (Talent Showcase, Job Board, Match & Rating) **tidak dapat dijalankan** meskipun sprint planning menandainya ✅ Selesai dan "End-to-End Testing" tercentang.

**Perbaikan:** Buat satu helper, misalnya `middlewares.GetUserID(c *fiber.Ctx) (uuid.UUID, bool)`, dan gunakan di semua controller. Tambahkan integration test yang memanggil setiap endpoint terproteksi dengan token valid.

---

### P0-2. Langganan aktif tanpa pembayaran (monetisasi tidak berfungsi)

**Lokasi:** `server/internal/services/subscription_service.go`, `midtrans_service.go`, `cmd/main.go`.

Empat masalah yang saling menumpuk:

1. `Checkout()` langsung menyimpan subscription dengan `Status: ACTIVE` dan `ExpiresAt: +1 bulan` **sebelum** ada konfirmasi pembayaran (komentar di kode: *"Auto-active for dev / sandbox mode"*).
2. `HandleMidtransWebhook()` tidak memverifikasi `signature_key` Midtrans (SHA-512 dari `order_id + status_code + gross_amount + ServerKey`) dan tidak melakukan apa pun terhadap data; siapa pun bisa POST ke `/subscriptions/webhook`.
3. `midtransService.CreateSnapTransaction()` memakai `config.AppConfig.JWTSecret` sebagai Midtrans Server Key (bukan `MIDTRANS_SERVER_KEY`), selalu memanggil URL **sandbox**, dan bila gagal apa pun alasannya **diam-diam mengembalikan token palsu** `SNAP-SIM-xxxx` seolah sukses.
4. Endpoint `POST /subscriptions/simulate-activate` (bypass pembayaran) terdaftar tanpa guard environment.
5. `order_id` tidak disimpan di tabel `subscriptions`, sehingga webhook asli pun tidak akan bisa dipetakan ke subscription.

**Bukti:**

```text
POST /subscriptions/checkout {"plan_type":"APPLY_JOB_5K"}
→ 201 {"status":"ACTIVE","snap_token":"SNAP-SIM-3d02bfec-167", ...}   ← aktif tanpa bayar

POST /subscriptions/webhook {"order_id":"FAKE-1","transaction_status":"settlement"}
→ 200 {"status":"ok"}                                                ← tanpa signature

POST /subscriptions/simulate-activate {"plan_type":"APPLY_JOB_5K"}
→ 200 "Simulasi aktivasi paket langganan berhasil!"
```

**Perbaikan:**
- Tambah kolom `order_id` (unique), `payment_status` (`PENDING/PAID/FAILED/EXPIRED`), `snap_token`, `paid_at` pada `subscriptions`.
- `Checkout` membuat record `PENDING`; aktivasi hanya di webhook setelah verifikasi `signature_key` dan `transaction_status IN ('settlement','capture')` dengan `fraud_status = 'accept'`.
- Baca `MIDTRANS_SERVER_KEY` & `MIDTRANS_IS_PRODUCTION` dari env; pilih URL sandbox/production sesuai flag; **hapus fallback token palsu** (kembalikan error).
- Hapus `simulate-activate` dari build production atau lindungi dengan `ENV != production` **dan** admin secret.
- Frontend: integrasikan `snap.js` (`window.snap.pay(snap_token)`) — saat ini `snap_token` tidak dipakai di client sama sekali.

---

### P0-3. Kebocoran data pribadi (PII) ke publik tanpa autentikasi

**Lokasi:** Semua repository yang melakukan `Preload("User"/"Employer"/"Applicant"/"Worker")` + `models.User` yang men-serialize seluruh kolom.

Endpoint publik `GET /skills`, `GET /skills/:id`, `GET /jobs`, `GET /jobs/:id` mengembalikan objek `user`/`employer` lengkap. Model `User` hanya menyembunyikan `password` (`json:"-"`); `email`, `phone`, `ktp_number`, `ktp_image_url`, `resume_url`, `address_detail`, `latitude/longitude` semuanya ikut keluar.

**Bukti (tanpa header Authorization):**

```text
GET /api/v1/skills?limit=1 → data[0].user =
{ email: 'siti.art@gmail.com', phone: '081987654321',
  ktp_number: '3174025508920003', address_detail: 'Jl. Fatmawati Raya No. 18, Cilandak', ... }
```

Ini bertentangan langsung dengan PRD (kontak baru dibuka setelah status TERHUBUNG) dan berisiko melanggar UU PDP No. 27/2022 (nomor KTP adalah data pribadi spesifik). Hal yang sama terjadi pada `GET /jobs/:id/applicants` (employer melihat HP/email pelamar sebelum menerima) dan `GET /offers/*`.

**Perbaikan:** Buat DTO/response struct terpisah (`PublicUserDTO`: `id, name, province, city, district, is_ktp_verified, avg_rating`) dan gunakan untuk semua nested user. Kontak lengkap hanya lewat `/connections/active`. Jangan pernah mengirim `ktp_number` / `ktp_image_url` ke pihak lain.

---

### P0-4. Secrets production ter-commit & fallback secret hardcoded

**Lokasi:** `server/.env`, `server/config/config.go`, tidak ada `server/.gitignore`.

- `server/.env` berisi Midtrans Server Key, Cloudinary API Secret, dan `CLOUDINARY_URL` lengkap. Folder `server/` **tidak memiliki `.gitignore`**, sehingga file ini akan ikut ter-commit.
- `JWT_SECRET` punya fallback hardcoded `"layanesia_super_secret_jwt_key_2026"` dan `DB_PASSWORD` fallback `"postgres"`. Bila env var lupa di-set di prod, server tetap jalan dengan secret yang sudah publik → siapa pun bisa memalsukan JWT.

**Perbaikan:** Rotasi semua key yang sudah terekspos (Midtrans, Cloudinary). Tambah `server/.gitignore` (`.env`, `uploads/`, binary). Di `LoadConfig()`, `log.Fatal` bila `JWT_SECRET`/`DB_PASSWORD`/`MIDTRANS_SERVER_KEY` kosong saat `ENV=production`. Baca semua variabel Midtrans/Cloudinary ke struct `Config` (saat ini tidak dibaca sama sekali).

---

### P0-5. Upload file tanpa validasi, disimpan di disk lokal, dan disajikan statis

**Lokasi:** `user_controller.go: UploadFile`, `cmd/main.go: app.Static("/uploads", "./uploads")`.

- Tidak ada whitelist ekstensi/MIME. File `.php`, `.html`, `.svg`, `.exe` diterima dan langsung bisa diakses di `/uploads/<nama>` → stored XSS / hosting file berbahaya di domain aplikasi.
- Batas hanya `BodyLimit` global 20 MB.
- Disimpan di `./uploads` lokal → hilang saat redeploy container / tidak bisa multi-instance. PRD mensyaratkan **Cloudinary**; env Cloudinary ada tetapi tidak ada satu baris kode pun yang memakainya.
- Foto KTP (data sensitif) disajikan publik tanpa autentikasi jika URL diketahui.

**Bukti:**

```text
POST /users/upload-file  (file=evil.php)
→ 200 {"file_url":"/uploads/1790566100944192000_41323d63.php"}
```

**Perbaikan:** Validasi MIME dengan sniffing (`image/jpeg`, `image/png`, `application/pdf`), batas ukuran per-jenis (misal 2 MB KTP, 5 MB PDF), upload ke Cloudinary dengan `type: authenticated`/signed URL untuk KTP, hapus `app.Static("/uploads")`.

---

### P0-6. Registrasi gagal pada database yang ada (schema drift)

**Lokasi:** Skema aktual `layanesia_db.users` vs `models.User`.

Tabel `users` di database saat ini memiliki **dua kolom**: `password_hash VARCHAR(255) NOT NULL` (legacy) dan `password VARCHAR(255) NOT NULL` (yang dipakai model). GORM hanya mengisi `password`, sehingga:

```text
POST /api/v1/auth/register → 400
{"error":"ERROR: null value in column \"password_hash\" of relation \"users\" violates not-null constraint (SQLSTATE 23502)"}
```

Seeder Go (`cmd/seed/main.go`) bahkan punya raw SQL khusus untuk mengisi `password_hash` sebagai workaround. Pada DB prod yang fresh masalah ini tidak akan muncul, tetapi ini menunjukkan tidak ada disiplin migrasi (lihat P1-1) dan pesan error SQL mentah bocor ke klien (lihat P2-2).

**Perbaikan:** Migrasi eksplisit `ALTER TABLE users DROP COLUMN password_hash`; hapus workaround di seeder.

---

## 4. 🟠 P1 - High

### P1-1. Dua sumber kebenaran skema: `db/migrations/*.sql` vs GORM `AutoMigrate`

File SQL di `db/migrations/` tidak pernah dijalankan oleh tool apa pun (tidak ada `golang-migrate`/`goose` di `go.mod`). Skema aktual dibuat oleh `AutoMigrate` saat startup. Akibatnya:

| Aspek | File migrasi SQL | Skema aktual (AutoMigrate) |
| :--- | :--- | :--- |
| FK `ON DELETE CASCADE` | Ada di semua FK | **Tidak ada** (semua FK tanpa CASCADE) |
| `CHECK (rating_stars BETWEEN 1 AND 5)` | Ada | **Tidak ada** |
| `uuid_generate_v4()` default | Ada | Tidak ada (ID diisi aplikasi - OK) |
| Kolom `password_hash` | Tidak ada | **Ada** (legacy) |

`AutoMigrate` di production juga berisiko (hanya menambah, tidak pernah drop/rename, tidak bisa rollback, dan berjalan dengan hak DDL pada setiap start).

**Perbaikan:** Pilih satu: gunakan `golang-migrate` dengan file SQL bernomor (disarankan), jalankan sebagai step deploy terpisah, dan **matikan `AutoMigrate` saat `ENV=production`**. Tambahkan migrasi untuk CASCADE dan CHECK constraint.

### P1-2. Verifikasi KTP hanya klaim mandiri

`POST /users/verify-ktp` langsung men-set `is_ktp_verified = true` dari input user tanpa review admin, OCR, atau validasi format NIK (16 digit). Padahal KTP verified adalah gerbang keamanan untuk melamar (PRD §3). Bukti: mengirim `ktp_number: "0000"` → `200 "KTP Anda berhasil diverifikasi!"`.

**Perbaikan:** Status `ktp_status: UNVERIFIED / PENDING_REVIEW / VERIFIED / REJECTED`; endpoint admin untuk approve; minimal validasi NIK 16 digit numerik.

### P1-3. Tidak ada validasi input dasar

`RegisterDTO` tidak memvalidasi format email, panjang/kekuatan password, format nomor HP, dan `active_mode` (string apa pun tersimpan, misal `"ADMIN"`). `UpdateProfile` mengizinkan ganti `phone` tanpa cek unik → error unik DB mentah ke klien. Rekomendasi: `go-playground/validator` dengan tag pada DTO.

### P1-4. Tidak ada rate limiting

`/auth/login`, `/auth/register`, `/subscriptions/webhook`, `/locations/nominatim/search` terbuka tanpa limiter → brute force kredensial dan penyalahgunaan proxy Nominatim (kebijakan Nominatim maks 1 req/detik; IP server bisa di-ban). Gunakan `fiber/middleware/limiter`.

### P1-5. Non-atomic state transition & error diabaikan

`AcceptApplicant` dan `RespondToOffer`: `UpdateStatus` lalu `CreateConnection` dijalankan tanpa transaksi dan error koneksi diabaikan (`_ = s.connRepo.Create(conn)`). Jika create koneksi gagal, lamaran sudah `ACCEPTED` tanpa koneksi → user tidak pernah "TERHUBUNG" dan tidak bisa rating. Bungkus dalam `db.Transaction`.

### P1-6. Unique constraint `job_offers` tidak sesuai ERD

ERD §2.6 mendefinisikan `UNIQUE(skill_posting_id, employer_id) WHERE status = 'PENDING'` (partial index). Implementasi memakai unique index penuh, sehingga employer **tidak akan pernah bisa** menawari worker yang sama lagi setelah tawaran pertama ditolak/selesai. Ganti dengan partial unique index via migrasi SQL (GORM tag tidak mendukung partial index secara portabel).

### P1-7. Lifecycle status tidak lengkap

`job_postings.status` tidak pernah berubah dari `OPEN` (tidak ada transisi ke `IN_PROGRESS`/`DONE`), `connections.status` tidak pernah `COMPLETED`, `subscriptions.status` tidak pernah `EXPIRED`. Akibat: lowongan yang sudah terisi tetap tampil & bisa dilamar; rating bisa diberikan sebelum pekerjaan selesai; riwayat langganan menampilkan `ACTIVE` untuk paket yang sudah lewat. Butuh endpoint `PATCH /jobs/:id/close`, `PATCH /connections/:id/complete`, dan job/cron expiry (atau hitung status saat baca).

### P1-8. CORS `AllowOrigins: "*"` + `DB_SSLMODE=disable`

Batasi origin ke domain frontend production dan wajibkan `sslmode=require` untuk koneksi DB di prod.

### P1-9. `CreateOffer` tidak memeriksa KTP employer

PRD §3 mewajibkan KTP verified untuk "melamar **& membuat tawaran**". `ApplyJob` sudah mengecek, `CreateOffer` belum.

---

## 5. 🟡 P2 - Medium

1. **Server tetap hidup saat DB gagal konek.** `ConnectDB` hanya `log.Printf` lalu `return`; `config.DB == nil` → setiap request panic (tertangkap `recover`, tapi respons 500 tidak informatif). Health check tetap menjawab `"database": "PostgreSQL GORM Connected"` tanpa ping. Gunakan `log.Fatal` saat gagal konek dan `DB.Exec("SELECT 1")` di `/health`.
2. **Error internal bocor ke klien.** Controller mengembalikan `err.Error()` langsung, termasuk pesan SQLSTATE (lihat P0-6). Pisahkan error domain (dikirim) vs error internal (log + pesan generik).
3. **GORM logger `logger.Info`** mencetak semua SQL beserta parameter (termasuk email, HP, hash password) ke stdout di semua environment. Set `Silent`/`Error` di production.
4. **JWT custom tanpa refresh/revocation**, umur 7 hari, `active_mode` di-embed dalam token sehingga basi setelah `switch-mode` (frontend sudah menangani dengan token baru; backend tidak memakai claim ini). Pertimbangkan library standar (`golang-jwt/jwt/v5`) + `iat`, `nbf`, `jti`.
5. **Race condition pada `Checkout`**: cek `FindActiveByUserAndPlan` lalu `Create` tanpa lock → dua request paralel menghasilkan dua langganan. Tambahkan partial unique index `(user_id, plan_type) WHERE status='ACTIVE'` atau `SELECT ... FOR UPDATE`.
6. **`work_date` gagal parse → diam-diam default besok**, tidak ada validasi tanggal masa lalu. Kembalikan `400`.
7. **Search `ILIKE '%term%'`** pada `title/description/province/city/district` tanpa index yang bisa dipakai (btree tidak membantu leading wildcard). Untuk MVP masih OK; siapkan `pg_trgm` GIN index atau paling tidak filter equality untuk `province/city/district` (nilainya berasal dari master data, seharusnya exact match, bukan `ILIKE %..%`).
8. **Tidak ada index komposit** untuk query terpanas: `subscriptions(user_id, plan_type, status, expires_at)`, `job_postings(status, province, city, district, created_at)`, `skill_postings(availability, province, city, district)`.
9. **`locations` tanpa `UNIQUE(province, city_or_regency, district)`**; idempotensi seed bergantung pada `FirstOrCreate` di aplikasi. Data seed aktual 229 baris, sesuai dengan file SQL.
10. **Seed SQL `db/seeds/000002_seed_users.sql` menyimpan password plaintext** (`'Password123!'`) → login user seed gagal jika seed SQL dipakai. Hanya seeder Go yang benar. Hapus atau ganti dengan hash bcrypt.

---

## 6. ⚪ P3 - Low / Hygiene

1. **0 unit test, 0 integration test**, tidak ada CI. Bug P0-1 akan langsung tertangkap oleh satu integration test sederhana.
2. Tidak ada `Dockerfile`, `docker-compose`, README server, atau dokumentasi deploy; tidak ada graceful shutdown (`app.ShutdownWithContext`) dan konfigurasi connection pool (`SetMaxOpenConns`, `SetConnMaxLifetime`).
3. `golang.org/x/crypto v0.14.0` (Okt 2023) sudah sangat tertinggal; jalankan `govulncheck ./...` dan `go get -u` untuk dependensi keamanan.
4. Logging belum terstruktur & tanpa request ID; sulit tracing di production.
5. Nama constraint di DB menggunakan prefiks `idx_*` untuk UNIQUE constraint (`idx_job_applicant`, `idx_rating_conn_rev`) — membingungkan, sebaiknya `uq_*`.
6. `ratingService` mendefinisikan `type ratingRepoAlias = repositories.RatingRepository` tanpa alasan; `GetConnectionByID` tidak dipakai; `ConnectionRepository.Create` dan `JobOfferRepository.CreateConnection` duplikat logika.
7. `README.md` client masih template Vite default.
8. Health endpoint dan root route mengembalikan `version: "1.0.0"` hardcoded; sebaiknya dari build flag (`-ldflags -X`).

---

## 7. Review Database (Skema Aktual vs ERD vs Migrasi)

**Yang sudah sesuai dan baik:**
- Sembilan tabel MVP sesuai ERD; `booster_blasts` benar di-exclude.
- Semua PK UUID, FK terpasang ke `users`/`job_postings`/`skill_postings`/`connections`.
- Unique constraint sesuai ERD ada untuk `job_applications(job_posting_id, applicant_id)`, `connections(source_id, employer_id, worker_id)`, `ratings(connection_id, reviewer_id)`, `users(email)`, `users(phone)`.
- Tipe `NUMERIC(12,2)` untuk uang, `DECIMAL(10,8)/(11,8)` untuk koordinat, `TIMESTAMPTZ` dengan TimeZone Asia/Jakarta — tepat.
- Index pada FK `user_id`, `employer_id`, `worker_id`, `reviewee_id`, dan kolom lokasi `locations`.

**Yang menyimpang (ringkasan dari temuan di atas):**

| # | Item | ERD / Migrasi SQL | Skema aktual | Ref |
| :---: | :--- | :--- | :--- | :---: |
| 1 | Kolom `users.password_hash` | Tidak ada | Ada, `NOT NULL` → registrasi gagal | P0-6 |
| 2 | FK `ON DELETE CASCADE` | Ada | Tidak ada | P1-1 |
| 3 | `CHECK rating_stars 1..5` | Ada | Tidak ada (hanya validasi aplikasi) | P1-1 |
| 4 | `job_offers` unique | Partial `WHERE status='PENDING'` | Full unique | P1-6 |
| 5 | `subscriptions` kolom pembayaran | `order_id`, `payment_status` tidak ada di ERD | Tidak ada | P0-2 |
| 6 | Mekanisme migrasi | File SQL bernomor | AutoMigrate saat startup, SQL tidak dipakai | P1-1 |
| 7 | Index komposit pencarian | - | Tidak ada | P2-8 |

**Catatan arsitektur yang sudah diakui ERD dan dapat diterima untuk MVP:** polymorphic `connections.source_type/source_id` tanpa FK native, serta lokasi denormalized string pada `users/skill_postings/job_postings`.

---

## 8. Review Frontend (`client/`)

**Blocker untuk deploy:**
1. **Base URL API hardcoded `http://localhost:8080`** di 15 file (`AuthContext.jsx`, `Home.jsx`, `LocationSelector.jsx`, semua modal). Proxy Vite di `vite.config.js` hanya berlaku saat `dev`. Build production akan memanggil localhost user. Ganti dengan `import.meta.env.VITE_API_BASE_URL` di satu modul `api.js`.
2. **Snap.js Midtrans tidak diintegrasikan.** `SubscriptionModal` memanggil checkout lalu langsung menampilkan "Transaksi Berhasil!" tanpa membuka popup pembayaran; `snap_token`/`redirect_url` tidak dipakai (konsisten dengan P0-2).

**Perlu diperbaiki:**
3. `LocationPickerModal.jsx` fallback memanggil `nominatim.openstreetmap.org` **langsung dari browser** tanpa header `User-Agent` khusus (browser tidak mengizinkan set UA) → melanggar usage policy Nominatim.
4. Tidak ada `ErrorBoundary`; satu error render mematikan seluruh SPA.
5. `react-router-dom` di-install tetapi tidak dipakai (single page `Home`); hapus dependensi atau gunakan.
6. Bundle utama 502 KB (147 KB gzip) dalam satu chunk; lazy-load `react-leaflet` dan modal.
7. 10 warning lint: unused imports (`Home.jsx`, `PostSkillModal.jsx`), `setState` sinkron dalam effect dan `exhaustive-deps` di `Home.jsx:198-200`.
8. Token JWT di `localStorage` (rentan XSS; relevan karena P0-5 memungkinkan stored XSS). Setelah P0-5 ditutup, risiko menurun; opsi lanjutan: cookie `HttpOnly` + CSRF token.

**Yang sudah baik:** `safeParseJSON` untuk respons non-JSON, penanganan `401` → auto logout, refresh token saat `switch-mode`, pemisahan komponen per fitur cukup rapi, theme persist.

---

## 9. Checklist Go-Live (urutan yang disarankan)

**Fase A - Fungsional (tanpa ini aplikasi tidak bisa dipakai)**
- [ ] P0-1 Helper `GetUserID` + refactor 6 controller + integration test seluruh endpoint terproteksi
- [ ] P0-6 / P1-1 Adopsi `golang-migrate`, migrasi drop `password_hash`, tambah CASCADE & CHECK, matikan AutoMigrate di prod
- [ ] Frontend: `VITE_API_BASE_URL` terpusat

**Fase B - Keamanan & Privasi**
- [ ] P0-3 DTO publik untuk user; audit semua `Preload`
- [ ] P0-4 Rotasi key, `server/.gitignore`, fail-fast config di prod
- [ ] P0-5 Upload ke Cloudinary + validasi MIME/ukuran; hapus static `/uploads`
- [ ] P1-2 Alur verifikasi KTP dengan review
- [ ] P1-3 Validator DTO; P1-4 rate limiter; P1-8 CORS & SSL; P1-9 cek KTP di `CreateOffer`

**Fase C - Monetisasi**
- [ ] P0-2 Skema pembayaran (`order_id`, `payment_status`), webhook dengan verifikasi signature, hapus fallback token palsu & `simulate-activate`, env Midtrans, Snap.js di frontend
- [ ] P2-5 Guard duplikasi langganan

**Fase D - Integritas & Operasional**
- [ ] P1-5 Transaksi DB pada accept/respond; P1-6 partial unique `job_offers`; P1-7 lifecycle status
- [ ] P2-1..P2-3 Fail-fast DB, error handling terpusat, log level prod
- [ ] Dockerfile, health check dengan DB ping, graceful shutdown, structured logging, CI (`go test`, `go vet`, `govulncheck`, `oxlint`, `vite build`)

Setelah Fase A-C selesai dan diverifikasi ulang dengan pengujian end-to-end (register → KTP → langganan berbayar sandbox → post skill/job → offer/apply → accept → connection → rating), review ini dapat diulang untuk keputusan GO.

---

## 10. Hal yang Sudah Baik (dipertahankan)

- Struktur layer `controllers → services → repositories → models` konsisten dan mudah diuji (interface pada semua service/repo).
- Bcrypt untuk password dengan hook `BeforeCreate`; perbandingan hash konstan waktu; pesan login generik ("email atau password tidak sesuai").
- Query GORM parametrik (tidak ditemukan SQL injection); UUID sebagai identifier publik.
- Otorisasi kepemilikan sudah dicek di service (`job.EmployerID != employerID`, `offer.WorkerID != workerID`, `conn.EmployerID/WorkerID == reviewerID`).
- Aturan bisnis PRD sebagian besar sudah diterjemahkan: syarat KTP + Resume + langganan untuk apply, langganan untuk post skill, larangan melamar/menawar diri sendiri, anti-duplikasi lamaran/offer/rating via unique constraint, fallback lokasi dari profil user.
- Master data lokasi Jabodetabek lengkap (229 kecamatan) dan seeder Go idempoten.
- Dokumentasi PRD/ERD/Sprint dan Postman collection tersedia dan cukup detail.
