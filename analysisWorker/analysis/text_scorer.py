
from __future__ import annotations

import logging
from typing import Any

import numpy as np
from sentence_transformers import SentenceTransformer

import config
from analysis.utils import chunk_text, extract_author_text

logger = logging.getLogger(__name__)

_model: SentenceTransformer | None = None


_WORDS_PER_TOKEN = 0.75


def get_model() -> SentenceTransformer:
    global _model
    if _model is None:
        logger.info("Loading embedding model: %s", config.MODEL_NAME)
        _model = SentenceTransformer(config.MODEL_NAME)
        logger.info("Model loaded.")
    return _model


def _embed(texts: list[str]) -> np.ndarray:
    model = get_model()
    prefixed = [f"passage: {t}" for t in texts]
    embeddings = model.encode(prefixed, normalize_embeddings=True, show_progress_bar=False)
    return np.array(embeddings)


def compute_embedding(structured_data: dict[str, Any]) -> np.ndarray | None:
   
    text = extract_author_text(structured_data)
    if not text.strip():
        logger.warning("No author text found in structured_data — text_score will be 0.")
        return None

    chunk_size_words = int(config.CHUNK_SIZE_TOKENS * _WORDS_PER_TOKEN)
    overlap_words = int(config.CHUNK_OVERLAP_TOKENS * _WORDS_PER_TOKEN)

    chunks = chunk_text(text, chunk_size=chunk_size_words, overlap=overlap_words)
    if not chunks:
        return None

    embeddings = _embed(chunks)         
    mean_vec = embeddings.mean(axis=0)  

    # Re-normalise the mean
    norm = np.linalg.norm(mean_vec)
    if norm > 0:
        mean_vec = mean_vec / norm

    return mean_vec


def cosine_similarity(vec_a: np.ndarray, vec_b: np.ndarray) -> float:
    """Cosine similarity for two L2-normalised vectors."""
    score = float(np.dot(vec_a, vec_b))
   
    return max(0.0, min(1.0, score))


def compute_text_score(
    emb_a: np.ndarray | None,
    emb_b: np.ndarray | None,
) -> float:
    if emb_a is None or emb_b is None:
        return 0.0
    return cosine_similarity(emb_a, emb_b)
