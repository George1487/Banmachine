# BanMachine

BanMachine is a system for checking student lab reports for suspicious similarity.

The flow is simple: a teacher creates a lab, students upload `.docx` reports, files are parsed in the background, and the system calculates similarity scores inside the same lab.

## How it works

- Teachers create labs with title, description, and deadline.
- Students upload `.docx` submissions in the frontend.
- API stores metadata in PostgreSQL and files in MinIO.
- `ingest-worker` parses uploaded DOCX files and writes parsed data.
- `analysis-worker` runs similarity analysis for lab submissions.
- Teachers see job status, top matches, and risk level in the dashboard.

## Built with

- Frontend: Vue 3, Vite, Pinia, Vue Router, Nginx
- API: ASP.NET Core 9 (`ApiService`)
- Ingest pipeline: Go (`ingestWorker`)
- Analysis pipeline: Python 3.11 (`analysisWorker`)
- Data and storage: PostgreSQL, MinIO
- Infra: Docker Compose

## Local Run

### Prerequisites

- Docker
- Docker Compose

### Start

Linux/macOS:
```bash
cp .env.example .env
docker compose up -d --build --remove-orphans
```

Windows PowerShell:
```powershell
Copy-Item .env.example .env
docker compose up -d --build --remove-orphans
```

### Useful URLs

- Frontend: `http://localhost:8080`
- Swagger: `http://localhost:8082/swagger`
- MinIO Console: `http://localhost:9001`

## Environment Variables

Use `.env.example` as template and create local `.env`.

| Variable | Purpose | Default in template |
| --- | --- | --- |
| `POSTGRES_DB` | PostgreSQL database name | `banmachine_db` |
| `POSTGRES_USER` | PostgreSQL user | `banmachine_user` |
| `POSTGRES_PASSWORD` | PostgreSQL password | `replace_with_strong_db_password` |
| `POSTGRES_PORT` | PostgreSQL host port | `5432` |
| `PGDATA` | PostgreSQL data path in container | `/var/lib/postgresql/data/pgdata` |
| `MINIO_ROOT_USER` | MinIO access key | `replace_with_minio_user` |
| `MINIO_ROOT_PASSWORD` | MinIO secret key | `replace_with_strong_minio_password` |
| `MINIO_BUCKET` | bucket for submissions | `submissions` |
| `MINIO_API_PORT` | MinIO API port | `9000` |
| `MINIO_CONSOLE_PORT` | MinIO console port | `9001` |
| `DOTNET_PORT` | .NET API host port | `8082` |
| `ASPNETCORE_ENVIRONMENT` | ASP.NET Core environment | `Development` |
| `AUTH_SECRET_KEY` | JWT signing key | `replace_with_very_long_random_secret_key_min_32_chars` |
| `IMAGE_REGISTRY` | image registry | `ghcr.io` |
| `IMAGE_NAMESPACE` | image namespace | `bandabenzogang/banmashine` |
| `FRONTEND_IMAGE_TAG` | frontend image tag | `latest` |
| `DOTNET_IMAGE_TAG` | dotnet image tag | `latest` |
| `INGEST_IMAGE_TAG` | ingest image tag | `latest` |

## Data model and migrations

Main SQL migrations are in `migrations/`:

- `001_init.up.sql` / `001_init.down.sql`
- `002_ingest_jobs_runtime_columns.up.sql` / `002_ingest_jobs_runtime_columns.down.sql`

Core tables created by migrations:

- `users`
- `labs`
- `submissions`
- `parsed_submissions`
- `ingest_jobs`
- `analysis_jobs`
- `analysis_job_snapshots`
- `pairwise_similarities`
- `submission_analysis_summaries`

## Security notes

- Runtime secrets are not committed to git.
- `.env` and local appsettings are ignored by `.gitignore`.
- Use `ApiService/Web/appsettings.template.json` only as template for local setup.

## What I worked on (ApiService and migrations)

I was responsible for the backend API and database migrations.

### ApiService

- Built the API service using ASP.NET Core 9 and layered structure.
- Structured backend modules as `Application`, `Domain`, `Infrastructure`, `Ports`, `Presentation`, and `Web`.
- Implemented API endpoints in `AuthenticationController`, `LabsController`, `SubmissionsController`, and `AnalysisController`.
- Wired service and repository layers for labs, submissions, analysis jobs, ingest jobs, summaries, and pairwise similarities.
- Configured PostgreSQL integration in `Infrastructure` (`AppDbContext` + repositories).
- Configured JWT authentication and authorization in API startup.
- Configured MinIO integration for submission file storage.
- Enabled Swagger for API inspection and testing.

### Migrations

- Wrote base migration `001_init` with the initial schema.
- Added indexes for key read paths and job/status queries.
- Added migration `002_ingest_jobs_runtime_columns` to extend ingest runtime tracking.
- Kept both `up` and `down` scripts for reproducible apply/rollback flow.
