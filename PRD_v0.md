# Product Requirement Document (PRD) - Layanesia

**Nama Proyek:** Layanesia  
**Tipe Platform:** On-Demand Gig & Temporary Staffing / Replacement Marketplace  
**Versi:** 1.0.0  
**Tanggal:** 30 Agustus 2026  
**Figma Flowchart:** [Figma Board - Layanesia Flowchart](https://www.figma.com/board/z0bEkpshXTalDn0ZtF8AK8/Untitled?node-id=0-1&t=LjyiOp8hOBdww0zA-1)  
**Lokasi File:** `/Users/fiqar/Documents/agustus/x-project/layanesia/PRD_v0.md`



---

## 1. Ringkasan Eksekutif & Tujuan Proyek

**Layanesia** adalah platform marketplace kerja dua arah (*dual-sided marketplace*) yang menghubungkan **Pencari Kerja** dan **Pemberi Kerja** untuk pekerjaan bersifat **sementara (on-demand / harian / shift / pengganti insidental)**, baik dalam skala **Serabutan** (tugas rumah tangga/harian) maupun **Profesional** (pengganti staf perusahaan/kantor).

Platform ini **BUKAN** untuk mencari/menawarkan pekerjaan tetap (*full-time permanent*), melainkan berfokus pada kecepatan menemukan tenaga kerja pengganti (*substitute*) atau bantuan mendesak.

### Contoh Kasus Penggunaan (Use Cases):
1. **Penggantian Tenaga Kerja Operasional Mendadak (Urgent Operational Replacement)**:
   * Sebuah perusahaan atau eksekutif membutuhkan **Supir Operasional / Pengemudi Pengganti** selama 1–2 hari karena pengemudi utama berhalangan hadir/sakit.
2. **Pemenuhan Tenaga Kerja Temporer Perusahaan (Corporate Temporary Staffing)**:
   * Tim HRD / Manajemen membutuhkan staf pengganti sementara seperti **Admin Operasional**, **Office Boy (OB)**, atau **Staff Logistik** untuk mengisi kekosongan posisi sementara selama periode cuti atau lonjakan beban kerja singkat.
3. **Layanan Bantuan Insidental & Tugas Harian (On-Demand Individual & Domestic Assistance)**:
   * Pengguna perorangan membutuhkan tenaga bantuan insidental untuk pemeliharaan properti/halaman, penanganan logistik (pindahan), atau pembersihan area harian.


---

## 2. Fitur Utama & Prinsip Desain Produk

### 2.1 Concept Single Account Dual-Role (Satu Akun Dua Peran)
* **Unified Account**: Setiap pengguna (*User*) memiliki 1 akun yang sama dan dapat bertindak sebagai **Pencari Kerja** (menawarkan keahlian / melamar) maupun **Pemberi Kerja** (membuat lowongan / mencari pekerja).
* **Role Mode Switcher**: Pengguna dapat dengan mudah berpindah konteks antara *"Saya Mencari Kerja"* dan *"Saya Membutuhkan Bantuan / Pekerja"*.

### 2.2 Kategori Pekerjaan
1. **Serabutan**: Pekerjaan ringan/harian tanpa kualifikasi formal rumit (contoh: cuci baju, bersihkan halaman, bantu pindahan, asisten rumah tangga harian).
2. **Profesional**: Pekerjaan yang membutuhkan keahlian/lisensi khusus (contoh: supir pribadi/perusahaan dengan SIM A/B, Admin kantor, Office Boy, Teknisi IT temporer).

---

## 3. Alur Kerja Aplikasi (Workflow Architecture)

Aplikasi Layanesia mengadopsi struktur alur kerja modular yang membedakan **Tahap Inisiasi (Pendaftaran & Verifikasi)** dengan **2 Model Interaksi Pasar (Talent Showcase vs Job Board)** sesuai dengan rancangan [Figma Board Layanesia Flowchart](https://www.figma.com/board/z0bEkpshXTalDn0ZtF8AK8/Untitled?node-id=0-1&t=LjyiOp8hOBdww0zA-1).

### 3.1 Diagram Alur Kerja Utama (Workflow Diagram)


```text
┌───────────────────────────────────────────────────────────┐
│                Dashboard (Public / Guest)                 │
└─────────────────────────────┬─────────────────────────────┘
                              │
                      Login / Register
                              │
┌─────────────────────────────┴─────────────────────────────┐
│                  Profil & Verifikasi KTP                  │
└─────────────────────────────┬─────────────────────────────┘
                              │
                              ▼
                [ Pilih Mode Transaksi ]
                              │
              ┌───────────────┴───────────────┐
              │                               │
              ▼                               ▼
┌──────────────────────────┐     ┌──────────────────────────┐
│ MODEL 1: TALENT SHOWCASE │     │    MODEL 2: JOB BOARD    │
│  (Worker Offer-Driven)   │     │  (Demand Apply-Driven)   │
├──────────────────────────┤     ├──────────────────────────┤
│ 1. Worker post keahlian  │     │ 1. Employer post lowongan│
│ 2. Employer cari talent  │     │ 2. Worker apply lowongan │
│ 3. Employer buat offer   │     │ 3. Employer terima worker│
│ 4. Worker setuju offer   │     │                          │
└─────────────┬────────────┘     └────────────┬─────────────┘
              │                               │
              └───────────────┬───────────────┘
                              │
                              ▼
┌───────────────────────────────────────────────────────────┐
│                     Status: TERHUBUNG                     │
└─────────────────────────────┬─────────────────────────────┘
                              │
                              ▼
┌───────────────────────────────────────────────────────────┐
│       Mutual Rating & Review System (Pekerja ↔ Employer)  │
└───────────────────────────────────────────────────────────┘
```

---

### 3.2 Matriks Perbandingan Alur Kerja (Workflow Comparison)

| Tahapan Flow | Model 1: Talent Showcase (*Worker Offer-Driven*) | Model 2: Job Board (*Demand Application-Driven*) |
| :--- | :--- | :--- |
| **Inisiator Awal** | **Pencari Kerja** (Memposting keahlian & tarif) | **Pemberi Kerja** (Memposting lowongan kerja) |
| **Pemicu Transaksi** | Pemberi Kerja mengirimkan *Direct Offer* | Pencari Kerja mengajukan *Apply Job* |
| **Persyaratan Pencari Kerja** | Profil Lengkap + Verifikasi KTP | Profil Lengkap + Verifikasi KTP + Resume + **Status Subscribed** |
| **Konfirmasi Akhir** | **Pencari Kerja** menyetujui tawaran masuk | **Pemberi Kerja** menerima pelamar kerja |
| **Output Match** | Status **TERHUBUNG** (Kontak direct & lokasi terbuka) | Status **TERHUBUNG** (Kontak direct & lokasi terbuka) |
| **Pasca Pekerjaan** | Reciprocal Rating (Ulasan 2 arah) | Reciprocal Rating (Ulasan 2 arah) |


---

## 4. Spesifikasi Fungsional Detail

### 4.1 Modul Pengguna & Autentikasi (`/server` & `/client`)
* **Registrasi & Login**: Email/Password atau OTP Phone Number.
* **Manajemen Profil**:
  * Informasi Pribadi & Foto Profil.
  * Dokumen Verifikasi KTP (wajib untuk melamar & membuat tawaran demi keamanan).
  * Upload Resume / Portfolio / Lisensi (misal SIM untuk supir).
* **Status Berlangganan (Subscription)**:
  * Fitur pelamaran kerja (*Apply Job*) dibatasi khusus untuk user yang berstatus **Subscribed**.

### 4.2 Modul Model 1: Worker Skill Showcase (Pencari Kerja Offer-Driven)
1. **Posting Keahlian (Worker)**:
   * Menentukan kategori: Serabutan / Profesional.
   * Deskripsi keahlian, pengalaman singkat, & kesiapan waktu (misal: "Siap kerja harian supir SIM A", "Bisa cuci & setrika harian").
   * Tarif acuan (per jam / per hari).
2. **Pencarian Talent (Employer)**:
   * Filter berdasarkan lokasi, kategori (serabutan/profesional), dan ketersediaan.
3. **Direct Offer (Pemberi Kerja -> Worker)**:
   * Pemberi kerja mengirim tawaran pekerjaan dengan rincian tugas, tanggal, dan nilai bayaran.
   * Pekerja menerima notifikasi dan memilih **Setujui / Tolak**.

### 4.3 Modul Model 2: Job Board (Pemberi Kerja Demand-Driven)
1. **Posting Lowongan Pekerjaan (Employer)**:
   * Judul tugas, Kategori (Serabutan/Profesional).
   * Deskripsi kebutuhan (contoh: "Dibutuhkan Pengemudi Operasional Pengganti - Shift 1 Hari" atau "Staf Admin Temporer - Coverage Cuti 3 Hari").
   * Tanggal pelaksanaan, durasi (jam/hari), lokasi, dan honor/gaji.
2. **Pencarian Lowongan & Melamar (Worker)**:
   * Guest / Unauthenticated user bisa melihat preview daftar lowongan di Dashboard.
   * Pekerja melengkapi profil + KTP + Resume.
   * Pekerja mengklik **Apply Pekerjaan** (Memeriksa syarat: Wajib Login & Status Subscription Aktif).
3. **Manajemen Pelamar (Employer)**:
   * Meninjau daftar pelamar (*list pelamar*) beserta profil, KTP terverifikasi, dan resume.
   * Menekan tombol **Terima Pekerja**.

### 4.4 Modul Match & Keterhubungan (Terhubung)
* Ketika tawaran disetujui (Model 1) ATAU pelamar diterima (Model 2), status berubah menjadi **TERHUBUNG**.
* Membuka akses kontak direct (No. HP/WhatsApp/Chat) & detail lokasi kerja.

### 4.5 Modul Rating & Reputasi Dua Arah
* **Beri Rating Pemberi Kerja**: Pekerja memberikan rating 1-5 bintang & ulasan tentang kelayakan pemberi kerja/pembayaran.
* **Beri Rating Pekerja**: Pemberi Kerja memberikan rating 1-5 bintang & ulasan tentang kinerja pekerja.

### 4.6 Modul Monetisasi & Berlangganan (Monetization Model)
* **Pemberi Kerja (Post Job Vacancy)**: **GRATIS (FREE)**
  * Memposting lowongan kerja insidental (Model 2) & mengelola pelamar tanpa dipungut biaya.
* **Pencari Kerja - Paket Melamar (Job Applicant Plan - Rp 5.000 / Bulan)**:
  * Berlangganan bulanan untuk mendapatkan hak melamar pekerjaan (*Apply Job*) pada lowongan yang tersedia di Model 2.
* **Pencari Kerja - Paket Posting Keahlian (Skill Showcase Plan - Rp 10.000 / Bulan)**:
  * Berlangganan bulanan untuk memposting profil keahlian, tarif, & kesiapan kerja (*Post Skill Showcase*) pada Model 1 agar dicari & ditawarkan langsung oleh Pemberi Kerja.
* **Fitur Add-On: Job & Talent Booster (Pay-Per-Blast - Rp 5.000 / Sekali Blast)**:
  * Biaya opsional **Rp 5.000 per sekali blast** untuk mengirim notifikasi pesan masif (Email & WhatsApp) ke seluruh user terdaftar di **Kecamatan / Kota target**.


### 4.7 Modul Filter Lokasi Hiperlokal (Location Hierarchy & OpenStreetMap)
* **Struktur Lokasi Bertingkat**: `Negara` > `Provinsi` > `Kota / Kabupaten` > `Kecamatan`.
* **Pewarisan Lokasi Postingan (Default Location Inheritance)**:
  * Pada UI Form Posting Job maupun Skill, lokasi secara otomatis mewarisi (*default*) lokasi domisili dari profil `users` (`province`, `city`, `district`, `address_detail`).
  * Form menyediakan input opsional (*override*) jika lokasi pekerjaan/keahlian berbeda dari lokasi profil user.
* **Integrasi OpenStreetMap + Leaflet.js**: Menggunakan petaan open-source Leaflet.js & Nominatim API (100% Gratis) untuk penentuan titik lokasi geografis, visualisasi peta interaktif, dan pencarian alamat presisi hingga tingkat **Kecamatan**.

### 4.8 Catatan Pengembangan MVP (MVP Scope Notes)
* **Fitur Broadcast Booster (WA & Email Blast)**: Ditangguhkan untuk pengembangan fase **Post-MVP** (setelah versi MVP rilis), sehingga fokus MVP 100% pada *Core Flow* (Auth, Post Job, Post Skill, Apply Job, Midtrans Payment, Match Connection, & Rating).

---

## 5. Arsitektur Teknologi & Folder Structure

### 5.1 Tech Stack MVP
* **Backend**:
  * Bahasa: **Go (Golang)**
  * Framework: **Fiber v2 / v3**
  * ORM: **GORM**
  * Database: **PostgreSQL**
  * Autentikasi: **JWT** (JSON Web Token) + Bcrypt
* **Layanan Pihak Ketiga (Third-Party MVP Services - 100% Free / Pay-per-Use)**:
  * **File Storage**: **Cloudinary** (Upload & Hosting Foto KTP & Resume PDF).
  * **Payment Gateway**: **Midtrans** (Sandbox/Production untuk Paket Langganan Rp 5.000 & Rp 10.000).
  * **Maps & Geocoding**: **OpenStreetMap + Leaflet.js** (Peta Interaktif Open-Source 100% Gratis) & **Nominatim API**.
* **Frontend**:
  * Library: **React.js** (Vite)
  * Styling: Vanilla CSS / Modular CSS System
  * Peta Frontend: `react-leaflet` / Leaflet.js
  * State Management & Fetcher: Context API / Axios


### 5.2 Layout Direktori Proyek
```text
layanesia/
├── PRD.md                       # Dokumen PRD Ini
├── image.png                    # Flowchart Model 1 (Worker Skill Posting)
├── image copy.png               # Flowchart Model 2 (Employer Job Posting)
├── server/                      # Backend (Go Fiber + GORM + PostgreSQL)
│   ├── cmd/
│   │   └── main.go
│   ├── config/                  # Database & Env Configurations
│   ├── internal/
│   │   ├── controllers/         # HTTP Handlers (Auth, Jobs, Skills, Subscription, Booster, Location)
│   │   ├── models/              # GORM Entities (User, Job, Skill, Subscription, Location, Booster, Rating)
│   │   ├── repositories/        # Database Layer
│   │   ├── services/            # Business Logic Layer (WA Gateway, Email Blast, Location Filter)
│   │   └── middlewares/         # Auth, Tier Check (5K vs 10K), & KTP Check
│   ├── go.mod
│   └── go.sum
└── client/                      # Frontend (React)
    ├── src/
    │   ├── components/
    │   ├── pages/
    │   ├── services/
    │   ├── context/
    │   └── styles/
    ├── package.json
    └── vite.config.js (atau next.config.js)
```

---

## 6. Rancangan Entity Relationship (Database Schema & ERD)

Dokumen rancangan Entity Relationship Diagram (ERD) lengkap beserta skema tabel GORM dan Kamus Data (Data Dictionary) telah dibuat pada file tersendiri: **[ERD_v0.md](file:///Users/fiqar/Documents/agustus/x-project/layanesia/ERD_v0.md)**.


### Ringkasan Tabel & Relasi Utama:
1. **`users`**: Entitas pengguna utama (Identitas, Auth, `active_mode` SEEKER/EMPLOYER, Verifikasi KTP, Resume, Lokasi Domisili).
2. **`subscriptions`**: Transaksi berlangganan (`APPLY_JOB_5K` & `POST_SKILL_10K`).
3. **`locations`**: Referensi master data lokasi hierarki (`Negara` > `Provinsi` > `Kota/Kabupaten` > `Kecamatan`).
4. **`skill_postings`**: Postingan etalase keahlian pekerja (Model 1 Showcase).
5. **`job_postings`**: Postingan lowongan kerja insidental pemberi kerja (Model 2 Job Board - Gratis).
6. **`job_applications`**: Pengajuan lamaran pencari kerja ke lowongan (Constraint Unique: 1 apply/job).
7. **`job_offers`**: Penawaran kerja langsung dari pemberi kerja ke keahlian pekerja (Constraint Unique: 1 pending offer/skill).
8. **`booster_blasts`**: 🔒 **Post-MVP** — Transaksi & log broadcast WhatsApp & Email Rp 5.000 / blast.
9. **`connections`**: Status **TERHUBUNG / MATCHED** setelah offer/application disetujui (Unique constraint on source).
10. **`ratings`**: Ulasan & Rating 1-5 bintang dua arah (Constraint Unique: 1 rating/connection/reviewer).


---

## 7. Metrik Keberhasilan (KPIs)
1. **Match Speed**: Waktu rata-rata dari posting lowongan/keahlian hingga terhubung (*Connected*) < 2 jam untuk kebutuhan insidental (diakselerasi dengan fitur Booster).
2. **KTP Verification Rate**: Persentase user terverifikasi KTP untuk keamanan transaksi kerja sementara.
3. **Hyperlocal Search Accuracy**: Filter lokasi hingga tingkat **Kecamatan** menghasilkan tingkat kecocokan jarak < 10 km.
4. **Reciprocal Rating Completion**: > 80% transaksi yang selesai mendapatkan ulasan dua arah.

