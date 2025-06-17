# blog4

tokuhirom の個人的なブログサービスです｡

## Quick Start

```bash
# Install dependencies
brew install go-task

# One-time setup
task gen
cd server/admin/frontend && npm install

# Run development servers
task dev
```

## Development Guide

### Prerequisites
- Go 1.21+
- Node.js 18+
- MySQL/MariaDB
- go-task (install with `brew install go-task`)

### Common Commands

```bash
# Generate all code (TypeSpec, SQLC, OpenAPI)
task gen

# Watch mode for code generation
task --watch gen

# Run all dev servers
task dev

# Frontend development only
task frontend

# Build production assets
task frontend-build
go build .

# Docker operations
task docker-build
task docker-run
```

### Running the Server

```bash
# Run the main server (includes both public and admin)
go run main.go

# The server starts on http://localhost:8181/
# - Public blog: http://localhost:8181/
# - Admin panel: http://localhost:8181/admin

# Debug with Delve
dlv debug main.go

# Build and run
go build -o blog4
./blog4
```

### Project Structure

```
blog4/
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
- Frontend: Svelte components in `/server/admin/frontend/`
- All code is auto-generated from TypeSpec and SQL files

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
