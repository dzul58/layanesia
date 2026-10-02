-- Seed Skill Postings (Model 1 Showcase)
INSERT INTO skill_postings (id, user_id, category, title, description, province, city, district, rate_type, rate_amount, availability) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b11', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'PROFESIONAL', 'Supir Operasional & Pengemudi Pribadi Pengganti', 'Pengalaman 6 tahun supir eksekutif & operasional boks/kantor. SIM A & B1 aktif. Siap kerja harian.', 'DKI Jakarta', 'Jakarta Selatan', 'Kebayoran Baru', 'PER_DAY', 250000.00, 'AVAILABLE'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380b22', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'SERABUTAN', 'Bantuan Asisten Rumah Tangga & Cuci Setrika Harian', 'Bisa membantu cuci setrika, bersihkan rumah harian, & persiapan konsumsi acara keluarga. Teliti.', 'DKI Jakarta', 'Jakarta Selatan', 'Cilandak', 'PER_HOUR', 35000.00, 'AVAILABLE')
ON CONFLICT (id) DO NOTHING;
