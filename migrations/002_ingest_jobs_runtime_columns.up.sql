ALTER TABLE ingest_jobs
    ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ;

ALTER TABLE ingest_jobs
    ADD COLUMN IF NOT EXISTS finished_at TIMESTAMPTZ;

ALTER TABLE ingest_jobs
    ADD COLUMN IF NOT EXISTS error_message TEXT;

CREATE INDEX IF NOT EXISTS idx_ingest_jobs_status_created_at
    ON ingest_jobs(status, created_at DESC);

CREATE UNIQUE INDEX IF NOT EXISTS uq_ingest_jobs_active_submission
    ON ingest_jobs(submission_id) WHERE status IN ('pending', 'processing');
