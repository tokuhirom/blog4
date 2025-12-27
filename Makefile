.PHONY: sqlc-admin sqlc-public sqlc gen docker-build help

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

# Show available targets
help:
	@echo "Available targets:"
	@echo "  sqlc-admin    - Generate SQLC code for admin database"
	@echo "  sqlc-public   - Generate SQLC code for public database"
	@echo "  sqlc          - Generate SQLC code for both databases"
	@echo "  gen           - Generate all code (SQLC only)"
	@echo "  docker-build  - Build production Docker image"
	@echo "  help          - Show this help message"
