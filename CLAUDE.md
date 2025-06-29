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

Blog4 is a full-stack blog application with an admin interface, built with Go backend and TypeScript/Svelte frontend. It features wiki-style linking, Amazon product integration, and automated image generation.

## Development Setup

The project now uses Docker Compose for local development:

```bash
# Start all services
docker-compose up

# Services and ports:
# - Frontend (Svelte admin): http://localhost:6173
# - Backend API: http://localhost:8181  
# - MariaDB: localhost:3306
# - MinIO (S3): http://localhost:9000
# - MinIO Console: http://localhost:9001
```

## Build Commands

```bash
# Build production Docker image
task docker-build
```

## Architecture

### Code Generation Pipeline
1. **TypeSpec** (`/typespec/*.tsp`) → OpenAPI spec
2. **OpenAPI** → Go server code via Ogen (`go generate ./cmd/blog4/main.go`)
3. **OpenAPI** → TypeScript client via Orval (no Java required)
4. **SQL** (`/db/*/queries/*.sql`) → Go code via SQLC

### Database Structure
- **MySQL/MariaDB** with two schemas:
  - `admindb`: Admin operations (write)
  - `publicdb`: Public queries (read)
- Key tables: `entry`, `entry_image`, `visibility`
- Uses n-gram parser for Japanese full-text search

### API Endpoints
All admin endpoints are prefixed with `/api/`:
- `/api/entries` - Entry CRUD operations
- `/api/entries/{path}/*` - Entry-specific operations (title, body, visibility, etc.)
- `/api/upload` - File upload

### Frontend Structure
- **Location**: `/web/admin`
- **Framework**: Svelte 5 with TypeScript
- **Key Components**:
  - `MarkdownEditor.svelte` - CodeMirror-based editor
  - `AdminEntryPage.svelte` - Entry editing page
- **API Client**: Auto-generated in `/src/generated-client/`

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

### Modify API
1. Edit TypeSpec files in `/typespec/`
2. Run `task tsp` to generate OpenAPI
3. Run `task ogen` to generate Go server
4. Run `task openapi-client` to generate TypeScript client with Orval

### Add Database Query
1. Add SQL query to `/db/admin/queries/*.sql` or `/db/public/queries/*.sql`
2. Run `task sqlc-admin` or `task sqlc-public`
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