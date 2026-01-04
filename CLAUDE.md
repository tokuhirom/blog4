# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Important Git Workflow Instructions

**CRITICAL: Always follow these git workflow rules:**
1. **Create a new branch for EVERY commit** - Never commit directly to main
2. **Create a PR after EVERY commit** - Each commit should have its own PR
3. **Do NOT enable auto-merge** - Manually review and merge PRs
4. **Branch naming convention** - Use descriptive names like `fix/bug-description` or `feat/feature-name`

Example workflow:
```bash
# 1. Create new branch
git checkout -b fix/some-bug

# 2. Make changes and commit
git add .
git commit -m "fix: description"

# 3. Push branch
git push -u origin fix/some-bug

# 4. Create PR
gh pr create --title "fix: description" --body "..."

# 5. Manually review and merge PR
```

## Project Overview

Blog4 is a full-stack blog application with an HTMX-based admin interface, built with Go backend and HTML templates. It features wiki-style linking, Amazon product integration, automated image generation, and PWA support with Web Share Target API for sharing content from Android devices.

## Development Setup

The project uses Docker Compose for local development:

```bash
# Start all services
docker-compose up

# Services and ports:
# - Backend API: http://localhost:8181
# - MariaDB: localhost:3306
# - LocalStack (S3): http://localhost:4566
```

### LocalStack S3 Storage

The project uses LocalStack for S3-compatible storage in local development:

- **Endpoint**: http://localhost:4566
- **Buckets**:
  - `blog4-attachments` (public) - for uploaded images and files
  - `blog4-backup` (private) - for database backups
- **Credentials**: `test` / `test` (LocalStack defaults)

```bash
# List buckets
aws --endpoint-url=http://localhost:4566 s3 ls

# List files in attachments bucket
aws --endpoint-url=http://localhost:4566 s3 ls s3://blog4-attachments

# Download a file
aws --endpoint-url=http://localhost:4566 s3 cp s3://blog4-attachments/attachments/2024/01/01/file.jpg ./
```

## Build Commands

```bash
# Generate SQLC code for admin database
make sqlc-admin

# Generate SQLC code for public database
make sqlc-public

# Generate all code (both databases)
make gen

# Build production Docker image
make docker-build

# JavaScript linting and formatting
make biome-check  # Check JavaScript files
make biome-fix    # Fix JavaScript files

# Database access
make db       # Access MariaDB console as blog4user
make db-root  # Access MariaDB console as root
```

## Architecture

### Code Generation Pipeline
1. **SQL** (`/db/*/queries/*.sql`) → Go code via SQLC
2. **Templates** (`/admin/templates/*.html`) → Rendered by Go html/template

### Database Structure
- **MySQL/MariaDB** with two schemas:
  - `admindb`: Admin operations (write)
  - `publicdb`: Public queries (read)
- Key tables: `entry`, `entry_image`, `visibility`
- Uses n-gram parser for Japanese full-text search

### Admin Interface
HTMX-based admin interface with server-side rendering:

**Key Routes** (`/internal/admin/htmx_router.go`):
- `/admin/entries/search` - Entry list page (GET)
- `/admin/entries/edit?path=...` - Entry edit page (GET)
- `/admin/entries/create` - Create new entry (POST)
- `/admin/entries/title` - Update entry title (POST, HTMX)
- `/admin/entries/body` - Update entry body (POST, HTMX)
- `/admin/entries/visibility` - Update visibility (POST, HTMX)
- `/admin/entries/delete` - Delete entry (DELETE, HTMX)
- `/admin/share-target` - Web Share Target endpoint (POST)

**Templates** (`/admin/templates/`):
- `layout.html` - Base layout with PWA meta tags
- `htmx_entries.html` - Entry list page
- `htmx_entry_edit.html` - Entry edit page with auto-save
- `htmx_login.html` - Login page

**Middleware** (`/internal/admin/htmx_router.go`):
- `NoCacheMiddleware` - Prevents caching of dynamic admin pages
- `GinSessionMiddleware` - Session authentication

### PWA Support
- **Manifest**: `/admin/manifest.webmanifest` - PWA configuration
- **Service Worker**: `/admin/static/sw.js` - Caching strategy and Web Share Target support
- **Web Share Target**: Android users can share content directly to Blog4
- **Cache Strategy**: Static assets cached, admin API endpoints always use network

## Key Features Implementation

### Wiki-Style Links
- Pattern: `[[Entry Title]]` in markdown
- Processed by `/markdown/wiki_link.go`
- Two-hop link tracking in `twohop.go`

### Amazon Product Links
- Pattern: `[asin:B00EXAMPLE]`
- Handled by `/server/admin/amazon.go`
- Uses PA-API5 for product data

### Entry Images
- Generated automatically via worker (`entry_image_worker.go`)
- Stored in S3-compatible storage
- Puppeteer-based screenshot generation

## Common Development Tasks

### Add Database Query
1. Add SQL query to `/db/admin/queries/*.sql` or `/db/public/queries/*.sql`
2. Run `make sqlc-admin` or `make sqlc-public`
3. Use generated code in Go handlers

### Run Specific Test
```bash
go test -run TestFunctionName ./path/to/package
```

## Environment Configuration

Docker Compose handles all environment variables automatically. Key variables include:
- `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_HOST`
- `S3_ATTACHMENTS_BUCKET_NAME`, `S3_ENDPOINT`
- `S3_ACCESS_KEY_ID`, `S3_SECRET_ACCESS_KEY`
- `AMAZON_PAAPI5_ACCESS_KEY`, `AMAZON_PAAPI5_SECRET_KEY`

See `docker-compose.yml` for the complete list and default values.

## Troubleshooting

### SQLC Generation Fails
- On macOS 15+, may encounter `strchrnul` compilation error
- This is a known issue with pg_query_go dependency

### TypeSpec Compilation Errors
- Ensure using TypeSpec v0.67.2+
- Object literals in decorators must use `#{}` syntax

### Frontend Build Issues
- Run `pnpm install` in `/web/admin/`
- Check Node.js version compatibility (v18+)

## Deployment
- Uses Docker with multi-stage build
- Deployed on Sakura Cloud App Run
- Health check endpoint: `/healthz`
- Backup runs daily via cron job to S3

## Git Workflow Best Practices

- after modify tsx or ts, run biome before the commit
- use mockgen for db testing
- before commit, run biome, go test.
- after send pr, sleep a while, and check the ci state. if it's failed, resolve the issue and commit & push again.
- run goimports before commit golang code.