from __future__ import annotations

import uuid
from dataclasses import dataclass
from typing import Any


@dataclass
class AnalysisJob:
    analysis_job_id: uuid.UUID
    lab_id: uuid.UUID
    status: str
    created_by: uuid.UUID


@dataclass
class SnapshotSubmission:
    """A parsed submission loaded for analysis."""

    submission_id: uuid.UUID
    student_id: uuid.UUID
    raw_text: str
    structured_data: dict[str, Any]


@dataclass
class PairwiseResult:
    """Result of comparing two submissions."""

    analysis_job_id: uuid.UUID
    lab_id: uuid.UUID
    left_submission_id: uuid.UUID
    right_submission_id: uuid.UUID
    text_score: float
    calculation_score: float
    images_score: float
    final_score: float


@dataclass
class SubmissionSummary:
    """Aggregated analysis summary for one submission."""

    analysis_job_id: uuid.UUID
    submission_id: uuid.UUID
    top_match_submission_id: uuid.UUID | None
    top_match_score: float | None
    final_score_risk_level: str  # low | medium | high
