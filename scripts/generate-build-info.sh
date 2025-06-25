#!/bin/bash

# Generate build information
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT=$(git rev-parse HEAD 2>/dev/null || echo "unknown")
GIT_SHORT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GIT_BRANCH=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "unknown")
GIT_TAG=$(git describe --tags --exact-match 2>/dev/null || echo "")
GITHUB_URL="https://github.com/tokuhirom/blog4/commit/${GIT_COMMIT}"

# Create build info JSON file
cat > build-info.json <<EOF
{
  "buildTime": "${BUILD_TIME}",
  "gitCommit": "${GIT_COMMIT}",
  "gitShortCommit": "${GIT_SHORT_COMMIT}",
  "gitBranch": "${GIT_BRANCH}",
  "gitTag": "${GIT_TAG}",
  "githubUrl": "${GITHUB_URL}"
}
EOF

echo "Build info generated:"
cat build-info.json