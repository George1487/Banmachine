# BanMachine

Security-first version of the project.

## Services

This repository is started by one `docker-compose.yml` and includes:
- `db` (PostgreSQL)
- `migrate` (DB migrations)
- `analysis-worker` (Python)
- `minio`
- `minio-init`
- `ingest-worker` (Go)
- `frontend` (Vue + Nginx)
- `backend-dotnet` (.NET)

## Secure setup

1. Create local env file from template:
```bash
cp .env.example .env
```

2. Open `.env` and replace all `replace_with_*` values with your real secrets.

3. For local `dotnet run` outside Docker, create local appsettings from template:
```bash
cp ApiService/Web/appsettings.template.json ApiService/Web/appsettings.Development.json
```
Then fill it with your local values.

## Run

```bash
docker compose up -d --build --remove-orphans
```

## Security notes

- Do not commit `.env` files.
- Do not commit `ApiService/Web/appsettings*.json` with real keys.
- Rotate keys immediately if they were ever pushed to a remote repository.
