

import logging
import time
import traceback

import config
from analysis.pipeline import run_analysis
from analysis.text_scorer import get_model
from db.connection import close_pool
from db.queries import claim_pending_job, complete_job

logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s [%(levelname)s] %(name)s: %(message)s",
)
logger = logging.getLogger("analysis_worker")


def main() -> None:
    logger.info("AnalysisWorker starting up.")


    logger.info("Pre-loading embedding model...")
    get_model()
    logger.info("Model ready. Entering polling loop.")

    try:
        while True:
            job = claim_pending_job()

            if job is None:
                logger.debug("No pending jobs. Sleeping %ds.", config.POLL_INTERVAL_SEC)
                time.sleep(config.POLL_INTERVAL_SEC)
                continue

            logger.info("Claimed job %s.", job.analysis_job_id)
            try:
                run_analysis(job)
            except Exception:
                error_msg = traceback.format_exc()
                logger.error(
                    "Job %s failed:\n%s", job.analysis_job_id, error_msg
                )
                complete_job(job.analysis_job_id, "failed", error_msg[:2000])
    finally:
        close_pool()
        logger.info("AnalysisWorker shut down.")


if __name__ == "__main__":
    main()
