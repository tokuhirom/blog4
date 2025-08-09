# Docker Compose Development Setup

This setup provides a complete development environment with MariaDB, Go backend (with Air hot-reload), and Vite frontend.

## Prerequisites

- Docker and Docker Compose
- Git

## Quick Start

1. Start all services:
   ```bash
   docker-compose up -d
   ```

2. Initialize the database (first time only):
   ```bash
   docker-compose exec -T mariadb mysql -uroot -prootpassword blog4 < db/init/01-schema.sql
   ```

3. Load dummy data (optional):
   ```bash
   docker-compose exec -T mariadb mysql -uroot -prootpassword --default-character-set=utf8mb4 blog4 < db/init/02-dummy-data.sql
   ```

4. Access the application:
   - Frontend (Admin): http://localhost:6173
   - Backend API: http://localhost:8181
   - Database: localhost:3306
   - MinIO Console: http://localhost:9001 (login: minioadmin/minioadmin)
   - MinIO API: http://localhost:9000

## Services

### MariaDB (v10.11.13)
- Port: 3306
- Database: blog4 (configurable via MYSQL_DATABASE)
- User: blog4user (configurable via MYSQL_USER)
- Data persisted in Docker volume

### Backend (Go with Air)
- Port: 8181
- Hot-reload enabled with Air
- Source code mounted as volume
- Automatically rebuilds on file changes

### Frontend (Vite)
- Port: 6173
- Hot-reload enabled
- Source code mounted as volume

### MinIO (S3-compatible storage)
- API Port: 9000
- Console Port: 9001
- Credentials: minioadmin/minioadmin
- Automatically creates buckets: blog4-attachments, blog4-backup
- Data persisted in Docker volume

## Common Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f mariadb

# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes database data)
docker-compose down -v

# Rebuild services
docker-compose build

# Access MariaDB shell
docker-compose exec mariadb mysql -u root -p

# Run backend tests
docker-compose exec backend go test ./...

# Generate code (TypeSpec, SQLC, etc.)
docker-compose exec backend task gen
```

## Development Workflow

1. Backend changes: Air automatically detects changes and rebuilds
2. Frontend changes: Vite automatically detects changes and reloads
3. Database schema changes: Update `db/init/01-schema.sql` and restart MariaDB container

## Troubleshooting

### Backend not starting
- Check logs: `docker-compose logs backend`
- Ensure database is healthy: `docker-compose ps`
- Verify environment variables in docker-compose.yml

### Frontend not accessible
- Check if port 6173 is already in use
- Verify npm install completed: `docker-compose logs frontend`

### Database connection issues
- Wait for MariaDB to be fully initialized (check health status)
- Verify credentials match in docker-compose.yml