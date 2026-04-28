

from __future__ import annotations

import re
from typing import Any

import config


_TRIVIAL_NUMBERS: frozenset[float] = frozenset(
    [
        0.0, 1.0, 2.0, 3.0, 4.0, 5.0, 6.0, 7.0, 8.0, 9.0,
        10.0, 100.0, 1000.0,
        3.14, 3.141, 3.1415, 3.14159,  
        9.8, 9.81, 9.807,             
        6.67e-11, 6.674e-11,           
        1.38e-23, 1.381e-23,          
        6.022e23, 6.02e23,            
        1.6e-19, 1.602e-19,          
    ]
)


_TEMPLATE_PHRASES: tuple[str, ...] = (
    "рабочий протокол",
    "отчет по лабораторной работе",
    "цель работы",
    "задачи, решаемые при выполнении",
    "объект исследования",
    "метод экспериментального исследования",
    "рабочие формулы и исходные данные",
    "измерительные приборы",
    "схема установки",
    "результаты прямых измерений",
    "расчет результатов косвенных измерений",
    "расчет погрешностей",
    "окончательные результаты",
    "перечень схем",
    "приложение",
)

_NUMBER_RE = re.compile(r"-?\d+[.,]\d+(?:[eE][+-]?\d+)?|-?\d{4,}")


def _is_trivial(value: float) -> bool:
    for tv in _TRIVIAL_NUMBERS:
        if tv == 0.0:
            if value == 0.0:
                return True
            continue
        if abs(value - tv) / abs(tv) < 0.001:
            return True
    return False


def extract_numbers(structured_data: dict[str, Any]) -> list[float]:
   
    raw_numbers: list[float] = []

   
    for formula in structured_data.get("formulas") or []:
        normalized = formula.get("normalized", "")
        for match in _NUMBER_RE.findall(normalized):
            try:
                raw_numbers.append(float(match.replace(",", ".")))
            except ValueError:
                pass

   
    for para in (structured_data.get("rawParts") or {}).get("paragraphs") or []:
        if "=" not in para:
            continue
        for match in _NUMBER_RE.findall(para):
            try:
                raw_numbers.append(float(match.replace(",", ".")))
            except ValueError:
                pass

   
    result: list[float] = []
    for num in raw_numbers:
        if _is_trivial(num):
            continue
      
        already = False
        for existing in result:
            if existing == 0.0:
                continue
            if abs(num - existing) / abs(existing) <= config.NUMBER_TOLERANCE:
                already = True
                break
        if not already:
            result.append(num)

    return result


def _is_template_line(text: str) -> bool:
    lower = text.lower()
    return any(phrase in lower for phrase in _TEMPLATE_PHRASES)


def extract_author_text(structured_data: dict[str, Any]) -> str:
  
    authored_paragraphs: list[str] = []

   
    in_author_section = False
    for section in structured_data.get("sections") or []:
        title = section.get("title", "").lower()
       
        is_author = any(
            kw in title
            for kw in ("вывод", "анализ результат", "окончательны", "заключен")
        )
        is_template = _is_template_line(title)

        if is_author:
            in_author_section = True
        elif is_template:
            in_author_section = False

        if in_author_section or is_author:
            for para in section.get("paragraphs") or []:
                if len(para.strip()) >= 20 and not _is_template_line(para):
                    authored_paragraphs.append(para.strip())

    
    if not authored_paragraphs:
        for para in (structured_data.get("rawParts") or {}).get("paragraphs") or []:
            stripped = para.strip()
            if len(stripped) >= 20 and not _is_template_line(stripped):
                authored_paragraphs.append(stripped)

    return "\n".join(authored_paragraphs)


def chunk_text(text: str, chunk_size: int, overlap: int) -> list[str]:
   
    words = text.split()
    if not words:
        return []

    chunks: list[str] = []
    start = 0
    while start < len(words):
        end = min(start + chunk_size, len(words))
        chunks.append(" ".join(words[start:end]))
        if end == len(words):
            break
        start += chunk_size - overlap

    return chunks
