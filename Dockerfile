# Stage 1: Build the frontend
FROM node:22 AS frontend-builder
WORKDIR /app
COPY internal/admin/frontend/package*.json ./
RUN npm install
COPY internal/admin/frontend/ ./
RUN npm run build

# Stage 2: Build the Go backend
FROM golang:1.24 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /app/server ./cmd/blog4

# Stage 3: Final stage
FROM ubuntu:24.04
WORKDIR /app
COPY --from=backend-builder /app/server/blog4 /app/
COPY --from=frontend-builder /app/dist /app/internal/admin/frontend/dist
COPY server/static /app/server/static
RUN apt-get update && apt-get install -y tzdata mysql-client openssl ca-certificates

ARG GIT_HASH
ENV GIT_HASH=$GIT_HASH

EXPOSE 8181
CMD ["/app/blog4"]
