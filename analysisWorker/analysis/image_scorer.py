

from __future__ import annotations

import logging
import re
from typing import Any

import numpy as np

from analysis.text_scorer import compute_text_score, _embed
from analysis.utils import chunk_text
import config

logger = logging.getLogger(__name__)

_MEANINGFUL_CONTENT_RE = re.compile(r"[а-яёА-ЯЁa-zA-Z]{3,}")


def _extract_ocr_text(structured_data: dict[str, Any]) -> str:
   
    image_texts: list[str] = (structured_data.get("rawParts") or {}).get("imageTexts") or []

    meaningful: list[str] = []
    for text in image_texts:
        stripped = text.strip()
        if len(stripped) < 10:
            continue
       
        if not _MEANINGFUL_CONTENT_RE.search(stripped):
            continue
        meaningful.append(stripped)

    return "\n".join(meaningful)


def compute_image_embedding(structured_data: dict[str, Any]) -> np.ndarray | None:
   
    text = _extract_ocr_text(structured_data)
    if not text.strip():
        return None

    chunk_size_words = int(config.CHUNK_SIZE_TOKENS * 0.75)
    overlap_words = int(config.CHUNK_OVERLAP_TOKENS * 0.75)

    chunks = chunk_text(text, chunk_size=chunk_size_words, overlap=overlap_words)
    if not chunks:
        return None

    embeddings = _embed(chunks)
    mean_vec = embeddings.mean(axis=0)

    norm = np.linalg.norm(mean_vec)
    if norm > 0:
        mean_vec = mean_vec / norm

    return mean_vec


def compute_image_score(
    emb_a: np.ndarray | None,
    emb_b: np.ndarray | None,
) -> float:
    return compute_text_score(emb_a, emb_b)
