# Stage 1: Build the Go backend
FROM golang:1.26 AS backend-builder
RUN apt-get update && apt-get install -y git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
# Generate build info
RUN bash scripts/generate-build-info.sh
RUN go build -o /app/blog4 ./cmd/blog4

# Stage 2: Final stage
FROM debian:trixie-slim
WORKDIR /app
COPY --from=backend-builder /app/blog4 /app/
COPY --from=backend-builder /app/build-info.json /app/
COPY public /app/public
COPY admin /app/admin
RUN apt-get update && apt-get install -y \
    tzdata \
    mariadb-client \
    openssl \
    ca-certificates \
    fonts-noto-cjk \
    fonts-ipafont-gothic \
    fonts-ipafont-mincho \
    && rm -rf /var/lib/apt/lists/*

ARG GIT_HASH
ENV GIT_HASH=$GIT_HASH

EXPOSE 8181
CMD ["/app/blog4"]
