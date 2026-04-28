
from __future__ import annotations

from typing import Any

import config
from analysis.utils import extract_numbers


def _numbers_match(a: float, b: float) -> bool:
    
    if a == 0.0 and b == 0.0:
        return True
    if a == 0.0 or b == 0.0:
        return False
    return abs(a - b) / max(abs(a), abs(b)) <= config.NUMBER_TOLERANCE


def _jaccard(nums_a: list[float], nums_b: list[float]) -> float:
    
    if not nums_a or not nums_b:
        return 0.0

    matched_b = [False] * len(nums_b)
    intersection = 0

    for a in nums_a:
        for j, b in enumerate(nums_b):
            if not matched_b[j] and _numbers_match(a, b):
                intersection += 1
                matched_b[j] = True
                break

    union = len(nums_a) + len(nums_b) - intersection
    return intersection / union if union > 0 else 0.0


def compute_numbers(structured_data: dict[str, Any]) -> list[float]:
    
    return extract_numbers(structured_data)


def compute_calculation_score(
    nums_a: list[float],
    nums_b: list[float],
) -> float:
    
    if (
        len(nums_a) < config.MIN_NUMBERS_FOR_CALC_SCORE
        or len(nums_b) < config.MIN_NUMBERS_FOR_CALC_SCORE
    ):
        return 0.0

    return _jaccard(nums_a, nums_b)
