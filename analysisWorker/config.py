import os


_database_url = os.environ.get("DATABASE_URL")
if not _database_url:
    _host = os.environ.get("POSTGRES_HOST", "localhost")
    _port = os.environ.get("POSTGRES_PORT", "5432")
    _db = os.environ.get("POSTGRES_DB", "banmachine")
    _user = os.environ.get("POSTGRES_USER", "postgres")
    _password = os.environ.get("POSTGRES_PASSWORD", "postgres")
    _database_url = f"postgresql://{_user}:{_password}@{_host}:{_port}/{_db}"

DATABASE_URL = _database_url

POLL_INTERVAL_SEC = int(os.environ.get("POLL_INTERVAL_SEC", "5"))


MODEL_NAME = os.environ.get("MODEL_NAME", "intfloat/multilingual-e5-large")


TEXT_WEIGHT = float(os.environ.get("TEXT_WEIGHT", "0.3"))
CALC_WEIGHT = float(os.environ.get("CALC_WEIGHT", "0.6"))
IMG_WEIGHT = float(os.environ.get("IMG_WEIGHT", "0.1"))


HIGH_THRESHOLD = float(os.environ.get("HIGH_THRESHOLD", "0.8"))
MEDIUM_THRESHOLD = float(os.environ.get("MEDIUM_THRESHOLD", "0.5"))


CHUNK_SIZE_TOKENS = int(os.environ.get("CHUNK_SIZE_TOKENS", "400"))
CHUNK_OVERLAP_TOKENS = int(os.environ.get("CHUNK_OVERLAP_TOKENS", "50"))


NUMBER_TOLERANCE = float(os.environ.get("NUMBER_TOLERANCE", "0.01"))  # ±1%
MIN_NUMBERS_FOR_CALC_SCORE = int(os.environ.get("MIN_NUMBERS_FOR_CALC_SCORE", "3"))
