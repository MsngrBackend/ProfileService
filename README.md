# Profile Service

Profile microservice for the Messenger backend.  
Handles profiles, contacts, favorites, privacy settings, and notification preferences.

## Stack

- **Go 1.23**
- **PostgreSQL 16** — profiles, contacts, favorites, privacy, notifications
- **MinIO** — avatar storage (S3-compatible)
- **goose** — database migrations
- **Docker + Compose** — local development

## Requirements

- [Docker](https://docs.docker.com/get-docker/) with Compose plugin
- Go 1.22+ (only needed for local development outside Docker)

## Run

```bash
git clone https://github.com/MsngrBackend/ProfileService
cd ProfileService

docker compose up --build
