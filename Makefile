.PHONY: sqlc-admin sqlc-public sqlc gen docker-build biome-check biome-fix help

# SQLC generation for admin database
sqlc-admin:
	go tool sqlc generate -f ./db/admin/sqlc-admin.yml

# SQLC generation for public database
sqlc-public:
	go tool sqlc generate -f ./db/public/sqlc-public.yml

# Run SQLC generation for both admin and public databases
sqlc: sqlc-admin sqlc-public

# Generate all code (SQLC only)
gen: sqlc

# Build production Docker image
docker-build:
	docker build -t blog4 .

# Run biome check (lint and format check) via Docker
biome-check:
	docker run --rm -v $(PWD):/app -w /app ghcr.io/biomejs/biome:latest check admin/static/sw.js

# Run biome fix (lint and format fix) via Docker
biome-fix:
	docker run --rm -v $(PWD):/app -w /app ghcr.io/biomejs/biome:latest check --write admin/static/sw.js

# Show available targets
help:
	@echo "Available targets:"
	@echo "  sqlc-admin    - Generate SQLC code for admin database"
	@echo "  sqlc-public   - Generate SQLC code for public database"
	@echo "  sqlc          - Generate SQLC code for both databases"
	@echo "  gen           - Generate all code (SQLC only)"
	@echo "  docker-build  - Build production Docker image"
	@echo "  biome-check   - Run biome lint and format check via Docker"
	@echo "  biome-fix     - Run biome lint and format fix via Docker"
	@echo "  help          - Show this help message"
