"""
Main analysis pipeline orchestrator.
"""

from __future__ import annotations

import logging
import uuid

import config
from analysis.calculation_scorer import compute_calculation_score, compute_numbers
from analysis.image_scorer import compute_image_embedding, compute_image_score
from analysis.text_scorer import compute_embedding, compute_text_score
from db.queries import (
    complete_job,
    load_snapshot_submissions,
    save_pairwise,
    save_snapshot,
    save_summaries,
)
from models.types import AnalysisJob, PairwiseResult, SubmissionSummary

logger = logging.getLogger(__name__)


def _risk_level(score: float) -> str:
    if score >= config.HIGH_THRESHOLD:
        return "high"
    if score >= config.MEDIUM_THRESHOLD:
        return "medium"
    return "low"


def run_analysis(job: AnalysisJob) -> None:
    logger.info("Starting analysis for job %s (lab %s)", job.analysis_job_id, job.lab_id)

   
    submissions = load_snapshot_submissions(job.lab_id)
    logger.info("Snapshot: %d submissions", len(submissions))

    submission_ids = [s.submission_id for s in submissions]
    save_snapshot(job.analysis_job_id, submission_ids)

   
    if len(submissions) < 2:
        logger.info("Less than 2 submissions — marking job done without comparison.")
        complete_job(job.analysis_job_id, "done")
        return

   
    logger.info("Computing embeddings and number sets...")
    text_embeddings = {}
    image_embeddings = {}
    number_sets = {}

    for sub in submissions:
        sid = sub.submission_id
        text_embeddings[sid] = compute_embedding(sub.structured_data)
        image_embeddings[sid] = compute_image_embedding(sub.structured_data)
        number_sets[sid] = compute_numbers(sub.structured_data)
        logger.debug(
            "sub %s: numbers=%d, text_emb=%s, img_emb=%s",
            sid,
            len(number_sets[sid]),
            "ok" if text_embeddings[sid] is not None else "none",
            "ok" if image_embeddings[sid] is not None else "none",
        )

   
    pairs: list[PairwiseResult] = []

    for i, sub_a in enumerate(submissions):
        for sub_b in submissions[i + 1:]:
            sid_a = sub_a.submission_id
            sid_b = sub_b.submission_id

           
            left_id = min(sid_a, sid_b)
            right_id = max(sid_a, sid_b)

            text_score = compute_text_score(
                text_embeddings[sid_a], text_embeddings[sid_b]
            )
            calc_score = compute_calculation_score(
                number_sets[sid_a], number_sets[sid_b]
            )
            img_score = compute_image_score(
                image_embeddings[sid_a], image_embeddings[sid_b]
            )

            final_score = (
                config.TEXT_WEIGHT * text_score
                + config.CALC_WEIGHT * calc_score
                + config.IMG_WEIGHT * img_score
            )
            final_score = max(0.0, min(1.0, final_score))

            logger.debug(
                "pair (%s, %s): text=%.4f calc=%.4f img=%.4f final=%.4f",
                left_id, right_id, text_score, calc_score, img_score, final_score,
            )

            pairs.append(
                PairwiseResult(
                    analysis_job_id=job.analysis_job_id,
                    lab_id=job.lab_id,
                    left_submission_id=left_id,
                    right_submission_id=right_id,
                    text_score=text_score,
                    calculation_score=calc_score,
                    images_score=img_score,
                    final_score=final_score,
                )
            )

    logger.info("Saving %d pairwise results...", len(pairs))
    save_pairwise(pairs)

   
    summaries: list[SubmissionSummary] = []

    for sub in submissions:
        sid = sub.submission_id

      
        related = [
            p for p in pairs
            if p.left_submission_id == sid or p.right_submission_id == sid
        ]

        if not related:
           
            summaries.append(
                SubmissionSummary(
                    analysis_job_id=job.analysis_job_id,
                    submission_id=sid,
                    top_match_submission_id=None,
                    top_match_score=None,
                    final_score_risk_level="low",
                )
            )
            continue

        top_pair = max(related, key=lambda p: p.final_score)
        top_match_id = (
            top_pair.right_submission_id
            if top_pair.left_submission_id == sid
            else top_pair.left_submission_id
        )

        summaries.append(
            SubmissionSummary(
                analysis_job_id=job.analysis_job_id,
                submission_id=sid,
                top_match_submission_id=top_match_id,
                top_match_score=top_pair.final_score,
                final_score_risk_level=_risk_level(top_pair.final_score),
            )
        )

    logger.info("Saving %d summaries...", len(summaries))
    save_summaries(summaries)

    complete_job(job.analysis_job_id, "done")
    logger.info("Job %s completed successfully.", job.analysis_job_id)
