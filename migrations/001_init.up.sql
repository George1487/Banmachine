CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    full_name TEXT NOT NULL,
    role TEXT NOT NULL,
    group_name TEXT NULL, 
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_users_role CHECK (role IN ('teacher', 'student')),
    CONSTRAINT chk_users_group_name CHECK (
        (role = 'teacher' AND group_name is NULL)
        OR
        (role = 'student' AND group_name is NOT NULL)
    )
);

CREATE TABLE IF NOT EXISTS labs (
    lab_id UUID PRIMARY KEY,
    teacher_id UUID NOT NULL REFERENCES users(user_id),
    title TEXT NOT NULL,
    description TEXT,
    status TEXT NOT NULL,
    deadline_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_labs_status CHECK (status IN ('active', 'closed'))
);

CREATE INDEX IF NOT EXISTS idx_labs_teacher_id_status ON labs(teacher_id, status);

CREATE TABLE IF NOT EXISTS submissions (
    submission_id UUID PRIMARY KEY,
    lab_id UUID NOT NULL REFERENCES labs(lab_id),
    student_id UUID NOT NULL REFERENCES users(user_id),
    status TEXT NOT NULL,
    source_file_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    storage_key TEXT NOT NULL,
    submitted_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_submissions_status CHECK (status IN ('uploaded', 'parsing', 'parsed', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_submissions_lab_id_student_id ON submissions(lab_id, student_id);
CREATE INDEX IF NOT EXISTS idx_submissions_student_id ON submissions(student_id);
CREATE INDEX IF NOT EXISTS idx_submissions_lab_student_submitted_at ON submissions(lab_id, student_id, submitted_at DESC);

CREATE TABLE IF NOT EXISTS parsed_submissions (
    parsed_submission_id UUID PRIMARY KEY,
    submission_id UUID NOT NULL UNIQUE REFERENCES submissions(submission_id) ON DELETE CASCADE,
    raw_text TEXT NOT NULL,
    structured_data JSONB NOT NULL,
    parsed_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_parsed_submissions_submission_id ON parsed_submissions(submission_id);

CREATE TABLE IF NOT EXISTS ingest_jobs (
    ingest_job_id UUID PRIMARY KEY,
    submission_id UUID NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,
    error_message TEXT,

    CONSTRAINT chk_ingest_job_status CHECK (status IN ('pending', 'processing', 'done', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_ingest_jobs_status_created_at ON ingest_jobs(status, created_at DESC);
CREATE UNIQUE INDEX IF NOT EXISTS uq_ingest_jobs_active_submission ON ingest_jobs(submission_id) WHERE status IN ('pending', 'processing');

CREATE TABLE IF NOT EXISTS analysis_jobs (
    analysis_job_id UUID PRIMARY KEY,
    lab_id UUID NOT NULL REFERENCES labs(lab_id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    created_by UUID NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMPTZ NOT NULL,
    started_at TIMESTAMPTZ NULL,
    finished_at TIMESTAMPTZ NULL,
    error_message TEXT NULL,

    CONSTRAINT chk_analysis_jobs_status CHECK (status IN ('pending', 'processing', 'done', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_analysis_jobs_lab_id ON analysis_jobs(lab_id);

CREATE TABLE IF NOT EXISTS analysis_job_snapshots(
    analysis_job_id UUID NOT NULL REFERENCES analysis_jobs(analysis_job_id) ON DELETE CASCADE,
    submission_id UUID NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    PRIMARY KEY (analysis_job_id, submission_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_analysis_job_snapshots_submission_id ON analysis_job_snapshots(submission_id);

CREATE TABLE IF NOT EXISTS pairwise_similarities(
    pairwise_similarity_id UUID PRIMARY KEY,
    analysis_job_id UUID NOT NULL REFERENCES analysis_jobs(analysis_job_id) ON DELETE CASCADE,
    lab_id UUID NOT NULL REFERENCES labs(lab_id),
    left_submission_id UUID NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    right_submission_id UUID NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    text_score NUMERIC(5,4) NOT NULL,
    calculation_score NUMERIC(5,4) NOT NULL,
    images_score NUMERIC(5,4) NOT NULL,
    final_score NUMERIC(5,4) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),


    CONSTRAINT chk_pairwise_similarities_text_score CHECK (text_score >= 0 AND text_score <= 1),
    CONSTRAINT chk_pairwise_similarities_calculation_score CHECK (calculation_score >= 0 AND calculation_score <= 1),
    CONSTRAINT chk_pairwise_similarities_images_score CHECK (images_score >= 0 AND images_score <= 1),
    CONSTRAINT chk_pairwise_similarities_final_score CHECK (final_score >= 0 AND final_score <= 1),
    CONSTRAINT chk_pairwise_similarities_left_right_submission CHECK (left_submission_id <> right_submission_id)
);

CREATE INDEX IF NOT EXISTS idx_pairwise_similarities_analysis_job_id ON pairwise_similarities(analysis_job_id);
CREATE INDEX IF NOT EXISTS idx_pairwise_similarities_lab_id ON pairwise_similarities(lab_id);
CREATE INDEX IF NOT EXISTS idx_pairwise_similarities_left_submission_id ON pairwise_similarities(left_submission_id);
CREATE INDEX IF NOT EXISTS idx_pairwise_similarities_right_submission_id ON pairwise_similarities(right_submission_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_pairwise_similarities_job_pair ON pairwise_similarities(analysis_job_id, left_submission_id, right_submission_id);

CREATE TABLE IF NOT EXISTS submission_analysis_summaries(
    submission_analysis_summary_id UUID PRIMARY KEY,
    analysis_job_id UUID NOT NULL REFERENCES analysis_jobs(analysis_job_id) ON DELETE CASCADE,
    submission_id UUID NOT NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    top_match_submission_id UUID NULL REFERENCES submissions(submission_id) ON DELETE CASCADE,
    top_match_score NUMERIC(5,4) NULL,
    final_score_risk_level TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_submission_analysis_summaries_top_match_score CHECK (top_match_score >= 0 AND top_match_score <= 1),
    CONSTRAINT chk_submission_analysis_summaries_final_score_risk_level CHECK (final_score_risk_level IN ('low', 'medium', 'high'))
);