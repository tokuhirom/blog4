# Stage 1: Build the Go backend
FROM golang:1.25 AS backend-builder
RUN apt-get update && apt-get install -y git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
# Generate build info
RUN bash scripts/generate-build-info.sh
RUN go build -o /app/blog4 ./cmd/blog4

# Stage 2: Final stage
FROM ubuntu:24.04
WORKDIR /app
COPY --from=backend-builder /app/blog4 /app/
COPY --from=backend-builder /app/build-info.json /app/
COPY web/static /app/web/static
COPY admin/static /app/admin/static
COPY web/templates /app/web/templates
RUN apt-get update && apt-get install -y tzdata mysql-client openssl ca-certificates

ARG GIT_HASH
ENV GIT_HASH=$GIT_HASH

EXPOSE 8181
CMD ["/app/blog4"]
