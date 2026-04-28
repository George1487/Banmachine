

from __future__ import annotations

import uuid

import psycopg2.extras

from db.connection import get_conn, put_conn
from models.types import AnalysisJob, PairwiseResult, SnapshotSubmission, SubmissionSummary


def claim_pending_job() -> AnalysisJob | None:
  
    conn = get_conn()
    try:
        with conn:
            with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
                cur.execute(
                    """
                    SELECT analysis_job_id, lab_id, status, created_by
                    FROM analysis_jobs
                    WHERE status = 'pending'
                    ORDER BY created_at
                    LIMIT 1
                    FOR UPDATE SKIP LOCKED
                    """,
                )
                row = cur.fetchone()
                if row is None:
                    return None

                cur.execute(
                    """
                    UPDATE analysis_jobs
                    SET status = 'processing', started_at = NOW()
                    WHERE analysis_job_id = %s
                    """,
                    (row["analysis_job_id"],),
                )

        return AnalysisJob(
            analysis_job_id=row["analysis_job_id"],
            lab_id=row["lab_id"],
            status="processing",
            created_by=row["created_by"],
        )
    finally:
        put_conn(conn)


def load_snapshot_submissions(lab_id: uuid.UUID) -> list[SnapshotSubmission]:
    
    conn = get_conn()
    try:
        with conn:
            with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as cur:
                cur.execute(
                    """
                    SELECT DISTINCT ON (s.student_id)
                        s.submission_id,
                        s.student_id,
                        ps.raw_text,
                        ps.structured_data
                    FROM submissions s
                    JOIN parsed_submissions ps ON ps.submission_id = s.submission_id
                    WHERE s.lab_id = %s
                      AND s.status = 'parsed'
                    ORDER BY s.student_id, s.submitted_at DESC
                    """,
                    (lab_id,),
                )
                rows = cur.fetchall()

        return [
            SnapshotSubmission(
                submission_id=row["submission_id"],
                student_id=row["student_id"],
                raw_text=row["raw_text"],
                structured_data=row["structured_data"],
            )
            for row in rows
        ]
    finally:
        put_conn(conn)


def save_snapshot(job_id: uuid.UUID, submission_ids: list[uuid.UUID]) -> None:
    """Record the exact set of submissions participating in this job."""
    if not submission_ids:
        return

    conn = get_conn()
    try:
        with conn:
            with conn.cursor() as cur:
                psycopg2.extras.execute_values(
                    cur,
                    """
                    INSERT INTO analysis_job_snapshots (analysis_job_id, submission_id)
                    VALUES %s
                    ON CONFLICT DO NOTHING
                    """,
                    [(job_id, sid) for sid in submission_ids],
                )
    finally:
        put_conn(conn)


def save_pairwise(results: list[PairwiseResult]) -> None:
    """Batch-insert pairwise similarity records."""
    if not results:
        return

    conn = get_conn()
    try:
        with conn:
            with conn.cursor() as cur:
                psycopg2.extras.execute_values(
                    cur,
                    """
                    INSERT INTO pairwise_similarities (
                        pairwise_similarity_id,
                        analysis_job_id,
                        lab_id,
                        left_submission_id,
                        right_submission_id,
                        text_score,
                        calculation_score,
                        images_score,
                        final_score
                    ) VALUES %s
                    ON CONFLICT DO NOTHING
                    """,
                    [
                        (
                            uuid.uuid4(),
                            r.analysis_job_id,
                            r.lab_id,
                            r.left_submission_id,
                            r.right_submission_id,
                            round(r.text_score, 4),
                            round(r.calculation_score, 4),
                            round(r.images_score, 4),
                            round(r.final_score, 4),
                        )
                        for r in results
                    ],
                )
    finally:
        put_conn(conn)


def save_summaries(summaries: list[SubmissionSummary]) -> None:
    """Batch-insert submission analysis summary records."""
    if not summaries:
        return

    conn = get_conn()
    try:
        with conn:
            with conn.cursor() as cur:
                psycopg2.extras.execute_values(
                    cur,
                    """
                    INSERT INTO submission_analysis_summaries (
                        submission_analysis_summary_id,
                        analysis_job_id,
                        submission_id,
                        top_match_submission_id,
                        top_match_score,
                        final_score_risk_level
                    ) VALUES %s
                    ON CONFLICT DO NOTHING
                    """,
                    [
                        (
                            uuid.uuid4(),
                            s.analysis_job_id,
                            s.submission_id,
                            s.top_match_submission_id,
                            round(s.top_match_score, 4) if s.top_match_score is not None else None,
                            s.final_score_risk_level,
                        )
                        for s in summaries
                    ],
                )
    finally:
        put_conn(conn)


def complete_job(
    job_id: uuid.UUID,
    status: str,
    error_message: str | None = None,
) -> None:
    """Mark the job as done or failed."""
    conn = get_conn()
    try:
        with conn:
            with conn.cursor() as cur:
                cur.execute(
                    """
                    UPDATE analysis_jobs
                    SET status = %s,
                        finished_at = NOW(),
                        error_message = %s
                    WHERE analysis_job_id = %s
                    """,
                    (status, error_message, job_id),
                )
    finally:
        put_conn(conn)
