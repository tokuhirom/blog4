# Stage 1: Build the frontend
FROM node:24 AS frontend-builder
RUN corepack enable && corepack prepare pnpm@10.12.2 --activate
WORKDIR /app
# Copy workspace configuration and lockfile
COPY pnpm-workspace.yaml pnpm-lock.yaml package.json ./
COPY web/admin/package.json ./web/admin/
# Install dependencies for the entire workspace
RUN pnpm install --frozen-lockfile
# Copy the frontend source code
COPY web/admin/ ./web/admin/
ENV VITE_API_BASE_URL=/admin/api
# Build from the web/admin directory
RUN cd web/admin && pnpm run build

# Stage 2: Build the Go backend
FROM golang:1.25 AS backend-builder
RUN apt-get update && apt-get install -y git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
# Generate build info
RUN bash scripts/generate-build-info.sh
RUN go build -o /app/blog4 ./cmd/blog4

# Stage 3: Final stage
FROM ubuntu:24.04
WORKDIR /app
COPY --from=backend-builder /app/blog4 /app/
COPY --from=backend-builder /app/build-info.json /app/
COPY --from=frontend-builder /app/web/admin/dist /app/web/admin/dist
COPY web/static /app/web/static
RUN apt-get update && apt-get install -y tzdata mysql-client openssl ca-certificates

ARG GIT_HASH
ENV GIT_HASH=$GIT_HASH

EXPOSE 8181
CMD ["/app/blog4"]
