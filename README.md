# blog4

tokuhirom の個人的なブログサービスです｡

## Quick Start

```bash
# Install Docker and Docker Compose
# Then run the entire stack
docker-compose up

# The services are available at:
# - Admin Interface: http://localhost:8181/admin/entries/search
# - Backend API: http://localhost:8181
# - MariaDB: localhost:3306
# - LocalStack (S3): http://localhost:4566
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
make docker-build

# Access MariaDB console
make db       # as blog4user
make db-root  # as root
```

### Port Numbers

When running with Docker Compose, the following ports are exposed:

| Service | Port | Description |
|---------|------|-------------|
| Backend | 8181 | Go API server with admin interface |
| MariaDB | 3306 | Database |
| LocalStack | 4566 | S3-compatible storage (LocalStack) |


### Project Structure

```
blog4/
├── cmd/
│   └── blog4/       # Main application entry point
├── db/               # Database schemas and queries
│   ├── admin/       # Admin database (write operations)
│   └── public/      # Public database (read operations)
├── internal/
│   ├── admin/       # Admin handlers and routes (Preact-based)
│   ├── public/      # Public handlers
│   └── middleware/  # HTTP middleware
├── admin/
│   ├── templates/   # HTML templates for admin interface
│   └── static/      # Static assets (CSS, JS, icons)
├── markdown/        # Markdown processing
└── scripts/         # Build and utility scripts
```

### Making Changes

For detailed development instructions, see [CLAUDE.md](./CLAUDE.md).

Key points:
- Database queries: Add SQL to `/db/*/queries/`, then run `make sqlc-admin` or `make sqlc-public`
- Admin interface: HTML templates in `/admin/templates/` with Preact apps
- Static assets: CSS, JavaScript, icons in `/admin/static/`
- Handlers: Go handlers in `/internal/admin/` and `/internal/public/`
- PWA: Manifest and Service Worker for Web Share Target support

### Code Formatting

#### Go Code
This project uses `goimports` for consistent Go code formatting:

```bash
# Format Go files
goimports -w .

# Install pre-commit hooks (optional)
pre-commit install
```

Import ordering follows `goimports` style:
1. Standard library imports
2. Third-party imports
3. Local imports (github.com/tokuhirom/blog4)

#### JavaScript Code
This project uses Biome for JavaScript linting and formatting:

```bash
# Check JavaScript files
make biome-check

# Fix JavaScript files
make biome-fix
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
