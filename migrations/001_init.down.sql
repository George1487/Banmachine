DROP TABLE IF EXISTS submission_analysis_summaries;

DROP INDEX IF EXISTS uq_pairwise_similarities_job_pair;
DROP INDEX IF EXISTS idx_pairwise_similarities_right_submission_id;
DROP INDEX IF EXISTS idx_pairwise_similarities_left_submission_id;
DROP INDEX IF EXISTS idx_pairwise_similarities_lab_id;
DROP INDEX IF EXISTS idx_pairwise_similarities_analysis_job_id;
DROP TABLE IF EXISTS pairwise_similarities;

DROP INDEX IF EXISTS idx_analysis_job_snapshots_submission_id;
DROP TABLE IF EXISTS analysis_job_snapshots;

DROP INDEX IF EXISTS idx_analysis_jobs_lab_id;
DROP TABLE IF EXISTS analysis_jobs;

DROP INDEX IF EXISTS idx_ingest_jobs_status_created_at;
DROP INDEX IF EXISTS uq_ingest_jobs_active_submission;
DROP TABLE IF EXISTS ingest_jobs;

DROP INDEX IF EXISTS idx_parsed_submissions_submission_id;
DROP TABLE IF EXISTS parsed_submissions;

DROP INDEX IF EXISTS idx_submissions_lab_student_submitted_at;
DROP INDEX IF EXISTS idx_submissions_student_id;
DROP INDEX IF EXISTS idx_submissions_lab_id_student_id;
DROP TABLE IF EXISTS submissions;

DROP INDEX IF EXISTS idx_labs_teacher_id_status;
DROP TABLE IF EXISTS labs;

DROP TABLE IF EXISTS users;