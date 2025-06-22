# blog4

tokuhirom の個人的なブログサービスです｡

## Quick Start

```bash
# Install dependencies (macOS)
brew install go-task

# Install dependencies (Linux)
go install github.com/go-task/task/v3/cmd/task@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# One-time setup
npm install  # Install TypeSpec dependencies
cd web/admin && npm install
task gen

# Run development servers
task dev
```

## Development Guide

### Prerequisites
- Go 1.21+
- Node.js 18+
- MySQL/MariaDB
- go-task
  - macOS: `brew install go-task`
  - Linux: `go install github.com/go-task/task/v3/cmd/task@latest`
- sqlc: `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
- TypeSpec: Installed via `npm install` (see package.json)

### Common Commands

```bash
# Generate all code (TypeSpec, SQLC, OpenAPI client)
task gen

# Watch mode for code generation
task --watch gen

# Run all dev servers
task dev

# Frontend development only
task frontend

# Build production assets
task frontend-build
go build ./cmd/blog4

# Docker operations
task docker-build
task docker-run
```

### Running the Server

```bash
# Run the main server (includes both public and admin)
go run ./cmd/blog4

# The server starts on http://localhost:8181/
# - Public blog: http://localhost:8181/
# - Admin panel: http://localhost:8181/admin

# Debug with Delve
dlv debug ./cmd/blog4

# Build and run
go build -o blog4 ./cmd/blog4
./blog4
```

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

# Or use npm directly in the frontend directory
cd web/admin
npm run format  # Format code
npm run lint    # Lint code
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

Personal project - not open for external contributions.
