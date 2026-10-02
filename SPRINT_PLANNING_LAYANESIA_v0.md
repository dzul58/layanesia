# Layanesia - Agile Development Roadmap & Sprint Planning

**Nama Proyek:** Layanesia  
**Platform:** On-Demand Gig & Temporary Replacement Marketplace  
**Metodologi:** Agile Scrum  
**Lokasi File:** `/Users/fiqar/Documents/agustus/x-project/layanesia/SPRINT_PLANNING_v0.md`  
**Referensi Dokumen:** [PRD_v0.md](file:///Users/fiqar/Documents/agustus/x-project/layanesia/PRD_v0.md) | [ERD_v0.md](file:///Users/fiqar/Documents/agustus/x-project/layanesia/ERD_v0.md) | [Figma Flowchart](https://www.figma.com/board/z0bEkpshXTalDn0ZtF8AK8/Untitled?node-id=0-1&t=LjyiOp8hOBdww0zA-1)


---

## 📊 Ringkasan Status Sprint MVP

| Sprint | Nama Sprint | Focus Area | Status | Output Utama |
| :--- | :--- | :--- | :--- | :--- |
| **Sprint 0** | Setup & Arsitektur dasar | Backend & Frontend | ✅ **Selesai** | Project Scaffold, GORM Migrations, React Setup |
| **Sprint 1** | Auth, Profile & Dual Role | Auth & KTP Cloudinary | ✅ **Selesai** | Login/Register (JWT & **Bcrypt**), Verifikasi KTP, Role Mode Switcher |
| **Sprint 2** | Lokasi & OpenStreetMap | Hiperlokal Location | ✅ **Selesai** | API Kecamatan/Kota, Leaflet.js Map Component & Postman Collection |
| **Sprint 3** | Monetisasi Midtrans | Midtrans Payment | ✅ **Selesai** | Package 5K (Apply) & 10K (Skill), Midtrans Webhook |
| **Sprint 4** | Model 1: Talent Showcase | Worker Skill & Offer | ✅ **Selesai** | Post Skill, Search Talent, Send Direct Offer & Accept |
| **Sprint 5** | Model 2: Job Board | Employer Job & Apply | ✅ **Selesai** | Post Job (Free), Guest Preview, Apply Job, Accept Applicant |
| **Sprint 6** | Match & Reciprocal Rating | Match & Mutual Rating | ✅ **Selesai** | Status TERHUBUNG, Unlock Kontak, Rating 2 Arah, QA |

---

## 📑 Rincian Ceklis Tugas Sprint (Checklists)

### 🟩 Sprint 0: Inisialisasi Project & Arsitektur Dasar
- [x] Inisialisasi backend Go Fiber di folder `/server`.
- [x] Setup koneksi PostgreSQL + GORM AutoMigrate untuk entitas MVP (`users`, `subscriptions`, `skill_postings`, `job_postings`, `job_applications`, `job_offers`, `connections`, `ratings`, `locations`).
  - [x] *Catatan AutoMigrate:* Exclude `booster_blasts` (🔒 Post-MVP).
  - [x] Pasang Unique DB Constraints pada `job_applications`, `ratings`, `job_offers`, dan `connections`.
- [x] Inisialisasi frontend React (Vite) di folder `/client`.
- [x] Setup struktur folder modular di backend (`controllers`, `services`, `repositories`, `models`, `middlewares`).
- [x] Setup sistem styling Vanilla/Modular CSS untuk tampilan UI modern & responsif.

---

### 🟩 Sprint 1: Autentikasi, Manajemen Profil & Single Account Dual Role
- [x] Endpoint `POST /api/v1/auth/register` (Register User & Hash Password menggunakan **Bcrypt** `golang.org/x/crypto/bcrypt`).
- [x] Endpoint `POST /api/v1/auth/login` (Login, Verifikasi Hash Bcrypt & Generate JWT Token).

- [x] Endpoint `GET /api/v1/users/me` & `PUT /api/v1/users/profile` (Update Data Diri & Lokasi Domisili: `province`, `city`, `district`).
- [x] Fitur **Single Account Dual-Role Switcher**:
  - [x] Implementasi kolom `active_mode` (`SEEKER` vs `EMPLOYER`) di tabel `users`.
  - [x] Endpoint `PATCH /api/v1/users/switch-mode` untuk berpindah mode kerja.
- [x] Fitur **Verifikasi KTP & Resume Upload**:
  - [x] Integrasi upload file multipart & URL link untuk gambar KTP & file Resume PDF.
  - [x] Endpoint `POST /api/v1/users/verify-ktp` & `POST /api/v1/users/upload-resume`.
- [x] UI React: Form Login/Register, Modal Upload KTP & Resume, serta Switcher Toggle Mode di Navbar.

---

### 🟩 Sprint 2: Master Data Lokasi & OpenStreetMap + Leaflet.js
- [x] Endpoint Master Data Lokasi: `GET /api/v1/locations/provinces`, `cities`, dan `districts` (Kecamatan).
- [x] Komponen Frontend Peta Interaktif (`react-leaflet` / OpenStreetMap):
  - [x] Integrasi Leaflet.js Map Picker & Pin Marker.
  - [x] Integrasi Nominatim API untuk Geocoding & Auto-suggest alamat (`osm_place_id`, `latitude`, `longitude`).
- [x] Menyimpan data lokasi (`country`, `province`, `city`, `district`, `address_detail`, `latitude`, `longitude`, `osm_place_id`) di backend.

---

### 🟩 Sprint 3: Monetisasi & Integrasi Midtrans Payment Gateway
- [x] Integrasi Midtrans Snap API / Core API (Sandbox Mode).
- [x] Paket Langganan:
  - [x] **Paket Melamar (Job Applicant Plan - Rp 5.000 / Bulan)**.
  - [x] **Paket Posting Keahlian (Skill Showcase Plan - Rp 10.000 / Bulan)**.
- [x] Endpoint `POST /api/v1/subscriptions/checkout` (Generate Midtrans Payment Token/URL).
- [x] Endpoint `POST /api/v1/subscriptions/webhook` (Menangani notifikasi status pembayaran Midtrans).
- [x] Middleware Backend & Helper: `CheckActiveSubscription(planType)` untuk mengecek status berlangganan aktif.
- [x] UI React: Halaman Checkout Langganan, Tab Paket Langganan di User Profile, & Modal Transaksi Pembayaran.

---

### 🟩 Sprint 4: Model 1 - Talent Showcase (Worker Offer-Driven)
- [x] **Modul Pekerja**:
  - [x] Endpoint `POST /api/v1/skills` (Post Keahlian, Tarif, Kategori Serabutan/Profesional - Proteksi Paket 10K. *Lokasi otomatis inherit dari profil user jika kosong*).
  - [x] Endpoint `GET /api/v1/skills/my-postings` (Kelola postingan keahlian saya).
- [x] **Modul Pemberi Kerja**:
  - [x] Endpoint `GET /api/v1/skills` (Search & Filter Talent berdasarkan Provinsi, Kota, Kecamatan, & Kategori).
  - [x] Endpoint `POST /api/v1/offers` (Kirim Direct Offer pekerjaan ke Pekerja).
- [x] **Modul Respon Worker**:
  - [x] Endpoint `GET /api/v1/offers/received` (Lihat daftar tawaran pekerjaan masuk).
  - [x] Endpoint `PATCH /api/v1/offers/:id/respond` (Setujui atau Tolak Tawaran).
- [x] UI React: Halaman Etalase Talent, Form Post Skill, Modal Kirim Offer, & Panel Tawaran Masuk.

---

### 🟩 Sprint 5: Model 2 - Job Board (Demand Application-Driven)
- [x] **Modul Pemberi Kerja (Gratis / Free)**:
  - [x] Endpoint `POST /api/v1/jobs` (Post Lowongan Kerja Insidental - Gratis. *Lokasi otomatis inherit dari profil user jika kosong*).
  - [x] Endpoint `GET /api/v1/jobs/my-postings` (Daftar lowongan kerja yang saya buka).
- [x] **Modul Pencari Kerja**:
  - [x] Endpoint `GET /api/v1/jobs` (Guest Preview & Filter Lowongan per Provinsi/Kota/Kecamatan).
  - [x] Endpoint `POST /api/v1/jobs/:id/apply` (Melamar Lowongan - Syarat: Login, KTP Verified, Resume Upload, & Active 5K Subscription).
- [x] **Modul Manajemen Pelamar**:
  - [x] Endpoint `GET /api/v1/jobs/:id/applicants` (Pemberi Kerja meninjau daftar pelamar & CV/KTP).
  - [x] Endpoint `PATCH /api/v1/jobs/applications/:id/accept` (Pemberi Kerja menerima pelamar).
- [x] UI React: Halaman Job Board Public, Form Post Job, Button Apply Job (dengan proteksi modal langganan), & Panel Pelamar Lowongan.

---

### 🟩 Sprint 6: Match Connection, Reciprocal Rating & Final Polish
- [x] **Struktur Match "TERHUBUNG"**:
  - [x] Membuat otomatis record `connections` saat Offer disetujui ATAU Application diterima.
  - [x] Endpoint `GET /api/v1/connections/active` (Buka kontak direct No. HP/WhatsApp & Detail Alamat Pelaksanaan).
- [x] **Sistem Rating Dua Arah (Reciprocal Rating)**:
  - [x] Endpoint `POST /api/v1/ratings` (Pekerja memberi rating Pemberi Kerja & Pemberi Kerja memberi rating Pekerja).
  - [x] Menghitung otomatis rata-rata rating bintang user.
- [x] **Final Polish & QA**:
  - [x] End-to-End Testing alur Model 1 dan Model 2.
  - [x] Refinement UI/UX & Responsive Web Layout di React.
  - [x] Dokumentasi API Postman Collection.
