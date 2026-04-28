import psycopg2
import psycopg2.pool
import psycopg2.extras

import config


psycopg2.extras.register_uuid()

_pool: psycopg2.pool.SimpleConnectionPool | None = None


def get_pool() -> psycopg2.pool.SimpleConnectionPool:
    global _pool
    if _pool is None:
        _pool = psycopg2.pool.SimpleConnectionPool(
            minconn=1,
            maxconn=5,
            dsn=config.DATABASE_URL,
        )
    return _pool


def get_conn():
    return get_pool().getconn()


def put_conn(conn) -> None:
    get_pool().putconn(conn)


def close_pool() -> None:
    global _pool
    if _pool is not None:
        _pool.closeall()
        _pool = None
