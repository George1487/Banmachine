

from __future__ import annotations

import json
import sys
from pathlib import Path


sys.path.insert(0, str(Path(__file__).parent))

from analysis.calculation_scorer import compute_calculation_score, compute_numbers
from analysis.image_scorer import compute_image_embedding, compute_image_score
from analysis.text_scorer import compute_embedding, compute_text_score
import config


PARSED_DIR = Path(__file__).parent / "data" / "parsed"


def load_parsed_file(path: Path) -> dict:
    """Read a parsed .txt file and extract the JSON structured_data block."""
    text = path.read_text(encoding="utf-8")
    separator = "STRUCTURED DATA:"
    if separator not in text:
        raise ValueError(f"No '{separator}' section found in {path.name}")

    raw_text, json_part = text.split(separator, 1)
    # Use raw_decode so trailing whitespace/chars after the JSON object are ignored
    decoder = json.JSONDecoder()
    structured_data, _ = decoder.raw_decode(json_part.strip())
    return {"raw_text": raw_text, "structured_data": structured_data, "name": path.name}


def main() -> None:
    files = sorted(PARSED_DIR.glob("*.txt"))
    if len(files) < 2:
        print("Need at least 2 parsed files in data/parsed/")
        sys.exit(1)

    print(f"Loading {len(files)} parsed submissions...")
    subs = [load_parsed_file(f) for f in files]

    print("Pre-computing embeddings and number sets...")
    for sub in subs:
        sd = sub["structured_data"]
        sub["text_emb"] = compute_embedding(sd)
        sub["img_emb"] = compute_image_embedding(sd)
        sub["numbers"] = compute_numbers(sd)
        print(
            f"  {sub['name']}: "
            f"numbers={len(sub['numbers'])}, "
            f"text_emb={'ok' if sub['text_emb'] is not None else 'none'}, "
            f"img_emb={'ok' if sub['img_emb'] is not None else 'none'}"
        )

    print("\n--- Pairwise scores ---")
    for i, a in enumerate(subs):
        for b in subs[i + 1:]:
            text_score = compute_text_score(a["text_emb"], b["text_emb"])
            calc_score = compute_calculation_score(a["numbers"], b["numbers"])
            img_score = compute_image_score(a["img_emb"], b["img_emb"])
            final = (
                config.TEXT_WEIGHT * text_score
                + config.CALC_WEIGHT * calc_score
                + config.IMG_WEIGHT * img_score
            )
            risk = (
                "HIGH" if final >= config.HIGH_THRESHOLD
                else "MEDIUM" if final >= config.MEDIUM_THRESHOLD
                else "low"
            )
            print(
                f"\n{a['name']}  vs  {b['name']}\n"
                f"  text_score:  {text_score:.6f}\n"
                f"  calc_score:  {calc_score:.6f}  (numbers: {len(a['numbers'])} vs {len(b['numbers'])})\n"
                f"  img_score:   {img_score:.6f}\n"
                f"  final_score: {final:.6f}  [{risk}]"
            )


if __name__ == "__main__":
    main()
