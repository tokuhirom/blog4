# blog4

tokuhirom の個人的なブログサービスです｡

## Quick Start

```bash
# Install Docker and Docker Compose
# Then run the entire stack
docker-compose up

# The services are available at:
# - Frontend: http://localhost:6173
# - Backend API: http://localhost:8181
# - MariaDB: localhost:3306
# - MinIO Console: http://localhost:9001
# - MinIO API: http://localhost:9000
```

## Development Guide

### Prerequisites
- Docker and Docker Compose

### Common Commands

```bash
# Start all services with Docker Compose
docker-compose up

# Start in background
docker-compose up -d

# Stop all services
docker-compose down

# View logs
docker-compose logs -f [service-name]

# Build production Docker image
task docker-build
```

### Port Numbers

When running with Docker Compose, the following ports are exposed:

| Service | Port | Description |
|---------|------|-------------|
| Frontend | 6173 | Svelte admin UI |
| Backend | 8181 | Go API server |
| MariaDB | 3306 | Database |
| MinIO API | 9000 | S3-compatible storage |
| MinIO Console | 9001 | MinIO web interface |


### Project Structure

```
blog4/
├── cmd/
│   └── blog4/       # Main application entry point
├── db/               # Database schemas and queries
│   ├── admin/       # Admin database (write operations)
│   └── public/      # Public database (read operations)
├── server/
│   ├── admin/       # Admin API server
│   │   └── frontend/  # Svelte admin UI
│   └── public/      # Public API server
├── typespec/        # API definitions
├── markdown/        # Markdown processing
└── workers/         # Background jobs
```

### Making Changes

For detailed development instructions, see [CLAUDE.md](./CLAUDE.md).

Key points:
- API changes: Edit TypeSpec files in `/typespec/`
- Database queries: Add SQL to `/db/*/queries/`
- Frontend: Svelte components in `/web/admin/`
- API client code is generated using Orval (no Java required)
- All code is auto-generated from TypeSpec and SQL files

### Code Formatting

#### Go Code
This project uses `goimports` for consistent Go code formatting:

```bash
# Format all Go files
task fmt

# Check formatting without making changes
task lint

# Install pre-commit hooks (optional)
pre-commit install
```

Import ordering follows `goimports` style:
1. Standard library imports
2. Third-party imports
3. Local imports (github.com/tokuhirom/blog4)

#### Frontend Code
This project uses Biome for frontend formatting and linting:

```bash
# Format frontend code
task biome-format

# Lint frontend code
task biome-lint

# Run both format and lint
task frontend-fmt

# Or use pnpm directly in the frontend directory
cd web/admin
pnpm run format  # Format code
pnpm run lint    # Lint code
```

### Testing

```bash
# Run all tests
go test ./...

# Run specific test
go test -run TestFunctionName ./path/to/package

# Run with race detector
go test -race ./...
```

### Contributing

1. Create a feature branch
2. Make your changes
3. Run tests and linters
4. Create a pull request

For AI-assisted development with Claude Code, refer to [CLAUDE.md](./CLAUDE.md) for specific workflow instructions.

## Deployment

The application is deployed on Sakura Cloud App Run using Docker. See `Dockerfile` for the build configuration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
