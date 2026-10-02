# Entity Relationship Diagram (ERD) - Layanesia (Box Visual Format)

**Proyek:** Layanesia  
**Database:** PostgreSQL (GORM ORM)  
**Versi Schema:** 1.0.1 (Revised with Active Mode, Province Hierarchy, OSM Labels & Unique Constraints)  
**File Location:** `/Users/fiqar/Documents/agustus/x-project/layanesia/ERD_v0.md`


---

## 1. Diagram Skema Relasi (Visual Box ERD)

```text
┌───────────────────────────────────────────────────────────┐
│                          USERS                            │
├───────────────────────────────────────────────────────────┤
│ PK  id               UUID                                 │
│ UK  email            VARCHAR(255)                         │
│ UK  phone            VARCHAR(50)                          │
│     name             VARCHAR(255)                         │
│     active_mode      VARCHAR(20) DEFAULT 'SEEKER'        │
│     is_ktp_verified  BOOLEAN                              │
│     country, province, city, district                     │
└─────────────────────────────┬─────────────────────────────┘
                              │
          ┌───────────────────┼───────────────────┬───────────────────┐
          │ (1:N)             │ (1:N)             │ (1:N)             │ (1:N)
          ▼                   ▼                   ▼                   ▼
┌───────────────────┐ ┌───────────────────┐ ┌───────────────────┐ ┌───────────────────┐
│   SUBSCRIPTIONS   │ │  SKILL_POSTINGS   │ │   JOB_POSTINGS    │ │  BOOSTER_BLASTS   │
├───────────────────┤ ├───────────────────┤ ├───────────────────┤ │  [🔒 Post-MVP]    │
│ PK id      UUID   │ │ PK id      UUID   │ │ PK id      UUID   │ ├───────────────────┤
│ FK user_id UUID   │ │ FK user_id UUID   │ │ FK empl_id UUID   │ │ PK id      UUID   │
│    plan_type      │ │    category       │ │    category       │ │ FK user_id UUID   │
│    amount  5k/10k │ │    rate_amount    │ │    budget         │ │    amount  5000   │
└───────────────────┘ └─────────┬─────────┘ └─────────┬─────────┘ └───────────────────┘
                                │                     │
                                │ (1:N)               │ (1:N)
                                ▼                     ▼
                      ┌───────────────────┐ ┌───────────────────┐
                      │    JOB_OFFERS     │ │ JOB_APPLICATIONS  │
                      ├───────────────────┤ ├───────────────────┤
                      │ PK id      UUID   │ │ PK id      UUID   │
                      │ FK skill_id UUID  │ │ FK job_id  UUID   │
                      │ FK empl_id  UUID  │ │ FK app_id  UUID   │
                      │ UK(skill,empl,PE) │ │ UK(job_id,app_id) │
                      └─────────┬─────────┘ └─────────┬─────────┘
                                │ (On Accept)         │ (On Accept)
                                └───────────┬─────────┘
                                            ▼
                                  ┌───────────────────┐
                                  │    CONNECTIONS    │
                                  ├───────────────────┤
                                  │ PK id      UUID   │
                                  │    source_type    │
                                  │ FK empl_id UUID   │
                                  │ FK worker_id UUID │
                                  │ UK(empl,work,src) │
                                  └─────────┬─────────┘
                                            │ (1:N)
                                            ▼
                                  ┌───────────────────┐
                                  │      RATINGS      │
                                  ├───────────────────┤
                                  │ PK id      UUID   │
                                  │ FK conn_id UUID   │
                                  │ FK reviewer_id    │
                                  │ UK(conn_id,rev_id)│
                                  └───────────────────┘
```

---

## 2. Struktur Tabel Detil (Box Format Kamus Data)

### 2.1 Tabel `users` (Pengguna Platform)
```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ Unique User ID           │
│ email            │ VARCHAR(255) │ UNIQUE, NOT │ Email Login              │
│ phone            │ VARCHAR(50)  │ UNIQUE, NOT │ No HP (OTP/WhatsApp)     │
│ name             │ VARCHAR(255) │ NOT NULL    │ Nama Lengkap             │
│ password         │ VARCHAR(255) │ NOT NULL    │ Pass Hash (Bcrypt)       │
│ active_mode      │ VARCHAR(20)  │ DEFAULT SEEK│ Mode: SEEKER / EMPLOYER  │
│ is_ktp_verified  │ BOOLEAN      │ DEFAULT F   │ Status Verifikasi KTP    │
│ ktp_number       │ VARCHAR(50)  │ NULLABLE    │ Nomor KTP                │
│ ktp_image_url    │ VARCHAR(500) │ NULLABLE    │ Foto KTP Storage URL     │
│ resume_url       │ VARCHAR(500) │ NULLABLE    │ Resume / CV File URL     │
│ country          │ VARCHAR(100) │ DEFAULT ID  │ Negara Domisili          │
│ province         │ VARCHAR(100) │ NOT NULL    │ Provinsi                 │
│ city             │ VARCHAR(100) │ NOT NULL    │ Kota / Kabupaten         │
│ district         │ VARCHAR(100) │ NOT NULL    │ Kecamatan                │
│ address_detail   │ TEXT         │ NULLABLE    │ Alamat Jalan / Patokan   │
│ latitude         │ DECIMAL(10,8)│ NULLABLE    │ [Opsional] Nominatim/OSM │
│ longitude        │ DECIMAL(11,8)│ NULLABLE    │ [Opsional] Nominatim/OSM │
│ formatted_address│ TEXT         │ NULLABLE    │ [Opsional] Nominatim/OSM │
│ osm_place_id     │ VARCHAR(255) │ NULLABLE    │ [Opsional] Nominatim/OSM │
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Buat Akun          │
│ updated_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Update Akun        │
└──────────────────┴──────────────┴─────────────┴──────────────────────────┘
```

### 2.2 Tabel `subscriptions` (Transaksi Berlangganan)
```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Transaksi Langganan   │
│ user_id          │ UUID         │ FK(users)   │ Relasi ke User           │
│ plan_type        │ VARCHAR(50)  │ NOT NULL    │ APPLY_JOB_5K /           │
│                  │              │             │ POST_SKILL_10K           │
│ amount           │ NUMERIC(12,2)│ NOT NULL    │ Rp 5.000 / Rp 10.000     │
│ status           │ VARCHAR(30)  │ NOT NULL    │ ACTIVE / EXPIRED         │
│ starts_at        │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Mulai              │
│ expires_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Berakhir           │
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Transaksi          │
└──────────────────┴──────────────┴─────────────┴──────────────────────────┘
```

### 2.3 Tabel `skill_postings` (Model 1: Workers Showcase)
> **Catatan Lokasi:** Jika input lokasi saat posting dikosongkan pada UI, sistem otomatis mengambil default dari lokasi profil `users` (`province`, `city`, `district`, `address_detail`).

```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Post Keahlian         │
│ user_id          │ UUID         │ FK(users)   │ Pekerja (Skill Owner)    │
│ category         │ VARCHAR(50)  │ NOT NULL    │ SERABUTAN / PROFESIONAL  │
│ title            │ VARCHAR(255) │ NOT NULL    │ Judul Keahlian           │
│ description      │ TEXT         │ NOT NULL    │ Deskripsi Pengalaman     │
│ country          │ VARCHAR(100) │ DEFAULT ID  │ Negara                   │
│ province         │ VARCHAR(100) │ NOT NULL    │ Provinsi (Inherit/Custom)│
│ city             │ VARCHAR(100) │ NOT NULL    │ Kota (Inherit/Custom)    │
│ district         │ VARCHAR(100) │ NOT NULL    │ Kecamatan (Inherit/Custom│
│ address_detail   │ TEXT         │ NULLABLE    │ Alamat Jalan / Patokan   │
│ latitude         │ DECIMAL(10,8)│ NULLABLE    │ [Opsional] Nominatim/OSM │
│ longitude        │ DECIMAL(11,8)│ NULLABLE    │ [Opsional] Nominatim/OSM │
│ formatted_address│ TEXT         │ NULLABLE    │ [Opsional] Nominatim/OSM │
│ osm_place_id     │ VARCHAR(255) │ NULLABLE    │ [Opsional] Nominatim/OSM │
│ rate_type        │ VARCHAR(30)  │ NOT NULL    │ PER_HOUR / PER_DAY       │
│ rate_amount      │ NUMERIC(12,2)│ NOT NULL    │ Tarif Acuan (Rp)         │
│ availability     │ VARCHAR(30)  │ DEFAULT AV  │ AVAILABLE / BUSY         │
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Post               │
│ updated_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Update             │
└──────────────────┴──────────────┴─────────────┴──────────────────────────┘
```

### 2.4 Tabel `job_postings` (Model 2: Demand Job Board - Gratis)
> **Catatan Lokasi:** Jika input lokasi saat posting dikosongkan pada UI, sistem otomatis mengambil default dari lokasi profil `users` (`province`, `city`, `district`, `address_detail`).

```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Lowongan Kerja        │
│ employer_id      │ UUID         │ FK(users)   │ Pemberi Kerja            │
│ category         │ VARCHAR(50)  │ NOT NULL    │ SERABUTAN / PROFESIONAL  │
│ title            │ VARCHAR(255) │ NOT NULL    │ Judul Lowongan           │
│ description      │ TEXT         │ NOT NULL    │ Kebutuhan Tugas          │
│ country          │ VARCHAR(100) │ DEFAULT ID  │ Negara                   │
│ province         │ VARCHAR(100) │ NOT NULL    │ Provinsi (Inherit/Custom)│
│ city             │ VARCHAR(100) │ NOT NULL    │ Kota (Inherit/Custom)    │
│ district         │ VARCHAR(100) │ NOT NULL    │ Kecamatan (Inherit/Custom│
│ address_detail   │ TEXT         │ NULLABLE    │ Alamat Jalan / Patokan   │
│ latitude         │ DECIMAL(10,8)│ NULLABLE    │ [Opsional] Nominatim/OSM │
│ longitude        │ DECIMAL(11,8)│ NULLABLE    │ [Opsional] Nominatim/OSM │
│ formatted_address│ TEXT         │ NULLABLE    │ [Opsional] Nominatim/OSM │
│ osm_place_id     │ VARCHAR(255) │ NULLABLE    │ [Opsional] Nominatim/OSM │
│ work_date        │ DATE         │ NOT NULL    │ Tanggal Pelaksanaan      │
│ duration_type    │ VARCHAR(30)  │ NOT NULL    │ HOURS / DAYS             │
│ duration_value   │ INT          │ NOT NULL    │ Nilai Durasi             │
│ budget           │ NUMERIC(12,2)│ NOT NULL    │ Total Honor / Gaji       │
│ status           │ VARCHAR(30)  │ DEFAULT OP  │ OPEN/IN_PROGRESS/DONE    │
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Post Lowongan      │
│ updated_at       │ TIMESTAMPTZ  │ NOT NULL    │ Waktu Update Lowongan    │
└──────────────────┴──────────────┴─────────────┴──────────────────────────┘
```

### 2.5 Tabel `job_applications` (Lamaran Pencari Kerja)
```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Lamaran               │
│ job_posting_id   │ UUID         │ FK(jobs)    │ Relasi Lowongan Target   │
│ applicant_id     │ UUID         │ FK(users)   │ Pekerja Pelamar          │
│ status           │ VARCHAR(30)  │ DEFAULT PE  │ PENDING/ACCEPTED/REJECTED│
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Tanggal Apply            │
│ updated_at       │ TIMESTAMPTZ  │ NOT NULL    │ Tanggal Status Changed   │
├──────────────────┴──────────────┴─────────────┴──────────────────────────┤
│ TABLE CONSTRAINT: UNIQUE(job_posting_id, applicant_id)                  │
│ (Mencegah 1 pelamar melamar berkali-kali ke lowongan yang sama)          │
└──────────────────────────────────────────────────────────────────────────┘
```

### 2.6 Tabel `job_offers` (Tawaran Direct Employer ke Worker)
```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Direct Offer          │
│ skill_posting_id │ UUID         │ FK(skills)  │ Relasi Post Skill Target │
│ employer_id      │ UUID         │ FK(users)   │ Pemberi Kerja            │
│ worker_id        │ UUID         │ FK(users)   │ Pekerja Target           │
│ offered_budget   │ NUMERIC(12,2)│ NOT NULL    │ Nominal Tawaran Honor    │
│ work_date        │ DATE         │ NOT NULL    │ Tanggal Kerja            │
│ status           │ VARCHAR(30)  │ DEFAULT PE  │ PENDING/ACCEPTED/REJECTED│
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Tanggal Offer            │
│ updated_at       │ TIMESTAMPTZ  │ NOT NULL    │ Tanggal Status Changed   │
├──────────────────┴──────────────┴─────────────┴──────────────────────────┤
│ TABLE CONSTRAINT: UNIQUE INDEX(skill_posting_id, employer_id)            │
│ WHERE status = 'PENDING'                                                 │
│ (Mencegah offer ganda berkategori PENDING dari employer yang sama)       │
└──────────────────────────────────────────────────────────────────────────┘
```

### 2.7 Tabel `booster_blasts` 🔒 [Post-MVP — Jangan di-include di AutoMigrate Sprint 0]
```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Booster Blast         │
│ user_id          │ UUID         │ FK(users)   │ Pembeli Booster          │
│ target_type      │ VARCHAR(30)  │ NOT NULL    │ JOB / SKILL              │
│ target_id        │ UUID         │ NOT NULL    │ Ref ID Job/Skill Target  │
│ amount           │ NUMERIC(12,2)│ DEFAULT 5000│ Biaya Rp 5.000           │
│ payment_status   │ VARCHAR(30)  │ DEFAULT UNP │ UNPAID / PAID            │
│ target_province  │ VARCHAR(100) │ NOT NULL    │ Provinsi Target Blast    │
│ target_city      │ VARCHAR(100) │ NOT NULL    │ Kota Target Blast        │
│ target_district  │ VARCHAR(100) │ NOT NULL    │ Kecamatan Target Blast   │
│ channels         │ VARCHAR(30)  │ NOT NULL    │ EMAIL / WHATSAPP / BOTH  │
│ total_sent       │ INT          │ DEFAULT 0   │ Total User Ter-blast     │
│ status           │ VARCHAR(30)  │ DEFAULT QUE │ QUEUED/PROCESSING/DONE   │
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Tanggal Transaksi Blast  │
└──────────────────┴──────────────┴─────────────┴──────────────────────────┘
```

### 2.8 Tabel `connections` (Status TERHUBUNG)
> **Known Architectural Risk Note:** Kolom `source_type` & `source_id` menggunakan *Polymorphic Association* ke `job_applications` atau `job_offers`. Karena PostgreSQL/GORM tidak dapat menerapkan native Foreign Key multi-tabel, integritas dijaga oleh Application Layer dan diperkuat oleh Unique Constraint di bawah ini.

```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Match Connection      │
│ source_type      │ VARCHAR(30)  │ NOT NULL    │ JOB_APPLICATION /        │
│                  │              │             │ JOB_OFFER                │
│ source_id        │ UUID         │ NOT NULL    │ Ref ID Offer/Application │
│ employer_id      │ UUID         │ FK(users)   │ ID Pemberi Kerja         │
│ worker_id        │ UUID         │ FK(users)   │ ID Pekerja               │
│ status           │ VARCHAR(30)  │ DEFAULT ACT │ ACTIVE / COMPLETED       │
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Tanggal Terhubung        │
├──────────────────┴──────────────┴─────────────┴──────────────────────────┤
│ TABLE CONSTRAINT: UNIQUE(employer_id, worker_id, source_id)             │
│ (Mencegah terbuatnya record koneksi duplikat dari sumber transaksi sama) │
└──────────────────────────────────────────────────────────────────────────┘
```

### 2.9 Tabel `ratings` (Rating Dua Arah)
```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ UUID         │ PRIMARY KEY │ ID Ulasan Rating         │
│ connection_id    │ UUID         │ FK(conn)    │ Relasi Connection Match  │
│ reviewer_id      │ UUID         │ FK(users)   │ User Pemberi Ulasan      │
│ reviewee_id      │ UUID         │ FK(users)   │ User Penerima Ulasan     │
│ reviewer_role    │ VARCHAR(30)  │ NOT NULL    │ EMPLOYER / WORKER        │
│ rating_stars     │ INT          │ NOT NULL    │ Bintang (1 - 5)          │
│ comment          │ TEXT         │ NULLABLE    │ Catatan Ulasan           │
│ created_at       │ TIMESTAMPTZ  │ NOT NULL    │ Tanggal Ulasan           │
├──────────────────┴──────────────┴─────────────┴──────────────────────────┤
│ TABLE CONSTRAINT: UNIQUE(connection_id, reviewer_id)                    │
│ (Mencegah 1 pihak memberikan rating lebih dari sekali per koneksi)       │
└──────────────────────────────────────────────────────────────────────────┘
```

### 2.10 Tabel `locations` (Master Data Referensi Wilayah)
> **Catatan Arsitektur Lokasi:** Tabel `locations` berfungsi sebagai **Master Data Referensi Autocomplete/Dropdown** pada API backend dan frontend. Kolom lokasi pada `users`, `skill_postings`, dan `job_postings` disimpan secara *denormalized string* (`country`, `province`, `city`, `district`) demi fleksibilitas pencarian berbasis OSM Nominatim tanpa melanggar DB FK constraint.

```text
┌──────────────────┬──────────────┬─────────────┬──────────────────────────┐
│ COLUMN NAME      │ DATA TYPE    │ CONSTRAINT  │ DESCRIPTION              │
├──────────────────┼──────────────┼─────────────┼──────────────────────────┤
│ id               │ SERIAL       │ PRIMARY KEY │ ID Master Lokasi         │
│ country          │ VARCHAR(100) │ DEFAULT ID  │ Nama Negara              │
│ province         │ VARCHAR(100) │ NOT NULL    │ Nama Provinsi            │
│ city_or_regency  │ VARCHAR(100) │ NOT NULL    │ Kota / Kabupaten         │
│ district         │ VARCHAR(100) │ NOT NULL    │ Kecamatan                │
│ postal_code      │ VARCHAR(20)  │ NULLABLE    │ Kode Pos                 │
└──────────────────┴──────────────┴─────────────┴──────────────────────────┘
```

