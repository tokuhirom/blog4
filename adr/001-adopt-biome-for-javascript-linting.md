# ADR 001: Adopt Biome for JavaScript Linting and Formatting

## Status

Accepted

## Date

2026-01-04

## Context

Blog4 includes JavaScript code (Service Worker) that requires linting and formatting to maintain code quality. We need a tool that:

- Provides fast linting and formatting for JavaScript
- Works well in CI/CD environments
- Doesn't require managing package.json or npm dependencies
- Integrates easily with our existing workflow

We evaluated the following options:

1. **ESLint** - The de facto standard for JavaScript linting
2. **oxc** - A fast Rust-based linter
3. **Biome** - A fast Rust-based formatter and linter

## Decision

We will adopt **Biome** for JavaScript linting and formatting.

## Rationale

### Docker Image Availability

Biome provides an official Docker image (`ghcr.io/biomejs/biome`), which allows us to:
- Run linting without installing Node.js or managing package.json
- Ensure consistent tool versions across local development and CI
- Avoid dependency management overhead
- Keep the repository lightweight

### Performance

Biome is written in Rust and offers excellent performance:
- Fast startup time
- Efficient parallel processing
- Suitable for CI/CD pipelines

### Comparison with Alternatives

| Tool | Docker Image | Speed | Maturity |
|------|--------------|-------|----------|
| ESLint | ❌ No official image | Slower (Node-based) | Very mature |
| oxc | ⚠️ Limited | Very fast | Less mature |
| Biome | ✅ Official image | Very fast | Growing maturity |

ESLint requires Node.js ecosystem and package.json management, which adds unnecessary complexity for our minimal JavaScript usage.

oxc is fast but lacks an official Docker image and has less ecosystem maturity.

Biome provides the best balance of:
- Docker image availability
- Performance
- Ease of integration

### Integration

Biome integrates seamlessly with our Makefile-based workflow:

```make
biome-check:
    docker run --rm -v $(PWD):/app -w /app ghcr.io/biomejs/biome:latest check admin/static/sw.js
```

This approach:
- Requires no additional installation
- Works identically in local development and CI
- Is easy to understand and maintain

## Consequences

### Positive

- No need to manage package.json or npm dependencies
- Fast linting in both local and CI environments
- Simple integration with existing Makefile workflow
- Consistent tooling across environments via Docker

### Negative

- Biome is less mature than ESLint (fewer plugins, smaller ecosystem)
- Limited to Biome's built-in rules (no plugin ecosystem)
- For complex JavaScript projects, ESLint might be more suitable

### Neutral

- Our JavaScript codebase is minimal (only Service Worker)
- Biome's built-in rules are sufficient for our needs
- Can reconsider if JavaScript usage grows significantly

## Notes

- Biome configuration: `biome.json`
- Make targets: `make biome-check` and `make biome-fix`
- CI integration: `.github/workflows/ci.yml` (biome-check job)
- Checked files: `admin/static/sw.js`
