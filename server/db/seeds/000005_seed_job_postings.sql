-- Seed Job Postings (Model 2 Demand Job Board - Free)
INSERT INTO job_postings (id, employer_id, category, title, description, province, city, district, work_date, duration_type, duration_value, budget, status) VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'PROFESIONAL', 'Dibutuhkan Supir Operasional Pengganti Shift 2 Hari', 'Dibutuhkan pengemudi boks operasional kantor area Jabodetabek selama 2 hari pengganti staf sakit.', 'DKI Jakarta', 'Jakarta Selatan', 'Kebayoran Baru', CURRENT_DATE + 2, 'DAYS', 2, 600000.00, 'OPEN'),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380c22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'SERABUTAN', 'Bantuan Insidental Bersihkan Kebun & Halaman Rumah', 'Dibutuhkan tenaga bantuan 1 hari (sekitar 4 jam) untuk potong rumput & pembersihan halaman belakang.', 'DKI Jakarta', 'Jakarta Selatan', 'Pondok Indah', CURRENT_DATE + 1, 'HOURS', 4, 150000.00, 'OPEN')
ON CONFLICT (id) DO NOTHING;
