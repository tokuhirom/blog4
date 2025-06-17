# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Blog4 is a full-stack blog application with an admin interface, built with Go backend and TypeScript/Svelte frontend. It features wiki-style linking, Amazon product integration, and automated image generation.

## Build Commands

```bash
# Install dependencies
brew install go-task

# One-time setup
task gen          # Generate all code (TypeSpec, SQLC, OpenAPI)
cd server/admin/frontend && npm install

# Development
task dev          # Run all dev servers and watchers
task frontend     # Run frontend dev server only
task --watch gen  # Watch mode for code generation

# Build
task frontend-build  # Build frontend assets
go build .          # Build Go binary

# Docker
task docker-build   # Build Docker image
task docker-run     # Run Docker container
```

## Architecture

### Code Generation Pipeline
1. **TypeSpec** (`/typespec/*.tsp`) → OpenAPI spec
2. **OpenAPI** → Go server code via Ogen (`go generate main.go`)
3. **OpenAPI** → TypeScript client via openapi-generator
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
- **Location**: `/server/admin/frontend`
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
4. Run `task openapi-client` to generate TypeScript client

### Add Database Query
1. Add SQL query to `/db/admin/queries/*.sql` or `/db/public/queries/*.sql`
2. Run `task sqlc-admin` or `task sqlc-public`
3. Use generated code in Go handlers

### Run Specific Test
```bash
go test -run TestFunctionName ./path/to/package
```

## Environment Configuration
Required environment variables (see `app.jsonnet` for full list):
- `MYSQL_USER`, `MYSQL_PASSWORD`, `MYSQL_HOST`
- `ENTRY_IMAGE_BUCKET_NAME`, `ENTRY_IMAGE_ENDPOINT`
- `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`
- `AMAZON_ACCESS_KEY_ID`, `AMAZON_SECRET_ACCESS_KEY`

## Troubleshooting

### SQLC Generation Fails
- On macOS 15+, may encounter `strchrnul` compilation error
- This is a known issue with pg_query_go dependency

### TypeSpec Compilation Errors
- Ensure using TypeSpec v0.67.2+
- Object literals in decorators must use `#{}` syntax

### Frontend Build Issues
- Run `npm install` in `/server/admin/frontend/`
- Check Node.js version compatibility (v18+)

## Deployment
- Uses Docker with multi-stage build
- Deployed on Sakura Cloud App Run
- Health check endpoint: `/healthz`
- Backup runs daily via cron job to S3