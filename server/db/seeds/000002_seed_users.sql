-- Seed Sample Users (DEV/STAGING ONLY). Password semua akun: Password123!
-- Kolom password berisi hash bcrypt (cost 10), bukan plaintext.
INSERT INTO users (id, email, phone, name, password, active_mode,
                   is_ktp_verified, ktp_status, ktp_number, ktp_submitted_at, ktp_verified_at,
                   province, city, district, address_detail) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'budi.worker@gmail.com',  '081234567890', 'Budi Santoso',                        '$2a$10$kl/4LvyaERHOAd3417Cq3.ZJBHS1jdaTqXOsy5n1PI2kDE/7DHpBW', 'SEEKER',   true, 'VERIFIED', '3174012304900001', NOW(), NOW(), 'DKI Jakarta', 'Jakarta Selatan', 'Kebayoran Baru', 'Jl. Wijaya I No. 42'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'siti.art@gmail.com',     '081987654321', 'Siti Rahmawati',                      '$2a$10$kl/4LvyaERHOAd3417Cq3.ZJBHS1jdaTqXOsy5n1PI2kDE/7DHpBW', 'SEEKER',   true, 'VERIFIED', '3174025508920003', NOW(), NOW(), 'DKI Jakarta', 'Jakarta Selatan', 'Cilandak',       'Jl. Fatmawati Raya No. 18'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33', 'hrd@logisticsjaya.co.id','081122334455', 'PT Logistics Jaya Mandiri (Hendra)',  '$2a$10$kl/4LvyaERHOAd3417Cq3.ZJBHS1jdaTqXOsy5n1PI2kDE/7DHpBW', 'EMPLOYER', true, 'VERIFIED', '3174091211850005', NOW(), NOW(), 'DKI Jakarta', 'Jakarta Selatan', 'Kebayoran Baru', 'Gedung Graha Logistics Lt. 4'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a44', 'anita.wibowo@gmail.com', '081555666777', 'Ibu Anita Wibowo',                    '$2a$10$kl/4LvyaERHOAd3417Cq3.ZJBHS1jdaTqXOsy5n1PI2kDE/7DHpBW', 'EMPLOYER', true, 'VERIFIED', NULL,               NOW(), NOW(), 'DKI Jakarta', 'Jakarta Selatan', 'Pondok Indah',   'Jl. Metro Pondok Indah Blok TA No. 12')
ON CONFLICT (id) DO NOTHING;
