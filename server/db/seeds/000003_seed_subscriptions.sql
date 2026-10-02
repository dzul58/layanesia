-- Seed Subscriptions (DEV/STAGING ONLY) - sudah dibayar & aktif 30 hari.
INSERT INTO subscriptions (id, user_id, plan_type, amount, status, payment_status, payment_provider, order_id, paid_at, starts_at, expires_at) VALUES
('1bafdaee-cfea-4894-82d1-8c163a1f6acd', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'POST_SKILL_10K', 10000.00, 'ACTIVE', 'PAID', 'SEED', 'LEGACY-SEED-1bafdaee', NOW(), NOW(), NOW() + INTERVAL '30 days'),
('097c880d-1d28-424e-aeab-aba9ff9d4550', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 'APPLY_JOB_5K',    5000.00, 'ACTIVE', 'PAID', 'SEED', 'LEGACY-SEED-097c880d', NOW(), NOW(), NOW() + INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;
