CREATE TABLE IF NOT EXISTS job_offers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    skill_posting_id UUID NOT NULL REFERENCES skill_postings(id) ON DELETE CASCADE,
    employer_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    worker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    offered_budget NUMERIC(12,2) NOT NULL,
    work_date DATE NOT NULL,
    status VARCHAR(30) DEFAULT 'PENDING',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT idx_offer_skill_employer UNIQUE(skill_posting_id, employer_id)
);

CREATE INDEX IF NOT EXISTS idx_job_offers_worker_id ON job_offers(worker_id);
