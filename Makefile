.PHONY: sqlc-admin sqlc-public sqlc gen docker-build biome-check biome-fix db db-root e2e-install e2e-install-browsers e2e-test e2e-test-ui help

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

# Access MariaDB console as blog4user
db:
	docker-compose exec mariadb mysql -ublog4user -pblog4password blog4

# Access MariaDB console as root
db-root:
	docker-compose exec mariadb mysql -uroot -prootpassword

# Install e2e test dependencies
e2e-install:
	cd e2e && npm install

# Install Playwright browsers for e2e tests
e2e-install-browsers:
	cd e2e && npx playwright install chromium

# Run e2e tests (requires docker-compose services to be running)
e2e-test:
	cd e2e && npm test

# Run e2e tests in UI mode
e2e-test-ui:
	cd e2e && npm run test:ui

# Show available targets
help:
	@echo "Available targets:"
	@echo "  sqlc-admin             - Generate SQLC code for admin database"
	@echo "  sqlc-public            - Generate SQLC code for public database"
	@echo "  sqlc                   - Generate SQLC code for both databases"
	@echo "  gen                    - Generate all code (SQLC only)"
	@echo "  docker-build           - Build production Docker image"
	@echo "  biome-check            - Run biome lint and format check via Docker"
	@echo "  biome-fix              - Run biome lint and format fix via Docker"
	@echo "  db                     - Access MariaDB console as blog4user"
	@echo "  db-root                - Access MariaDB console as root"
	@echo "  e2e-install            - Install e2e test dependencies"
	@echo "  e2e-install-browsers   - Install Playwright browsers"
	@echo "  e2e-test               - Run e2e tests"
	@echo "  e2e-test-ui            - Run e2e tests in UI mode"
	@echo "  help                   - Show this help message"
